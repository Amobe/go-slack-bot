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

// Eaxmple:
// @bot_name help
// Commands:
//   msgid    Decode the message id
//   help     Show help information
var helpReplyFmt = "```\n" +
	"Commands:\n" +
	"  msgid    Decode the message id\n" +
	"  help     Show help information\n" +
	"```\n"

// Example:
// @bot_name msgid 5577006791947779410
// |      offset |          id |                      date |
// |  1298498081 |   134020434 |      2011-02-23T21:54:41Z |
var msgIDReplyFmt = "```\n" +
	"|      offset |          id |                 date |\n" +
	"| %11d | %11d | %20s |\n" +
	"```\n"

func handleHelp() (string, error) {
	return helpReplyFmt, nil
}

func handleMsgID(midStr string) (string, error) {
	mid, err := strconv.ParseUint(midStr, 10, 64)
	if err != nil {
		return "", fmt.Errorf("fail to parse msg id: %w", err)
	}
	offset := uint32(mid >> 32)
	id := uint32(mid)
	date := time.Unix(int64(offset), 0).UTC().Format(time.RFC3339)
	return fmt.Sprintf(msgIDReplyFmt, offset, id, date), nil
}
