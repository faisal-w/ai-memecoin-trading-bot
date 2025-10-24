# Meme Bot 🎭

A versatile project featuring:
1. **Meme Generation Bot** - Discord and Telegram bot for creating memes with a modern web interface
2. **🚀 NEW: Meme Coin Trading Bot** - Automated trading system for Solana and Base meme coins with AI-powered strategy evaluation

Built with **Go (Golang)** for the backend and **Rust** for the frontend!

## Features

### Meme Generation Bot
- 🎨 Create memes from predefined templates
- 🌐 Generate memes from any image URL
- 💬 Support for both Discord and Telegram
- 🖼️ Classic meme text styling (white text with black outline)
- 📝 Easy-to-use command interface
- 🚀 Web-based frontend built with Rust (Yew framework)
- ⚡ High-performance Go backend
- 🔌 RESTful API for frontend integration

### 🆕 Meme Coin Trading Bot
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

**Backend:**
- Go (Golang) 1.20+
- Discord bot library: discordgo
- Telegram bot library: telegram-bot-api
- HTTP server: gorilla/mux
- Image processing: golang/freetype

**Frontend:**
- Rust with Yew framework (WebAssembly)
- Modern, responsive UI
- Real-time meme generation

## Prerequisites

- Go 1.20 or higher
- Rust 1.70 or higher (with cargo) - *Optional, only for frontend*
- Trunk (Rust WASM bundler): `cargo install trunk` - *Optional, only for frontend*

**For Meme Generation Bot:**
- A Discord Bot Token (for Discord) or Telegram Bot Token (for Telegram)

**For Trading Bot:**
- Solana and Base RPC endpoints
- OKX Wallet SDK or private key (for live trading)
- API keys (CoinGecko, Twitter - optional)

## Quick Start

### Meme Generation Bot

See the [original documentation below](#meme-generation-bot-usage) for Discord/Telegram/Web interface setup.

### 🚀 Meme Coin Trading Bot

**⚠️ Start in Dry-Run Mode (Recommended)**

```bash
# 1. Configure environment
cp .env.example .env
# Edit .env: ensure DRY_RUN=true, AUTO_EXECUTE=false

# 2. Build
make build-backend

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

## Meme Generation Bot Usage

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

3. Edit `.env` file and add your bot tokens:
   - For Discord: Add your `DISCORD_TOKEN`
   - For Telegram: Add your `TELEGRAM_TOKEN`

4. Build the project:
```bash
make build
```

This will build both the Go backend and Rust frontend.

## Configuration

Edit the `.env` file with your credentials:

```env
# Discord Bot Token (if using Discord)
DISCORD_TOKEN=your_discord_bot_token_here

# Telegram Bot Token (if using Telegram)
TELEGRAM_TOKEN=your_telegram_bot_token_here

# Bot Settings
BOT_PREFIX=!

# Server Settings (for web interface)
SERVER_PORT=8080

# Directories
TEMPLATES_DIR=meme_templates
OUTPUT_DIR=output
```

### Getting Bot Tokens

**Discord:**
1. Go to [Discord Developer Portal](https://discord.com/developers/applications)
2. Create a new application
3. Go to the "Bot" section and create a bot
4. Copy the bot token
5. Enable "Message Content Intent" in Bot settings

**Telegram:**
1. Open Telegram and search for [@BotFather](https://t.me/botfather)
2. Send `/newbot` command
3. Follow the instructions to create your bot
4. Copy the bot token provided

## Usage

### Running the Web Interface

```bash
make run-server
```

Then open your browser to `http://localhost:8080` to access the web interface!

### Running the Discord Bot

```bash
make run-discord
```

### Running the Telegram Bot

```bash
make run-telegram
```

## Commands

### Discord Commands

- `!meme <template_name> "top text" "bottom text"` - Create a meme from a template
- `!templates` - List all available meme templates
- `!memeurl <image_url> "top text" "bottom text"` - Create a meme from an image URL
- `!help_meme` - Show help information

**Examples:**
```
!meme drake.jpg "studying for exams" "browsing memes"
!templates
!memeurl https://example.com/funny-image.jpg "when you" "find a bug in production"
```

### Telegram Commands

- `/start` - Welcome message and bot introduction
- `/meme <template_name>` - Create a meme from a template (then reply with text)
- `/templates` - List all available meme templates
- `/memeurl <image_url>` - Create a meme from an image URL (then reply with text)
- `/help` - Show help information

**Examples:**
```
/meme drake.jpg
(then reply): studying for exams | browsing memes

/memeurl https://example.com/funny-image.jpg
(then reply): when you | find a bug in production
```

## Adding Meme Templates

1. Create a `meme_templates` directory in the project root (if it doesn't exist)
2. Add your meme template images (PNG, JPG, or JPEG format)
3. Use the filename (without path) in your commands

Example:
```
meme_templates/
  ├── drake.jpg
  ├── distracted-boyfriend.jpg
  └── two-buttons.jpg
```

## Project Structure

```
meme_bot/
├── cmd/                    # Go command-line applications
│   ├── server/            # HTTP API server with frontend
│   ├── discord/           # Discord bot
│   └── telegram/          # Telegram bot
├── internal/              # Internal Go packages
│   └── config/           # Configuration management
├── pkg/                   # Public Go packages
│   └── meme/             # Meme generation logic
├── frontend/              # Rust frontend (Yew)
│   ├── src/              # Rust source code
│   ├── Cargo.toml        # Rust dependencies
│   └── index.html        # HTML template
├── meme_templates/        # Directory for meme templates
├── output/                # Generated memes (auto-created)
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
# Build everything
make build

# Build only backend
make build-backend

# Build only frontend
make build-frontend
```

### Running in Development Mode

For the frontend with hot reload:
```bash
cd frontend
trunk serve
```

For the backend:
```bash
# Build first
make build-backend

# Then run
./bin/server
```

### Testing

```bash
# Go tests
go test ./...

# Rust tests  
cd frontend
cargo test
```

## Troubleshooting

### Common Issues

1. **Bot doesn't respond:**
   - Check if your bot token is correctly set in `.env`
   - Ensure the bot has necessary permissions (Discord: Message Content Intent enabled)
   - Verify the bot is online and connected

2. **Template not found:**
   - Make sure the template file exists in the `meme_templates` directory
   - Check the filename spelling and case sensitivity
   - Use the `!templates` or `/templates` command to list available templates

3. **Build errors:**
   - Go: Make sure you have Go 1.20+ installed
   - Rust: Ensure you have Rust 1.70+ and trunk installed
   - Run `go mod tidy` and `cargo clean` then rebuild

4. **Frontend doesn't load:**
   - Make sure you built the frontend: `make build-frontend`
   - Check that the server is running on the correct port
   - Verify `frontend/dist` directory exists and contains files

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

Created with ❤️ for meme lovers everywhere!

## Support

If you encounter any issues or have questions, please open an issue on GitHub.

---

**Happy Meme Making! 🎉**