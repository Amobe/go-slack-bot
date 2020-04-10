package slack

import "fmt"

func StartSlackBot(tokenID string) error {
	c := NewBotClient(tokenID)
	if err := c.Init(); err != nil {
		return fmt.Errorf("fail to initial new bot client: %w", err)
	}
	c.Start()
	return nil
}
