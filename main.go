package main

import (
	"miniEdward/bot"
	"miniEdward/cfg"
	"miniEdward/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

func main() {
	cfg, err := cfg.NewCfg()
	if err != nil {
		panic(err)
	}

	log := logrus.New()
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

	bot := bot.NewBot(httpClient, log)

	// server := srv.NewSrv(bot)
	log.Info("starting message checking")
	go bot.Start()
	defer bot.Stop()

	// TODO: add ability to kill this f thing remotely
	// go server.Start()

	var stopChan = make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-stopChan // wait for SIGINT
}
