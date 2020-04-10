package slack

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/slack-go/slack"
)

var (
	ErrAlreadyClose   = fmt.Errorf("client already closed")
	ErrInvalidMessage = fmt.Errorf("invalid message")
	ErrInvalidCommand = fmt.Errorf("invalid command")
)

type BotClient struct {
	tokenID      string
	botID        string
	atBotID      string
	api          *slack.Client
	rtm          *slack.RTM
	stopClient   chan bool
	isClosed     chan struct{}
	isClosedSync *sync.Once
}

func NewBotClient(tokenID string) *BotClient {
	c := &BotClient{
		stopClient:   make(chan bool),
		isClosed:     make(chan struct{}),
		isClosedSync: &sync.Once{},
	}
	c.api = slack.New(tokenID)
	c.rtm = c.api.NewRTM()
	return c
}

func (c *BotClient) requestBotID() (string, error) {
	resp, err := c.api.AuthTest()
	if err != nil {
		return "", fmt.Errorf("fail to auth: %w", err)
	}
	return resp.UserID, nil
}

func (c *BotClient) setBotID(botID string) {
	c.botID = botID
	c.atBotID = fmt.Sprintf("<@%s>", botID)
}

func (c *BotClient) Init() error {
	botID, err := c.requestBotID()
	if err != nil {
		return fmt.Errorf("fail to request bot id: %w", err)
	}
	c.setBotID(botID)
	return nil
}

func (c *BotClient) Start() {
	go c.rtm.ManageConnection()
	for msg := range c.rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			text := ev.Msg.Text
			if !c.ValidateMessageEvent(text, ev.Channel) {
				continue
			}
			resp, err := c.ParseMessage(text, ev.Channel)
			if err != nil {
				log.Printf("[Error] message: %s, err: %v", text, err)
				continue
			}
			c.rtm.SendMessage(c.rtm.NewOutgoingMessage(resp, ev.Channel))
		}
	}
}

func (c *BotClient) close() {
	c.isClosedSync.Do(func() {
		close(c.isClosed)
	})
}

func (c *BotClient) Close() error {
	select {
	case c.stopClient <- true:
		c.close()
		if err := c.rtm.Disconnect(); err != nil {
			return fmt.Errorf("fail to disconnect rtm: %w", err)
		}
		return nil
	case <-c.isClosed:
		return ErrAlreadyClose
	}
}

func (c *BotClient) ValidateMessageEvent(text, channelID string) bool {
	return strings.HasPrefix(text, c.atBotID)
}

func (c *BotClient) ParseMessage(text, channelID string) (string, error) {
	msg := strings.Split(strings.TrimSpace(text), " ")[1:]
	if len(msg) == 0 {
		return "", ErrInvalidMessage
	}
	cmd := _BotCmd(msg[0])
	switch cmd {
	case _MsgIDCmd:
		return handleMsgID(msg[1])
	case _HelpCmd:
		return handleHelp()
	}
	return "", ErrInvalidCommand
}
