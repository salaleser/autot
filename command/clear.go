package command

import (
	"github.com/nlopes/slack"
	"salaleser.ru/autot/gui"
	"salaleser.ru/autot/util"
)

// ClearHandler очищает список отправляемых файлов и удаляет файл-бэкап с диска
func ClearHandler(c *slack.Client, rtm *slack.RTM, ev *slack.MessageEvent, data []string) {
	util.Files = map[string]string{}
	util.UpdateBackupFile()

	params := slack.PostMessageParameters{}
	attachment := slack.Attachment{
		Color: gui.Green,
		Text:  "Список файлов очищен",
	}
	params.Attachments = []slack.Attachment{attachment}
	util.API.PostMessage(ev.Channel, "", params)
}
