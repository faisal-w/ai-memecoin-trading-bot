# Quick Start Guide

Get the meme coin trading bot up and running in minutes!

## Prerequisites

- **Go:** Version 1.20 or higher

## Step 1: Installation

```bash
# Clone the repository
git clone https://github.com/mumugogoing/meme_bot.git
cd meme_bot

# Build the project
make build
```

## Step 2: Configuration

**⚠️ IMPORTANT: Start in DRY-RUN mode**

```bash
# Copy the example environment file
cp .env.example .env

# Edit .env with your configuration
# Make sure DRY_RUN=true and AUTO_EXECUTE=false
```

See [TRADING_BOT.md](TRADING_BOT.md) for detailed configuration options.

## Step 3: Run the Trading Bot

```bash
make run-trading
```

The bot will start monitoring Solana and Base chains for new meme coins.

## Step 4: Monitor via API

Access the monitoring API:

- **Health Check:** http://localhost:8080/api/health
- **Bot Status:** http://localhost:8080/api/status
- **Candidates:** http://localhost:8080/api/candidates
- **Metrics:** http://localhost:8080/api/metrics
- **Risk Status:** http://localhost:8080/api/risk

## Quick Commands Reference

```bash
make build              # Build the trading bot
make run-trading        # Run the trading bot
make test               # Run tests
make clean              # Clean build artifacts
```

## Configuration Tips

### For Testing (Recommended First):
```env
DRY_RUN=true
AUTO_EXECUTE=false
WIN_PROBABILITY_THRESHOLD=0.80
ACCOUNT_BALANCE=10000.0
```

### For Small Live Trading:
```env
DRY_RUN=false
AUTO_EXECUTE=true
ACCOUNT_BALANCE=100.0
SINGLE_POSITION_PCT=0.01
WIN_PROBABILITY_THRESHOLD=0.85
DAILY_LOSS_LIMIT=10.0
```

## Troubleshooting

**Build fails?**
- Install Go 1.20+: https://go.dev/dl/
- Run `go mod tidy` and try again

**Bot won't start?**
- Check if port 8080 is available
- Verify RPC endpoints in `.env`
- Check logs for error messages

**No trades executing?**
- Verify DRY_RUN and AUTO_EXECUTE settings
- Check WIN_PROBABILITY_THRESHOLD (lower = more trades)
- Review API endpoints for candidate tokens

That's it! You're ready to start trading! 🚀

For more details, see [TRADING_BOT.md](TRADING_BOT.md) and [README.md](README.md)
