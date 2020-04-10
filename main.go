package main

import (
	"log"
	"os"

	"github.com/amobe/go-slack-bot/slack"
)

func main() {
	os.Exit(mainWithCode())
}

func mainWithCode() int {
	var (
		tokenID = "REPLACE_BY_YOUR_TOKEN"
	)

	err := slack.StartSlackBot(tokenID)
	log.Printf("[Error] fail to start slack bot: %v", err)

	return 0
}
