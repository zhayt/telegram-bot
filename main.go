package main

import (
	"context"
	"flag"
	tgClient "github.com/zhayt/read-adviser-bot/clients/telegram"
	"github.com/zhayt/read-adviser-bot/consumer/event-consumer"
	"github.com/zhayt/read-adviser-bot/events/telegram"
	"github.com/zhayt/read-adviser-bot/storage/sqlite"
	"log"
)

const (
	tgBotHost         = "api.telegram.org"
	sqliteStoragePath = "date/sqlite/storage.db"
	batchSize         = 100
)

func main() {
	//s := files.New(storagePath)
	s, err := sqlite.New(sqliteStoragePath)
	if err != nil {
		log.Fatalf("can't connect to storage: %s", err.Error())
	}

	if err := s.Init(context.TODO()); err != nil {
		log.Fatalf("can't init storage: %s", err.Error())
	}

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		s,
	)

	log.Print("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)
	if err := consumer.Start(); err != nil {
		log.Fatal(err)
	}
}

// mustToken use for takes token value as an argument from terminal
// if no-token is specified, the program terminates
func mustToken() string {
	token := flag.String(
		"tg-token",
		"",
		"token for access to telegram bot",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}
