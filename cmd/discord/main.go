package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/mumugogoing/meme_bot/internal/config"
	"github.com/mumugogoing/meme_bot/pkg/meme"
)

type DiscordBot struct {
	session   *discordgo.Session
	generator *meme.Generator
	config    *config.Config
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if cfg.DiscordToken == "" {
		log.Fatal("DISCORD_TOKEN not found in environment variables")
	}

	generator, err := meme.NewGenerator(cfg.TemplatesDir, cfg.OutputDir)
	if err != nil {
		log.Fatalf("Failed to create meme generator: %v", err)
	}

	session, err := discordgo.New("Bot " + cfg.DiscordToken)
	if err != nil {
		log.Fatalf("Failed to create Discord session: %v", err)
	}

	bot := &DiscordBot{
		session:   session,
		generator: generator,
		config:    cfg,
	}

	session.AddHandler(bot.messageHandler)
	session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent

	if err := session.Open(); err != nil {
		log.Fatalf("Failed to open Discord session: %v", err)
	}
	defer session.Close()

	log.Println("Discord bot is now running. Press Ctrl+C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func (b *DiscordBot) messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages from the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Check if message starts with prefix
	if !strings.HasPrefix(m.Content, b.config.BotPrefix) {
		return
	}

	// Parse command
	content := strings.TrimPrefix(m.Content, b.config.BotPrefix)
	parts := strings.Fields(content)
	if len(parts) == 0 {
		return
	}

	command := parts[0]
	args := parts[1:]

	switch command {
	case "meme":
		b.handleMeme(s, m, args)
	case "templates":
		b.handleTemplates(s, m)
	case "memeurl":
		b.handleMemeURL(s, m, args)
	case "help_meme":
		b.handleHelp(s, m)
	}
}

func (b *DiscordBot) handleMeme(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) < 1 {
		s.ChannelMessageSend(m.ChannelID, "Usage: !meme <template_name> \"top text\" \"bottom text\"")
		return
	}

	template := args[0]
	topText := ""
	bottomText := ""

	// Parse quoted text
	content := strings.TrimPrefix(m.Content, b.config.BotPrefix+"meme "+template+" ")
	texts := parseQuotedStrings(content)
	if len(texts) > 0 {
		topText = texts[0]
	}
	if len(texts) > 1 {
		bottomText = texts[1]
	}

	outputPath, err := b.generator.CreateMeme(template, topText, bottomText)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("❌ Error: %v", err))
		return
	}
	defer os.Remove(outputPath)

	file, err := os.Open(outputPath)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "❌ Failed to open generated meme")
		return
	}
	defer file.Close()

	s.ChannelFileSend(m.ChannelID, outputPath, file)
}

func (b *DiscordBot) handleTemplates(s *discordgo.Session, m *discordgo.MessageCreate) {
	templates, err := b.generator.ListTemplates()
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "❌ Failed to list templates")
		return
	}

	if len(templates) == 0 {
		s.ChannelMessageSend(m.ChannelID, "No templates found. Add templates to the meme_templates directory.")
		return
	}

	message := "**Available Templates:**\n"
	for _, t := range templates {
		message += fmt.Sprintf("• %s\n", t)
	}
	s.ChannelMessageSend(m.ChannelID, message)
}

func (b *DiscordBot) handleMemeURL(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) < 1 {
		s.ChannelMessageSend(m.ChannelID, "Usage: !memeurl <image_url> \"top text\" \"bottom text\"")
		return
	}

	imageURL := args[0]
	topText := ""
	bottomText := ""

	// Parse quoted text
	content := strings.TrimPrefix(m.Content, b.config.BotPrefix+"memeurl "+imageURL+" ")
	texts := parseQuotedStrings(content)
	if len(texts) > 0 {
		topText = texts[0]
	}
	if len(texts) > 1 {
		bottomText = texts[1]
	}

	s.ChannelMessageSend(m.ChannelID, "🎨 Generating meme...")

	outputPath, err := b.generator.CreateMemeFromURL(imageURL, topText, bottomText)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("❌ Error: %v", err))
		return
	}
	defer os.Remove(outputPath)

	file, err := os.Open(outputPath)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "❌ Failed to open generated meme")
		return
	}
	defer file.Close()

	s.ChannelFileSend(m.ChannelID, outputPath, file)
}

func (b *DiscordBot) handleHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	help := `**Meme Bot Commands:**

` + "`" + b.config.BotPrefix + `meme <template> "top text" "bottom text"` + "`" + ` - Create a meme from a template
` + "`" + b.config.BotPrefix + `templates` + "`" + ` - List all available templates
` + "`" + b.config.BotPrefix + `memeurl <url> "top text" "bottom text"` + "`" + ` - Create a meme from an image URL
` + "`" + b.config.BotPrefix + `help_meme` + "`" + ` - Show this help message

**Examples:**
` + "`" + b.config.BotPrefix + `meme drake.jpg "studying" "watching memes"` + "`" + `
` + "`" + b.config.BotPrefix + `memeurl https://example.com/image.jpg "top" "bottom"` + "`"

	s.ChannelMessageSend(m.ChannelID, help)
}

func parseQuotedStrings(s string) []string {
	var result []string
	var current strings.Builder
	inQuotes := false

	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '"' {
			if inQuotes {
				result = append(result, current.String())
				current.Reset()
			}
			inQuotes = !inQuotes
		} else if inQuotes {
			current.WriteByte(c)
		}
	}

	return result
}
