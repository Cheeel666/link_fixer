package bot

import "miniEdward/http"

type Bot struct {
	httpClient http.TelegramAPIClient
}

func NewBot(httpClient http.TelegramAPIClient) *Bot {
	return &Bot{httpClient: httpClient}
}
