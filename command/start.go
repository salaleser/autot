package command

import (
	"github.com/nlopes/slack"
	"salaleser.ru/autot/gui"
	"salaleser.ru/autot/util"
)

// StartHandler содержит функцию, которая запустит службу
func StartHandler(c *slack.Client, rtm *slack.RTM, ev *slack.MessageEvent, data []string) {
	var text string
	switch util.Status {
	case util.StatusStopped:
		go util.Execute("start")
		text = "Запускаю…"
	case util.StatusStartPending:
		text = "Терпение, служба уже запускается!"
	case util.StatusStopPending:
		text = "Подождите, служба еще не остановлена!"
	case util.StatusRunning:
		text = "Служба уже запущена!"
	}
	params := slack.PostMessageParameters{}
	attachment := slack.Attachment{
		Color: gui.Green,
		Title: text,
	}
	params.Attachments = []slack.Attachment{attachment}
	params.AsUser = true
	util.API.PostMessage(ev.Channel, "", params)
}
