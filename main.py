"""Main entry point for the meme bot."""
import sys
import argparse
from discord_bot import DiscordMemeBot
from telegram_bot import TelegramMemeBot


def main():
    """Main function to run the bot."""
    parser = argparse.ArgumentParser(description='Meme Bot - Discord and Telegram bot for generating memes')
    parser.add_argument(
        'platform',
        choices=['discord', 'telegram'],
        help='Platform to run the bot on (discord or telegram)'
    )
    
    args = parser.parse_args()
    
    if args.platform == 'discord':
        print("Starting Discord Meme Bot...")
        bot = DiscordMemeBot()
        bot.run()
    elif args.platform == 'telegram':
        print("Starting Telegram Meme Bot...")
        bot = TelegramMemeBot()
        bot.run()


if __name__ == '__main__':
    main()
