"""Telegram bot implementation for meme bot."""
from telegram import Update
from telegram.ext import Application, CommandHandler, ContextTypes, MessageHandler, filters
import os
from meme_generator import MemeGenerator, generate_meme_from_url
import config


class TelegramMemeBot:
    """Telegram meme bot class."""
    
    def __init__(self):
        """Initialize Telegram bot."""
        self.meme_gen = MemeGenerator()
        self.app = None
    
    async def start_command(self, update: Update, context: ContextTypes.DEFAULT_TYPE):
        """Handle /start command."""
        welcome_text = """
👋 Welcome to Meme Bot!

**Commands:**
/meme - Create a meme from template
/templates - List available templates
/memeurl - Create meme from URL
/help - Show this help message

**Usage Examples:**
/meme drake.jpg
Then send: top text | bottom text

/memeurl https://example.com/image.jpg
Then send: top text | bottom text
        """
        await update.message.reply_text(welcome_text)
    
    async def help_command(self, update: Update, context: ContextTypes.DEFAULT_TYPE):
        """Handle /help command."""
        help_text = """
**Meme Bot Help:**

**Commands:**
• `/meme <template_name>` - Create a meme from a template
  Then reply with: `top text | bottom text`
  
• `/templates` - List all available meme templates

• `/memeurl <image_url>` - Create a meme from an image URL
  Then reply with: `top text | bottom text`

• `/help` - Show this help message

**Examples:**
1. `/meme drake.jpg`
   Reply: `studying | watching memes`

2. `/memeurl https://example.com/image.jpg`
   Reply: `top text | bottom text`
        """
        await update.message.reply_text(help_text)
    
    async def templates_command(self, update: Update, context: ContextTypes.DEFAULT_TYPE):
        """Handle /templates command."""
        templates = self.meme_gen.list_templates()
        
        if templates:
            template_list = '\n'.join([f"• {t}" for t in templates])
            await update.message.reply_text(f"**Available Templates:**\n\n{template_list}")
        else:
            await update.message.reply_text("No templates found. Add templates to the meme_templates directory.")
    
    async def meme_command(self, update: Update, context: ContextTypes.DEFAULT_TYPE):
        """Handle /meme command."""
        if not context.args:
            await update.message.reply_text("❌ Please provide a template name.\nUsage: /meme <template_name>")
            return
        
        template_name = context.args[0]
        
        # Store template name in user data for next message
        context.user_data['pending_meme'] = {
            'type': 'template',
            'template_name': template_name
        }
        
        await update.message.reply_text(
            f"📝 Template: {template_name}\n"
            "Now send the text in format: `top text | bottom text`"
        )
    
    async def memeurl_command(self, update: Update, context: ContextTypes.DEFAULT_TYPE):
        """Handle /memeurl command."""
        if not context.args:
            await update.message.reply_text("❌ Please provide an image URL.\nUsage: /memeurl <image_url>")
            return
        
        image_url = context.args[0]
        
        # Store URL in user data for next message
        context.user_data['pending_meme'] = {
            'type': 'url',
            'image_url': image_url
        }
        
        await update.message.reply_text(
            "📝 Image URL received!\n"
            "Now send the text in format: `top text | bottom text`"
        )
    
    async def handle_text_message(self, update: Update, context: ContextTypes.DEFAULT_TYPE):
        """Handle text messages for meme generation."""
        if 'pending_meme' not in context.user_data:
            return
        
        text = update.message.text
        parts = text.split('|')
        
        top_text = parts[0].strip() if len(parts) > 0 else ''
        bottom_text = parts[1].strip() if len(parts) > 1 else ''
        
        pending = context.user_data['pending_meme']
        
        try:
            await update.message.reply_text("🎨 Generating meme...")
            
            if pending['type'] == 'template':
                output_path = self.meme_gen.create_meme(
                    pending['template_name'],
                    top_text,
                    bottom_text
                )
                
                with open(output_path, 'rb') as f:
                    await update.message.reply_photo(photo=f)
                
                # Clean up
                if os.path.exists(output_path):
                    os.remove(output_path)
                    
            elif pending['type'] == 'url':
                img = generate_meme_from_url(
                    pending['image_url'],
                    top_text,
                    bottom_text
                )
                
                # Save to temporary file
                temp_path = os.path.join(config.OUTPUT_DIR, 'temp_meme.png')
                img.save(temp_path)
                
                with open(temp_path, 'rb') as f:
                    await update.message.reply_photo(photo=f)
                
                # Clean up
                if os.path.exists(temp_path):
                    os.remove(temp_path)
            
            # Clear pending meme
            del context.user_data['pending_meme']
            
        except FileNotFoundError as e:
            await update.message.reply_text(f"❌ Error: {str(e)}")
            del context.user_data['pending_meme']
        except Exception as e:
            await update.message.reply_text(f"❌ An error occurred: {str(e)}")
            del context.user_data['pending_meme']
    
    def run(self):
        """Start the Telegram bot."""
        if not config.TELEGRAM_TOKEN:
            print("Error: TELEGRAM_TOKEN not found in environment variables")
            print("Please set your Telegram token in .env file")
            return
        
        print("Starting Telegram bot...")
        
        # Create application
        self.app = Application.builder().token(config.TELEGRAM_TOKEN).build()
        
        # Register handlers
        self.app.add_handler(CommandHandler("start", self.start_command))
        self.app.add_handler(CommandHandler("help", self.help_command))
        self.app.add_handler(CommandHandler("templates", self.templates_command))
        self.app.add_handler(CommandHandler("meme", self.meme_command))
        self.app.add_handler(CommandHandler("memeurl", self.memeurl_command))
        self.app.add_handler(MessageHandler(filters.TEXT & ~filters.COMMAND, self.handle_text_message))
        
        # Start bot
        print("Bot is running... Press Ctrl+C to stop")
        self.app.run_polling()


if __name__ == '__main__':
    bot = TelegramMemeBot()
    bot.run()
