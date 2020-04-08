package main

import (
	"os"

	"github.com/amobe/go-slack-bot/slack"
)

func main() {
	os.Exit(mainWithCode())
}

func mainWithCode() int {
	var (
		tokenID    = ""
		botID      = ""
		channelIDs = []string{""}
	)

	slack.StartSlackBot(tokenID, botID, channelIDs...)

	return 0
}
