package slack

func StartSlackBot(tokenID string, botID string, channelIDs ...string) {
	c := NewClient(tokenID, botID, channelIDs...)
	c.Start()
}
