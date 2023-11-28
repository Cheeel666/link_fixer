package main

import (
	"errors"
	"fmt"
	"miniEdward/bot"
	"miniEdward/cfg"
	"miniEdward/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
)

var instaReel = regexp.MustCompile(`https://www.instagram.com/reel(.*?)`)
var log = logrus.New()

const (
	oldPrefix     = "https://www.instagram.com/reel"
	newPrefix     = "https://www.ddinstagram.com/reel"
	replaceAmount = 1
)

func main() {
	cfg, err := cfg.NewCfg()
	if err != nil {
		panic(err)
	}

	log.Out = os.Stdout

	file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.Out = file
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	httpClient, err := http.NewStdTelegramAPIClient(cfg.Token)
	if err != nil {
		log.Fatal(err)
	}

	bot := bot.NewBot(httpClient)

	offset := int64(1)
	gcCount := 0

	errChan := make(chan error)
	abortChan := make(chan bool, 2)
	go abortOnError(errChan, abortChan)
	gcCount += 1

	log.Info("starting message checking")
	gcCount += 1
	go func() {
		for {
			select {
			case <-abortChan:
				log.Info("Exiting update checking")
				return
			default:
				updates, err := bot.GetUpdates(offset)
				if err != nil {
					errChan <- err // send error to chan
				}

				for _, update := range updates {
					log.Info("update: ", spew.Sdump(update.Message))
					if update.Message == nil {
						continue
					}
					text := update.Message.Text

					if instaReel.MatchString(text) {
						resp := replaceString(text)
						log.Info("message matched: user ", update.Message.From.Username, update.Message.From.ID, "chat -", update.Message.Chat.ID)
						log.Info("message: ", text, "\n", "response:", resp)

						response := fmt.Sprintf("@%s прислал:\n%s", update.Message.From.Username, resp)
						err := bot.SendMessage(update.Message.Chat.ID, response)
						if err != nil {
							errChan <- errors.New(fmt.Sprintf("failed to send message: %s", err.Error()))
						}
						err = bot.DeleteMessage(update.Message.Chat.ID, update.Message.ID)
						if err != nil {
							errChan <- errors.New(fmt.Sprintf("failed to delete message: %s", err.Error()))
						}

					}
					offset = update.ID + 1
				}
			}
		}

	}()

	var stopChan = make(chan os.Signal, 2)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-stopChan // wait for SIGINT

	for i := 0; i < gcCount; i++ {
		abortChan <- true
	}

	log.Info("Waiting till processes die")
	time.Sleep(time.Second * 1)
}

func replaceString(msg string) string {
	return strings.Replace(msg, oldPrefix, newPrefix, replaceAmount)
}

func abortOnError(errChan <-chan error, abort <-chan bool) {
	for {
		select {
		case err := <-errChan:
			log.Error(err)

		case <-abort:
			log.Warn("Exiting error checking")
			return
		default:
			time.Sleep(time.Second * 1)
		}
	}

}
