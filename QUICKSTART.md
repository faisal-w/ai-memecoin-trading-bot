# Quick Start Guide

Get the meme bot up and running in minutes!

## Prerequisites

- **Go:** Version 1.20 or higher
- **Rust:** Version 1.70 or higher  
- **Trunk:** Install with `cargo install trunk`

## Step 1: Installation

```bash
# Clone the repository
git clone https://github.com/mumugogoing/meme_bot.git
cd meme_bot

# Build the project
make build
```

This will build both the Go backend and Rust frontend.

## Step 2: Configuration (Optional)

For Discord/Telegram bots, you'll need tokens:

```bash
# Copy the example environment file
cp .env.example .env

# Edit .env and add your tokens
```

### For Discord:
- Visit https://discord.com/developers/applications
- Create a new application
- Go to "Bot" section and create a bot
- Copy the token and paste it in `.env` as `DISCORD_TOKEN`
- Enable "Message Content Intent" in Bot settings

### For Telegram:
- Open Telegram and message @BotFather
- Send `/newbot` and follow instructions
- Copy the token and paste it in `.env` as `TELEGRAM_TOKEN`

## Step 3: Run the Application

### Option 1: Web Interface (Easiest!)

```bash
make run-server
```

Then open your browser to: **http://localhost:8080**

You can now generate memes through the web interface!

### Option 2: Discord Bot

```bash
make run-discord
```

Commands in Discord:
- `!meme sample.png "Top Text" "Bottom Text"`
- `!templates`
- `!help_meme`

### Option 3: Telegram Bot

```bash
make run-telegram
```

Commands in Telegram:
- `/start`
- `/meme sample.png` (then reply with: Top Text | Bottom Text)
- `/templates`
- `/help`

## Step 4: Add Your Own Templates (Optional)

Add meme template images to the `meme_templates/` directory:

```bash
# Copy your meme images
cp your-meme.jpg meme_templates/

# The templates will be automatically available!
```

Supported formats: PNG, JPG, JPEG, GIF

## Quick Commands Reference

```bash
make build              # Build everything
make run-server         # Run web interface
make run-discord        # Run Discord bot
make run-telegram       # Run Telegram bot
make clean             # Clean build artifacts
```

## Example Usage

### Web Interface:
1. Select a template from the dropdown
2. Enter your top and bottom text
3. Click "Generate Meme"
4. Download your meme!

### Discord:
```
!meme sample.png "When you" "Use Go and Rust"
```

### Telegram:
```
/meme sample.png
(reply with): When you | Use Go and Rust
```

## Troubleshooting

**Build fails?**
- Install Go 1.20+: https://go.dev/dl/
- Install Rust: https://rustup.rs/
- Install trunk: `cargo install trunk`
- Add wasm target: `rustup target add wasm32-unknown-unknown`

**Server won't start?**
- Check if port 8080 is available
- Set a different port in `.env`: `SERVER_PORT=3000`

**Bots don't respond?**
- Verify tokens in `.env` file
- Check bot permissions on Discord
- Make sure bot is running in terminal

That's it! You're ready to create memes! 🎉

For more details, see [README.md](README.md)
