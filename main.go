package main

import (
	"flag"
	"github.com/zhayt/read-adviser-bot/clients/telegram"
	"log"
)

const (
	tgBotHost = "api.telegram.org"
)

func main() {
	// client use for connect telegram api, token for identified user
	tgClient := telegram.New(tgBotHost, mustToken())

	// fetcher = fetcher.New()         -> fetcher for gets events

	// processor = processor.New()        -> processor for handler events

	// consumer.Start(fetcher, processor)   -> consumer gets events and handlers its
}

// mustToken use for takes token value as an argument from terminal
// if no-token is specified, the program terminates
func mustToken() string {
	token := flag.String(
		"token-bot-token",
		"",
		"token for access to telegram bot",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}
