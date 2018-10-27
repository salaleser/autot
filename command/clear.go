package command

import (
	"github.com/nlopes/slack"
	"salaleser.ru/autot/poster"
	"salaleser.ru/autot/util"
)

// ClearHandler очищает список отправляемых файлов и удаляет файл-бэкап с диска
func ClearHandler(c *slack.Client, rtm *slack.RTM, ev *slack.MessageEvent, data []string) {
	util.Files = map[string]string{}
	util.UpdateBackupFile()

	poster.Post(ev.Channel, "", "Список файлов очищен", "")
}
