package bot

import (
	"errors"
	"fmt"
	"miniEdward/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	prefix = "instagram.com"

	newPrefix = "ddinstagram.com"

	replaceAmount = 1
)

var instaReel = regexp.MustCompile(`https://www.instagram.com/reel(.*?)`)

type Bot struct {
	httpClient http.TelegramAPIClient
	readFlag   bool
	stopSignal chan bool

	logger *logrus.Logger
}

func NewBot(httpClient http.TelegramAPIClient, logger *logrus.Logger) *Bot {
	return &Bot{
		httpClient: httpClient,
		readFlag:   true,
	}
}

func (b *Bot) Start() {
	offset := int64(1)

	for {
		select {
		case <-b.stopSignal:
			b.logger.Info("Exiting update checking")
			return
		default:
			updates, err := b.GetUpdates(offset)
			if err != nil {
				b.logger.Error(err) // send error to chan
			}

			for _, update := range updates {
				if update.Message == nil {
					continue
				}
				text := update.Message.Text

				if instaReel.MatchString(text) && b.readFlag {
					resp := replaceString(text)
					b.logger.Info("response:", resp)

					response := fmt.Sprintf("%s", resp)
					err := b.SendMessage(update.Message.Chat.ID, update.Message.ID, response)
					if err != nil {
						b.logger.Error(errors.New(fmt.Sprintf("failed to send message: %s", err.Error())))
					}

				}
				offset = update.ID + 1
			}
		}
	}

}

func (b *Bot) Stop() {
	b.stopSignal <- true
}

func (b *Bot) SetRead(ctx *gin.Context) {
	b.logger.Warn("starting reading")
	b.readFlag = true
}

func (b *Bot) UnsetRead(ctx *gin.Context) {
	b.logger.Warn("stopping reading")
	b.readFlag = false
}

func replaceString(msg string) string {
	return strings.Replace(msg, prefix, newPrefix, replaceAmount)
}
