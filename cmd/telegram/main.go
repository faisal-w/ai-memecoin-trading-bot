package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mumugogoing/meme_bot/internal/config"
	"github.com/mumugogoing/meme_bot/pkg/meme"
)

type TelegramBot struct {
	bot       *tgbotapi.BotAPI
	generator *meme.Generator
	config    *config.Config
	// Store pending meme requests per user
	pendingMemes map[int64]*PendingMeme
}

type PendingMeme struct {
	Type     string // "template" or "url"
	Template string
	ImageURL string
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if cfg.TelegramToken == "" {
		log.Fatal("TELEGRAM_TOKEN not found in environment variables")
	}

	generator, err := meme.NewGenerator(cfg.TemplatesDir, cfg.OutputDir)
	if err != nil {
		log.Fatalf("Failed to create meme generator: %v", err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatalf("Failed to create Telegram bot: %v", err)
	}

	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)

	tb := &TelegramBot{
		bot:          bot,
		generator:    generator,
		config:       cfg,
		pendingMemes: make(map[int64]*PendingMeme),
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			tb.handleCommand(update.Message)
		} else {
			tb.handleText(update.Message)
		}
	}
}

func (tb *TelegramBot) handleCommand(message *tgbotapi.Message) {
	switch message.Command() {
	case "start":
		tb.handleStart(message)
	case "help":
		tb.handleHelp(message)
	case "templates":
		tb.handleTemplates(message)
	case "meme":
		tb.handleMemeCommand(message)
	case "memeurl":
		tb.handleMemeURLCommand(message)
	}
}

func (tb *TelegramBot) handleStart(message *tgbotapi.Message) {
	text := `👋 Welcome to Meme Bot!

**Commands:**
/meme - Create a meme from template
/templates - List available templates
/memeurl - Create meme from URL
/help - Show this help message

**Usage Examples:**
/meme drake.jpg
Then send: top text | bottom text

/memeurl https://example.com/image.jpg
Then send: top text | bottom text`

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ParseMode = "Markdown"
	tb.bot.Send(msg)
}

func (tb *TelegramBot) handleHelp(message *tgbotapi.Message) {
	text := `**Meme Bot Help:**

**Commands:**
• /meme <template_name> - Create a meme from a template
  Then reply with: top text | bottom text
  
• /templates - List all available meme templates

• /memeurl <image_url> - Create a meme from an image URL
  Then reply with: top text | bottom text

• /help - Show this help message

**Examples:**
1. /meme drake.jpg
   Reply: studying | watching memes

2. /memeurl https://example.com/image.jpg
   Reply: top text | bottom text`

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ParseMode = "Markdown"
	tb.bot.Send(msg)
}

func (tb *TelegramBot) handleTemplates(message *tgbotapi.Message) {
	templates, err := tb.generator.ListTemplates()
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "❌ Failed to list templates")
		tb.bot.Send(msg)
		return
	}

	if len(templates) == 0 {
		msg := tgbotapi.NewMessage(message.Chat.ID, "No templates found. Add templates to the meme_templates directory.")
		tb.bot.Send(msg)
		return
	}

	text := "**Available Templates:**\n\n"
	for _, t := range templates {
		text += fmt.Sprintf("• %s\n", t)
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ParseMode = "Markdown"
	tb.bot.Send(msg)
}

func (tb *TelegramBot) handleMemeCommand(message *tgbotapi.Message) {
	args := strings.Fields(message.CommandArguments())
	if len(args) == 0 {
		msg := tgbotapi.NewMessage(message.Chat.ID, "❌ Please provide a template name.\nUsage: /meme <template_name>")
		tb.bot.Send(msg)
		return
	}

	template := args[0]
	tb.pendingMemes[message.Chat.ID] = &PendingMeme{
		Type:     "template",
		Template: template,
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("📝 Template: %s\nNow send the text in format: `top text | bottom text`", template))
	msg.ParseMode = "Markdown"
	tb.bot.Send(msg)
}

func (tb *TelegramBot) handleMemeURLCommand(message *tgbotapi.Message) {
	args := strings.Fields(message.CommandArguments())
	if len(args) == 0 {
		msg := tgbotapi.NewMessage(message.Chat.ID, "❌ Please provide an image URL.\nUsage: /memeurl <image_url>")
		tb.bot.Send(msg)
		return
	}

	imageURL := args[0]
	tb.pendingMemes[message.Chat.ID] = &PendingMeme{
		Type:     "url",
		ImageURL: imageURL,
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, "📝 Image URL received!\nNow send the text in format: `top text | bottom text`")
	msg.ParseMode = "Markdown"
	tb.bot.Send(msg)
}

func (tb *TelegramBot) handleText(message *tgbotapi.Message) {
	pending, ok := tb.pendingMemes[message.Chat.ID]
	if !ok {
		return
	}

	// Parse text
	parts := strings.Split(message.Text, "|")
	topText := ""
	bottomText := ""

	if len(parts) > 0 {
		topText = strings.TrimSpace(parts[0])
	}
	if len(parts) > 1 {
		bottomText = strings.TrimSpace(parts[1])
	}

	// Send generating message
	statusMsg := tgbotapi.NewMessage(message.Chat.ID, "🎨 Generating meme...")
	tb.bot.Send(statusMsg)

	var outputPath string
	var err error

	if pending.Type == "template" {
		outputPath, err = tb.generator.CreateMeme(pending.Template, topText, bottomText)
	} else if pending.Type == "url" {
		outputPath, err = tb.generator.CreateMemeFromURL(pending.ImageURL, topText, bottomText)
	}

	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("❌ Error: %v", err))
		tb.bot.Send(msg)
		delete(tb.pendingMemes, message.Chat.ID)
		return
	}
	defer os.Remove(outputPath)

	photo := tgbotapi.NewPhoto(message.Chat.ID, tgbotapi.FilePath(outputPath))
	tb.bot.Send(photo)

	// Clear pending meme
	delete(tb.pendingMemes, message.Chat.ID)
}
