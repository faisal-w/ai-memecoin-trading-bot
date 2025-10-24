"""Discord bot implementation for meme bot."""
import discord
from discord.ext import commands
import os
from meme_generator import MemeGenerator, generate_meme_from_url
import config


class DiscordMemeBot:
    """Discord meme bot class."""
    
    def __init__(self):
        """Initialize Discord bot."""
        intents = discord.Intents.default()
        intents.message_content = True
        
        self.bot = commands.Bot(command_prefix=config.BOT_PREFIX, intents=intents)
        self.meme_gen = MemeGenerator()
        
        # Register commands
        self._register_commands()
    
    def _register_commands(self):
        """Register bot commands."""
        
        @self.bot.event
        async def on_ready():
            print(f'{self.bot.user} has connected to Discord!')
            print(f'Bot is in {len(self.bot.guilds)} guilds')
        
        @self.bot.command(name='meme', help='Create a meme: !meme <template_name> "top text" "bottom text"')
        async def create_meme(ctx, template_name: str, top_text: str = '', bottom_text: str = ''):
            """Create a meme from a template."""
            try:
                output_path = self.meme_gen.create_meme(template_name, top_text, bottom_text)
                await ctx.send(file=discord.File(output_path))
                
                # Clean up generated file
                if os.path.exists(output_path):
                    os.remove(output_path)
                    
            except FileNotFoundError as e:
                await ctx.send(f"❌ Error: {str(e)}")
            except Exception as e:
                await ctx.send(f"❌ An error occurred: {str(e)}")
        
        @self.bot.command(name='templates', help='List available meme templates')
        async def list_templates(ctx):
            """List all available meme templates."""
            templates = self.meme_gen.list_templates()
            
            if templates:
                template_list = '\n'.join([f"• {t}" for t in templates])
                await ctx.send(f"**Available Templates:**\n{template_list}")
            else:
                await ctx.send("No templates found. Add templates to the meme_templates directory.")
        
        @self.bot.command(name='memeurl', help='Create meme from URL: !memeurl <image_url> "top text" "bottom text"')
        async def create_meme_from_url(ctx, image_url: str, top_text: str = '', bottom_text: str = ''):
            """Create a meme from an image URL."""
            try:
                await ctx.send("🎨 Generating meme...")
                img = generate_meme_from_url(image_url, top_text, bottom_text)
                
                # Save to temporary file
                temp_path = os.path.join(config.OUTPUT_DIR, 'temp_meme.png')
                img.save(temp_path)
                
                await ctx.send(file=discord.File(temp_path))
                
                # Clean up
                if os.path.exists(temp_path):
                    os.remove(temp_path)
                    
            except Exception as e:
                await ctx.send(f"❌ An error occurred: {str(e)}")
        
        @self.bot.command(name='help_meme', help='Show meme bot help')
        async def help_meme(ctx):
            """Show help information."""
            help_text = """
**Meme Bot Commands:**

`!meme <template> "top text" "bottom text"` - Create a meme from a template
`!templates` - List all available templates
`!memeurl <url> "top text" "bottom text"` - Create a meme from an image URL
`!help_meme` - Show this help message

**Examples:**
`!meme drake.jpg "studying" "watching memes"`
`!memeurl https://example.com/image.jpg "top" "bottom"`
            """
            await ctx.send(help_text)
    
    def run(self):
        """Start the Discord bot."""
        if not config.DISCORD_TOKEN:
            print("Error: DISCORD_TOKEN not found in environment variables")
            print("Please set your Discord token in .env file")
            return
        
        print("Starting Discord bot...")
        self.bot.run(config.DISCORD_TOKEN)


if __name__ == '__main__':
    bot = DiscordMemeBot()
    bot.run()
