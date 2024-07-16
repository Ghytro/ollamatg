package main

import (
	"context"
	"grishabot/internal/handlers"
	"grishabot/internal/ollama"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("919910748:AAEYlyLMpywKlgnixu0yBsX7B_vamg5vlDA")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	ollamaApi := ollama.NewApi("http://localhost:11434")
	hands := handlers.NewAnyHandler(ollamaApi)

	for update := range updates {
		if update.Message != nil { // If we got a message
			go func() {
				response, err := hands.Handle(context.Background(), update.Message)
				if err != nil {
					log.Println(err)
				}
				if response == nil {
					return
				}

				// response := "я пидор"
				response.ReplyToMessageID = update.Message.MessageID

				bot.Send(response)
			}()
		}
	}
}
