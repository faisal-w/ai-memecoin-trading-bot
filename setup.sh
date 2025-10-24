#!/bin/bash
# Setup script for meme_bot

echo "🎭 Meme Bot Setup Script"
echo "========================"
echo ""

# Check if Python is installed
if ! command -v python3 &> /dev/null; then
    echo "❌ Python 3 is not installed. Please install Python 3.8 or higher."
    exit 1
fi

echo "✅ Python 3 found: $(python3 --version)"
echo ""

# Check if pip is installed
if ! command -v pip3 &> /dev/null; then
    echo "❌ pip3 is not installed. Please install pip."
    exit 1
fi

echo "✅ pip3 found"
echo ""

# Install dependencies
echo "📦 Installing dependencies..."
pip3 install -r requirements.txt

if [ $? -eq 0 ]; then
    echo "✅ Dependencies installed successfully"
else
    echo "❌ Failed to install dependencies"
    exit 1
fi

echo ""

# Create .env file if it doesn't exist
if [ ! -f .env ]; then
    echo "📝 Creating .env file from template..."
    cp .env.example .env
    echo "✅ .env file created"
    echo "⚠️  Please edit .env file and add your bot token(s)"
else
    echo "ℹ️  .env file already exists"
fi

echo ""

# Create directories if they don't exist
echo "📁 Creating directories..."
mkdir -p meme_templates output
echo "✅ Directories created"

echo ""
echo "🎉 Setup complete!"
echo ""
echo "Next steps:"
echo "1. Edit .env file and add your bot token(s)"
echo "2. Add meme templates to meme_templates/ directory"
echo "3. Run the bot:"
echo "   - For Discord: python3 main.py discord"
echo "   - For Telegram: python3 main.py telegram"
echo ""
echo "For more information, see README.md or QUICKSTART.md"
