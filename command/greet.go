package command

import (
	"fmt"

	"github.com/nlopes/slack"
)

// GreetHandler приветствует в ответ
func GreetHandler(c *slack.Client, rtm *slack.RTM, ev *slack.MessageEvent, data []string) {
	rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("Привет, <@%s>", ev.User), ev.Channel))
}
