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
	prefix = "www.instagram.com"

	newPrefix = "www.ddinstagram.com"

	replaceAmount = 1
)

var instaReel = regexp.MustCompile(`www.instagram.com/(.*?)/`)

type Bot struct {
	httpClient http.TelegramAPIClient
	readFlag   bool
	stopSignal chan bool

	offset int64

	logger *logrus.Logger
}

func NewBot(httpClient http.TelegramAPIClient, logger *logrus.Logger) *Bot {
	return &Bot{
		offset:     1,
		httpClient: httpClient,
		readFlag:   true,
		logger:     logger,
		stopSignal: make(chan bool),
	}
}

// recalculateOffset to prevent spam with bot start
func (b *Bot) recalculateOffset() {
	updates, err := b.GetUpdates(b.offset)
	if err != nil {
		b.logger.Error(err)
	}
	for _, update := range updates {
		b.offset = update.ID + 1
	}
}

func (b *Bot) Start() {
	b.recalculateOffset()

	for {
		select {
		case <-b.stopSignal:
			b.logger.Info("Exiting update checking")
			return
		default:
			updates, err := b.GetUpdates(b.offset)
			if err != nil {
				b.logger.Error(err)
			}
			for _, update := range updates {
				if update.Message == nil || update.Message.Chat == nil {
					b.offset = update.ID + 1
					continue
				}
				text := update.Message.Text

				if instaReel.MatchString(text) && b.readFlag {
					newStr := replaceString(text)
					b.logger.Info("response:", newStr)

					response := fmt.Sprintf("%s", newStr)

					err := b.SendMessage(update.Message.Chat.ID, update.Message.ID, response)
					if err != nil {
						b.logger.Error(errors.New(fmt.Sprintf("failed to send message: %s", err.Error())))
					}

				}
				b.offset = update.ID + 1
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
