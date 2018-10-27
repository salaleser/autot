package command

import (
	"github.com/nlopes/slack"
	"salaleser.ru/autot/util"
)

// VoteNegativetHandler отменяет остановку службы
func VoteNegativetHandler(c *slack.Client, rtm *slack.RTM, ev *slack.MessageEvent, data []string) {
	if util.Status == util.StatusRunning {
		util.OpStatus <- true
	}
}
