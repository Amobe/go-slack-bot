package slack

import (
	"fmt"
	"strconv"
	"time"
)

type _BotCmd string

const (
	_HelpCmd  _BotCmd = "help"
	_MsgIDCmd _BotCmd = "msgid"
)

var helpMsg = `
Usage: @bot_name COMMAND

Commands:
  msgid    Decode the message id
  help     Show help information

`

func handleHelp() string {
	return helpMsg
}

func handleMsgID(midStr string) string {
	mid, err := strconv.ParseUint(midStr, 10, 64)
	if err != nil {
		return fmt.Errorf("fail to parse msg id: %v", err).Error()
	}
	offset := uint32(mid >> 32)
	id := uint32(mid)
	date := time.Unix(int64(offset), 0).Format(time.RFC3339)
	return fmt.Sprintf("offset: %d, id: %d, date: %s", offset, id, date)
}
