package main

import (
	"context"
	"grishabot/internal/config"
	"grishabot/internal/handlers"
	"grishabot/internal/ollama"
	"grishabot/internal/tenor"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	config := config.FromEnv()
	bot, err := tgbotapi.NewBotAPI(config.TgBotToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	ollamaApi := ollama.NewApi("http://localhost:11434")
	tenorApi := tenor.NewApi()
	hands := handlers.NewAnyHandler(ollamaApi, tenorApi)

	for update := range updates {
		if update.Message != nil {
			go func() {
				response, err := hands.Handle(context.Background(), update.Message)
				if err != nil {
					log.Println(err)
				}
				if response == nil {
					return
				}

				response = setReplyToMsgId(response, update.Message.MessageID)

				bot.Send(response)
			}()
		}
	}
}

// setReplyToMsgId оч противоречивый метод. С одной стороны с логики
// должен возвращаться какой то своего внутреннего типа респонс. С другой стороны
// шо то хня шо это хня писать либо ифы с type-assert либо интерфейс+обертка с сеттером
func setReplyToMsgId(chattable tgbotapi.Chattable, replyMsgId int) tgbotapi.Chattable {
	if textReply, ok := chattable.(tgbotapi.MessageConfig); ok {
		textReply.ReplyToMessageID = replyMsgId
		return textReply
	}
	if animationReply, ok := chattable.(tgbotapi.AnimationConfig); ok {
		animationReply.ReplyToMessageID = replyMsgId
		return animationReply
	}
	return chattable
}
