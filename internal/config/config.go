package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}

type Config struct {
	TgBotToken string
}

func FromEnv() *Config {
	return &Config{
		TgBotToken: os.Getenv("GRISHA_BOT_TOKEN"),
	}
}
