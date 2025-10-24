# Meme Coin Trading Bot 🚀

Automated trading system for Solana and Base meme coins with AI-powered strategy evaluation.

Built with **Go (Golang)**!

## Features

- 🔍 Automated scanning of Solana and Base chains for new tokens
- 🛡️ Comprehensive honeypot detection and safety checks
- 📊 Win probability calculation (≥80% threshold)
- 💰 Automated trade execution with OKX Wallet SDK support
- ⚠️ Advanced risk management and circuit breakers
- 📈 Real-time metrics and monitoring via API
- 🔐 Security-first design with dry-run mode
- 📱 Multi-agent architecture for scalability

**[📖 See Trading Bot Documentation](TRADING_BOT.md)** | **[🏗️ Architecture Guide](ARCHITECTURE.md)**

## Tech Stack

- Go (Golang) 1.20+
- HTTP server: gorilla/mux
- Multi-agent architecture

## Prerequisites

- Go 1.20 or higher
- Solana and Base RPC endpoints
- OKX Wallet SDK or private key (for live trading)
- API keys (CoinGecko, Twitter - optional)

## Quick Start

**⚠️ Start in Dry-Run Mode (Recommended)**

```bash
# 1. Configure environment
cp .env.example .env
# Edit .env: ensure DRY_RUN=true, AUTO_EXECUTE=false

# 2. Build
make build

# 3. Run trading bot
make run-trading
```

**Access the API:**
- Health: http://localhost:8080/api/health
- Status: http://localhost:8080/api/status
- Candidates: http://localhost:8080/api/candidates
- Metrics: http://localhost:8080/api/metrics

**📚 Complete Guide:** See [TRADING_BOT.md](TRADING_BOT.md) for comprehensive documentation.

---

## Installation

1. Clone the repository:
```bash
git clone https://github.com/mumugogoing/meme_bot.git
cd meme_bot
```

2. Set up environment variables:
```bash
cp .env.example .env
```

3. Edit `.env` file with your trading bot configuration (see [TRADING_BOT.md](TRADING_BOT.md) for details)

4. Build the project:
```bash
make build
```

## Project Structure

```
meme_bot/
├── cmd/                    # Go command-line applications
│   └── trading/           # Trading bot
├── internal/              # Internal Go packages
│   └── config/           # Configuration management
├── pkg/                   # Public Go packages
│   ├── agents/           # Trading agents
│   ├── orchestrator/     # Orchestration logic
│   ├── models/           # Data models
│   └── config/           # Configuration
├── go.mod                 # Go module definition
├── go.sum                 # Go dependencies
├── Makefile               # Build automation
├── .env.example          # Example environment variables
├── .gitignore            # Git ignore rules
└── README.md             # This file
```

## Development

### Building the Project

```bash
make build
```

### Testing

```bash
go test ./...
```

## Troubleshooting

### Common Issues

1. **Build errors:**
   - Make sure you have Go 1.20+ installed
   - Run `go mod tidy` then rebuild

2. **Environment configuration:**
   - Check that all required environment variables are set in `.env`
   - See [TRADING_BOT.md](TRADING_BOT.md) for configuration details

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is open source and available under the MIT License.

## Credits

Created for automated meme coin trading!

## Support

If you encounter any issues or have questions, please open an issue on GitHub.

---

**Happy Trading! 🚀**