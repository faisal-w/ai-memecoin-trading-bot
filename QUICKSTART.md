# Quick Start Guide

This guide will help you get the meme bot up and running quickly.

## Step 1: Install Dependencies

```bash
pip install -r requirements.txt
```

## Step 2: Configure Bot Token

1. Copy the example environment file:
```bash
cp .env.example .env
```

2. Choose your platform and get a token:

### For Discord:
- Visit https://discord.com/developers/applications
- Create a new application
- Go to "Bot" section and create a bot
- Copy the token and paste it in `.env` as `DISCORD_TOKEN`
- Enable "Message Content Intent" in Bot settings
- Invite bot to your server using OAuth2 URL generator

### For Telegram:
- Open Telegram and message @BotFather
- Send `/newbot` and follow instructions
- Copy the token and paste it in `.env` as `TELEGRAM_TOKEN`

## Step 3: Add Meme Templates (Optional)

Create a `meme_templates` directory and add your meme images:

```bash
mkdir meme_templates
# Copy your meme template images here
```

## Step 4: Run the Bot

For Discord:
```bash
python main.py discord
```

For Telegram:
```bash
python main.py telegram
```

## Step 5: Test the Bot

### Discord:
Type in any channel where the bot is present:
```
!help_meme
!templates
```

### Telegram:
Message your bot:
```
/start
/help
/templates
```

## Common Commands

### Creating a Meme from Template:

**Discord:**
```
!meme template_name.jpg "Top Text" "Bottom Text"
```

**Telegram:**
```
/meme template_name.jpg
(then reply with): Top Text | Bottom Text
```

### Creating a Meme from URL:

**Discord:**
```
!memeurl https://example.com/image.jpg "Top Text" "Bottom Text"
```

**Telegram:**
```
/memeurl https://example.com/image.jpg
(then reply with): Top Text | Bottom Text
```

That's it! You're ready to create memes! 🎉
