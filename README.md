# Meme Bot 🎭

A versatile meme generation bot that works with both Discord and Telegram platforms. Create hilarious memes from templates or custom images with simple commands!

## Features

- 🎨 Create memes from predefined templates
- 🌐 Generate memes from any image URL
- 💬 Support for both Discord and Telegram
- 🖼️ Classic meme text styling (white text with black outline)
- 📝 Easy-to-use command interface
- 🔧 Configurable with environment variables

## Prerequisites

- Python 3.8 or higher
- pip (Python package manager)
- A Discord Bot Token (for Discord) or Telegram Bot Token (for Telegram)

## Installation

1. Clone the repository:
```bash
git clone https://github.com/mumugogoing/meme_bot.git
cd meme_bot
```

2. Install required dependencies:
```bash
pip install -r requirements.txt
```

3. Set up environment variables:
```bash
cp .env.example .env
```

4. Edit `.env` file and add your bot tokens:
   - For Discord: Add your `DISCORD_TOKEN`
   - For Telegram: Add your `TELEGRAM_TOKEN`

## Configuration

Edit the `.env` file with your credentials:

```env
# Discord Bot Token (if using Discord)
DISCORD_TOKEN=your_discord_bot_token_here

# Telegram Bot Token (if using Telegram)
TELEGRAM_TOKEN=your_telegram_bot_token_here

# Bot Settings
BOT_PREFIX=!
ADMIN_IDS=123456789,987654321
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

### Running the Discord Bot

```bash
python main.py discord
```

### Running the Telegram Bot

```bash
python main.py telegram
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
├── main.py                 # Main entry point
├── discord_bot.py          # Discord bot implementation
├── telegram_bot.py         # Telegram bot implementation
├── meme_generator.py       # Meme generation logic
├── config.py               # Configuration management
├── requirements.txt        # Python dependencies
├── .env.example           # Example environment variables
├── .gitignore             # Git ignore rules
├── README.md              # This file
├── meme_templates/        # Directory for meme templates (create this)
└── output/                # Generated memes (auto-created)
```

## Development

### Running Tests

```bash
python -m pytest tests/
```

### Code Style

This project follows PEP 8 guidelines. You can check your code with:

```bash
flake8 .
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

3. **Font issues:**
   - The bot uses system fonts. On Linux, install `fonts-dejavu` package
   - On Windows/Mac, default fonts should work automatically

4. **Image generation fails:**
   - Check if Pillow is correctly installed: `pip install --upgrade Pillow`
   - Ensure write permissions for the `output` directory

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