package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	DiscordToken    string
	TelegramToken   string
	BotPrefix       string
	TemplatesDir    string
	OutputDir       string
	ServerPort      string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	cfg := &Config{
		DiscordToken:  os.Getenv("DISCORD_TOKEN"),
		TelegramToken: os.Getenv("TELEGRAM_TOKEN"),
		BotPrefix:     getEnvOrDefault("BOT_PREFIX", "!"),
		TemplatesDir:  getEnvOrDefault("TEMPLATES_DIR", "meme_templates"),
		OutputDir:     getEnvOrDefault("OUTPUT_DIR", "output"),
		ServerPort:    getEnvOrDefault("SERVER_PORT", "8080"),
	}

	return cfg, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
