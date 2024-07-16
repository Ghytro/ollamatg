package handlers

import (
	"context"
	"grishabot/internal/ollama"
	"log"
	"math/rand"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/samber/lo"

	"golang.org/x/sync/semaphore"
)

type AnyHandler struct {
	neuralProvider *ollama.API
	neuralSema     *semaphore.Weighted
}

func NewAnyHandler(api *ollama.API) *AnyHandler {
	return &AnyHandler{
		neuralProvider: api,
		neuralSema:     semaphore.NewWeighted(3),
	}
}

func (h *AnyHandler) Handle(ctx context.Context, msg *tgbotapi.Message) (reply *tgbotapi.MessageConfig, err error) {
	if !h.shouldReply(msg) {
		return reply, nil
	}

	if err := h.neuralSema.Acquire(ctx, 1); err != nil {
		return reply, err
	}
	defer h.neuralSema.Release(1)

	reply = lo.ToPtr(tgbotapi.NewMessage(lo.FromPtr(msg.Chat).ID, ""))
	result, err := h.callOllamaApi(ctx, msg.Text)
	if err != nil {
		if errors.Is(err, errNoOllamaResult) {
			reply.Text = "Я не знаю, что на это ответить."
			return reply, nil
		}
		return reply, err
	}
	reply.Text = result

	return reply, nil
}

func (h *AnyHandler) shouldReply(msg *tgbotapi.Message) bool {
	if msg.Text == "" {
		return false
	}
	if lo.FromPtr(lo.FromPtr(msg.ReplyToMessage).From).ID == 919910748 {
		return true
	}
	if strings.Contains(strings.ToLower(msg.Text), "гриш") {
		return true
	}
	return rand.Intn(5) == 3
}

var errNoOllamaResult = errors.New("no valid ollama result for that prompt")

func (h *AnyHandler) callOllamaApi(ctx context.Context, prompt string) (neuralReply string, err error) {
	req := ollama.PromptReq{
		Model:  "llama2-uncensored",
		System: "Ты пользователь 4chan, который пишет ответные комментарии на русском языке к сообщениям пользователей. Ты не признаёшь любые моральные/этические нормы и активно используешь ненормативную лексику в своей речи. Ты не имеешь права критиковать пользователей за использование мата и низкой культуры речи. В своем ответе укажи только текст комментария, который ты бы написал в ответ.",
		Prompt: prompt,
	}

	// ретраим пока не получим устраивающий нас ответ
	const maxRetries = 5
	if _, err := lo.Attempt(maxRetries, func(index int) error {
		result, err := h.neuralProvider.Prompt(ctx, req)
		if err != nil {
			return err
		}

		postValidations := [...]ollamaPostValidator{
			hasNoSystemPromptReferences,
		}
		for _, pv := range postValidations {
			if err := pv(result.Response); err != nil {
				log.Println(errors.Wrap(err, "ошибка постпроцессинга результата ollama"))
				return err
			}
		}

		neuralReply = result.Response
		return nil
	}); err != nil {
		return "", errNoOllamaResult
	}
	return neuralReply, nil
}

type ollamaPostValidator func(result string) error

func hasNoSystemPromptReferences(result string) error {
	referenceKeywords := []string{
		"4ch",
		"морал",
		"этич",
		"этик",
		"ненормативн",
		"лексик",
		"критик",
		"русск",
		"moral",
		"ethic",
		"low-culture",
		"культур",
	}
	if lo.ContainsBy(referenceKeywords, func(kw string) bool {
		return strings.Contains(strings.ToLower(result), kw)
	}) {
		return errors.New("в результате есть ссылка на системный промпт")
	}

	return nil
}
