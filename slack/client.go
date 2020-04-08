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

type Client struct {
	tokenID       string
	botID         string
	validChannels []string
	rtm           *slack.RTM
	stopClient    chan bool
	isClosed      chan struct{}
	isClosedSync  *sync.Once
}

func NewClient(tokenID, botID string, channelIDs ...string) *Client {
	c := &Client{
		tokenID:      tokenID,
		botID:        botID,
		stopClient:   make(chan bool),
		isClosed:     make(chan struct{}),
		isClosedSync: &sync.Once{},
	}
	for _, cid := range channelIDs {
		c.validChannels = append(c.validChannels, cid)
	}

	api := slack.New(c.tokenID)
	c.rtm = api.NewRTM()
	return c
}

func (c *Client) Start() {
	go c.rtm.ManageConnection()
	for msg := range c.rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			if !c.ValidateMessageEvent(ev.Msg.Text, ev.Channel, ev.BotID) {
				continue
			}
			resp, err := c.ParseMessage(ev.Msg.Text, ev.Channel)
			if err != nil {
				log.Printf("[Error] Fail to handle message: %v", err)
				continue
			}
			c.rtm.SendMessage(c.rtm.NewOutgoingMessage(resp, ev.Channel))
		}
	}
}

func (c *Client) close() {
	c.isClosedSync.Do(func() {
		close(c.isClosed)
	})
}

func (c *Client) Close() error {
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

func (c *Client) ValidateMessageEvent(text, channelID, botID string) bool {
	if !c.validateChannel(channelID) {
		log.Printf("invlid channel, cid: %s, msg: %s", channelID, text)
		return false
	}
	if !strings.HasPrefix(text, botID) {
		log.Printf("not to bot, bid: %s, msg: %s", botID, text)
		return false
	}
	return true
}

func (c *Client) validateChannel(channelID string) bool {
	for _, cid := range c.validChannels {
		if cid == channelID {
			return true
		}
	}
	return false
}

func (c *Client) ParseMessage(text, channelID string) (string, error) {
	msg := strings.Split(strings.TrimSpace(text), " ")
	if len(msg) == 0 {
		return "", ErrInvalidMessage
	}
	cmd := _BotCmd(msg[0])
	switch cmd {
	case _MsgIDCmd:
		return handleMsgID(msg[1]), nil
	case _HelpCmd:
		return handleHelp(), nil
	}
	return "", ErrInvalidCommand
}
