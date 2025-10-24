# Deployment Guide

This guide covers different deployment scenarios for the Meme Coin Trading Bot.

## Table of Contents
- [Local Development](#local-development)
- [Production Server](#production-server)
- [Docker Deployment](#docker-deployment)
- [Cloud Deployment](#cloud-deployment)
- [Monitoring Setup](#monitoring-setup)

## Local Development

### Prerequisites
```bash
# Install Go
curl -OL https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# Verify installation
go version
```

### Running Locally
```bash
# 1. Clone repository
git clone https://github.com/mumugogoing/meme_bot.git
cd meme_bot

# 2. Configure environment
cp .env.example .env
# Edit .env with your settings (use DRY_RUN=true initially)

# 3. Build
make build-backend

# 4. Run
./bin/trading

# Or use make
make run-trading
```

### Development Mode
```bash
# Watch for changes and rebuild
while true; do
  make build-backend && ./bin/trading
  sleep 2
done
```

## Production Server

### Server Requirements
- **OS**: Ubuntu 20.04+ or similar Linux distribution
- **CPU**: 2+ cores
- **RAM**: 4GB minimum, 8GB recommended
- **Storage**: 50GB SSD
- **Network**: Stable connection, low latency to RPC endpoints

### Setup Script
```bash
#!/bin/bash
# setup.sh - Production server setup

set -e

echo "Installing dependencies..."
sudo apt update
sudo apt install -y git build-essential curl

echo "Installing Go..."
cd /tmp
curl -OL https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

echo "Creating app directory..."
sudo mkdir -p /opt/meme_bot
sudo chown $USER:$USER /opt/meme_bot
cd /opt/meme_bot

echo "Cloning repository..."
git clone https://github.com/mumugogoing/meme_bot.git .

echo "Building application..."
make build-backend

echo "Creating systemd service..."
sudo tee /etc/systemd/system/trading-bot.service > /dev/null <<EOF
[Unit]
Description=Meme Coin Trading Bot
After=network.target

[Service]
Type=simple
User=$USER
WorkingDirectory=/opt/meme_bot
Environment="PATH=/usr/local/go/bin:/usr/bin:/bin"
ExecStart=/opt/meme_bot/bin/trading
Restart=always
RestartSec=10
StandardOutput=append:/var/log/trading-bot/output.log
StandardError=append:/var/log/trading-bot/error.log

[Install]
WantedBy=multi-user.target
EOF

echo "Creating log directory..."
sudo mkdir -p /var/log/trading-bot
sudo chown $USER:$USER /var/log/trading-bot

echo "Setup complete!"
echo "Next steps:"
echo "1. Edit /opt/meme_bot/.env with your configuration"
echo "2. sudo systemctl enable trading-bot"
echo "3. sudo systemctl start trading-bot"
echo "4. sudo systemctl status trading-bot"
```

### Systemd Service Management
```bash
# Enable service
sudo systemctl enable trading-bot

# Start service
sudo systemctl start trading-bot

# Check status
sudo systemctl status trading-bot

# View logs
journalctl -u trading-bot -f

# Stop service
sudo systemctl stop trading-bot

# Restart service
sudo systemctl restart trading-bot
```

### Log Rotation
```bash
# Create logrotate configuration
sudo tee /etc/logrotate.d/trading-bot > /dev/null <<EOF
/var/log/trading-bot/*.log {
    daily
    rotate 7
    compress
    delaycompress
    missingok
    notifempty
    create 0640 $USER $USER
    sharedscripts
    postrotate
        systemctl reload trading-bot > /dev/null 2>&1 || true
    endscript
}
EOF
```

## Docker Deployment

### Dockerfile
Create `Dockerfile` in project root:

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o trading ./cmd/trading

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/trading .

# Expose API port
EXPOSE 8080

# Run
CMD ["./trading"]
```

### Docker Compose
Create `docker-compose.yml`:

```yaml
version: '3.8'

services:
  trading-bot:
    build: .
    container_name: meme-trading-bot
    restart: unless-stopped
    ports:
      - "8080:8080"
      - "9090:9090"  # Prometheus metrics
    env_file:
      - .env
    volumes:
      - ./data:/root/data
      - ./logs:/root/logs
    networks:
      - trading-network
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"

networks:
  trading-network:
    driver: bridge
```

### Docker Commands
```bash
# Build image
docker-compose build

# Start container
docker-compose up -d

# View logs
docker-compose logs -f

# Stop container
docker-compose down

# Restart
docker-compose restart

# Access container shell
docker-compose exec trading-bot sh
```

## Cloud Deployment

### AWS EC2

1. **Launch EC2 Instance**
   ```bash
   # Instance type: t3.medium or larger
   # OS: Ubuntu 20.04 LTS
   # Storage: 50GB SSD
   # Security Group: Allow port 8080, 22
   ```

2. **Connect and Setup**
   ```bash
   ssh -i your-key.pem ubuntu@your-ec2-ip
   
   # Run setup script
   curl -sSL https://raw.githubusercontent.com/mumugogoing/meme_bot/main/setup.sh | bash
   ```

3. **Configure Secrets (Using AWS Secrets Manager)**
   ```bash
   # Install AWS CLI
   sudo apt install awscli
   
   # Configure credentials
   aws configure
   
   # Store secrets
   aws secretsmanager create-secret \
     --name trading-bot-private-key \
     --secret-string "your-private-key"
   
   # Retrieve in application
   # Modify code to fetch from Secrets Manager
   ```

### Google Cloud Platform (GCP)

1. **Create Compute Engine Instance**
   ```bash
   gcloud compute instances create trading-bot \
     --machine-type=n1-standard-2 \
     --image-family=ubuntu-2004-lts \
     --image-project=ubuntu-os-cloud \
     --boot-disk-size=50GB \
     --tags=http-server
   ```

2. **Setup Firewall**
   ```bash
   gcloud compute firewall-rules create allow-trading-bot \
     --allow=tcp:8080 \
     --target-tags=http-server
   ```

3. **Use Secret Manager**
   ```bash
   # Store secret
   echo -n "your-private-key" | \
     gcloud secrets create trading-bot-key --data-file=-
   
   # Grant access to service account
   gcloud secrets add-iam-policy-binding trading-bot-key \
     --member="serviceAccount:your-service-account" \
     --role="roles/secretmanager.secretAccessor"
   ```

### Kubernetes Deployment

Create `k8s-deployment.yaml`:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: trading-bot-config
data:
  DRY_RUN: "true"
  AUTO_EXECUTE: "false"
  # Add other non-sensitive config

---
apiVersion: v1
kind: Secret
metadata:
  name: trading-bot-secrets
type: Opaque
stringData:
  PRIVATE_KEY: "your-private-key"
  OKX_API_KEY: "your-okx-key"
  # Add other secrets

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: trading-bot
spec:
  replicas: 1
  selector:
    matchLabels:
      app: trading-bot
  template:
    metadata:
      labels:
        app: trading-bot
    spec:
      containers:
      - name: trading-bot
        image: your-registry/trading-bot:latest
        ports:
        - containerPort: 8080
        - containerPort: 9090
        envFrom:
        - configMapRef:
            name: trading-bot-config
        - secretRef:
            name: trading-bot-secrets
        resources:
          requests:
            memory: "512Mi"
            cpu: "500m"
          limits:
            memory: "2Gi"
            cpu: "1000m"
        livenessProbe:
          httpGet:
            path: /api/health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10

---
apiVersion: v1
kind: Service
metadata:
  name: trading-bot-service
spec:
  selector:
    app: trading-bot
  ports:
  - name: api
    port: 8080
    targetPort: 8080
  - name: metrics
    port: 9090
    targetPort: 9090
  type: LoadBalancer
```

Deploy:
```bash
kubectl apply -f k8s-deployment.yaml
kubectl get pods
kubectl logs -f deployment/trading-bot
```

## Monitoring Setup

### Prometheus

1. **Install Prometheus**
   ```bash
   # Download
   wget https://github.com/prometheus/prometheus/releases/download/v2.40.0/prometheus-2.40.0.linux-amd64.tar.gz
   tar xvfz prometheus-2.40.0.linux-amd64.tar.gz
   cd prometheus-2.40.0.linux-amd64
   ```

2. **Configure** (`prometheus.yml`):
   ```yaml
   global:
     scrape_interval: 15s
   
   scrape_configs:
     - job_name: 'trading-bot'
       static_configs:
         - targets: ['localhost:9090']
   ```

3. **Run**:
   ```bash
   ./prometheus --config.file=prometheus.yml
   # Access at http://localhost:9090
   ```

### Grafana

1. **Install**:
   ```bash
   sudo apt-get install -y software-properties-common
   sudo add-apt-repository "deb https://packages.grafana.com/oss/deb stable main"
   wget -q -O - https://packages.grafana.com/gpg.key | sudo apt-key add -
   sudo apt-get update
   sudo apt-get install grafana
   ```

2. **Start**:
   ```bash
   sudo systemctl enable grafana-server
   sudo systemctl start grafana-server
   # Access at http://localhost:3000 (admin/admin)
   ```

3. **Configure**:
   - Add Prometheus data source
   - Import trading bot dashboard
   - Set up alerts

### Alerting

Create `alerts.yml` for Prometheus:

```yaml
groups:
  - name: trading_bot
    interval: 30s
    rules:
      - alert: TradingHalted
        expr: trading_halted == 1
        for: 1m
        annotations:
          summary: "Trading has been halted"
          
      - alert: HighFailureRate
        expr: rate(execution_failures[5m]) > 0.5
        for: 5m
        annotations:
          summary: "High execution failure rate"
          
      - alert: DailyLossLimit
        expr: daily_loss > daily_loss_limit * 0.9
        for: 1m
        annotations:
          summary: "Approaching daily loss limit"
```

## Backup and Recovery

### Database Backup
```bash
# Backup script
#!/bin/bash
DATE=$(date +%Y%m%d_%H%M%S)
sqlite3 /opt/meme_bot/meme_bot.db ".backup /opt/meme_bot/backups/meme_bot_$DATE.db"

# Keep only last 7 days
find /opt/meme_bot/backups -name "meme_bot_*.db" -mtime +7 -delete
```

### Configuration Backup
```bash
# Backup .env (excluding sensitive data)
cp .env .env.backup
```

### Automated Backups
```bash
# Add to crontab
0 */6 * * * /opt/meme_bot/backup.sh
```

## Security Hardening

### Firewall Setup
```bash
# UFW firewall
sudo ufw allow ssh
sudo ufw allow 8080/tcp  # API
sudo ufw enable
```

### SSL/TLS (Nginx Reverse Proxy)
```bash
# Install nginx
sudo apt install nginx certbot python3-certbot-nginx

# Configure nginx
sudo tee /etc/nginx/sites-available/trading-bot <<EOF
server {
    listen 80;
    server_name your-domain.com;
    
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
    }
}
EOF

# Enable site
sudo ln -s /etc/nginx/sites-available/trading-bot /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx

# Get SSL certificate
sudo certbot --nginx -d your-domain.com
```

## Troubleshooting

### Common Issues

**Service won't start:**
```bash
# Check logs
journalctl -u trading-bot -n 50

# Check permissions
ls -la /opt/meme_bot

# Verify Go installation
go version
```

**High memory usage:**
```bash
# Monitor resources
htop

# Adjust limits in systemd service
[Service]
MemoryLimit=2G
```

**RPC connection issues:**
```bash
# Test connectivity
curl -X POST -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"getHealth"}' \
  $SOLANA_RPC_URL

# Use multiple endpoints
# Add fallback RPC URLs in configuration
```

## Maintenance

### Updates
```bash
# Stop service
sudo systemctl stop trading-bot

# Backup
./backup.sh

# Update code
cd /opt/meme_bot
git pull

# Rebuild
make build-backend

# Start service
sudo systemctl start trading-bot

# Check status
sudo systemctl status trading-bot
```

### Monitoring Checklist
- [ ] Check service status daily
- [ ] Review metrics weekly
- [ ] Update dependencies monthly
- [ ] Rotate logs
- [ ] Verify backups
- [ ] Review security advisories

---

**Need Help?** Open an issue on GitHub or consult the documentation.
