package command

import (
	"github.com/nlopes/slack"
	"salaleser.ru/autot/gui"
	"salaleser.ru/autot/util"
)

// VoteNegativetHandler отменяет остановку службы
func VoteNegativetHandler(c *slack.Client, rtm *slack.RTM, ev *slack.MessageEvent, data []string) {
	if util.Status == util.StatusRunning {
		util.OpStatus <- true
	} else {
		errorParams := slack.PostMessageParameters{}
		attachment := slack.Attachment{
			Color: gui.Red,
			Title: "Службу не планировалось останавливать",
		}
		errorParams.Attachments = []slack.Attachment{attachment}
		util.API.PostMessage(ev.Channel, "", errorParams)
	}
}
