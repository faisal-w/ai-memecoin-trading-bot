"""Configuration module for meme bot."""
import os
from dotenv import load_dotenv

# Load environment variables from .env file
load_dotenv()

# Bot tokens
DISCORD_TOKEN = os.getenv('DISCORD_TOKEN')
TELEGRAM_TOKEN = os.getenv('TELEGRAM_TOKEN')

# Bot settings
BOT_PREFIX = os.getenv('BOT_PREFIX', '!')
ADMIN_IDS = os.getenv('ADMIN_IDS', '').split(',')

# Meme API credentials
IMGFLIP_USERNAME = os.getenv('IMGFLIP_USERNAME')
IMGFLIP_PASSWORD = os.getenv('IMGFLIP_PASSWORD')

# Paths
MEME_TEMPLATES_DIR = 'meme_templates'
OUTPUT_DIR = 'output'
