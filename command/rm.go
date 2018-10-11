package command

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/nlopes/slack"
	"salaleser.ru/autot/gui"
	"salaleser.ru/autot/util"
)

// RmHandler удаляет элемент из списка отправляемых файлов по его номеру (ключу в мапе)
func RmHandler(c *slack.Client, rtm *slack.RTM, ev *slack.MessageEvent, data []string) {
	a := strings.Split(ev.Msg.Text, " ")

	if len(a) < 2 {
		errorParams := slack.PostMessageParameters{}
		attachment := slack.Attachment{
			Color: gui.Red,
			Text:  "Не указан номер",
		}
		errorParams.Attachments = []slack.Attachment{attachment}
		util.API.PostMessage(ev.Channel, "", errorParams)
		return
	}
	key := a[1]
	filename := util.Files[key]

	if _, err := strconv.Atoi(key); err != nil {
		errorParams := slack.PostMessageParameters{}
		attachment := slack.Attachment{
			Color: gui.Red,
			Title: "Ошибка!",
			Text:  fmt.Sprintf("%q не является числом!", key),
		}
		errorParams.Attachments = []slack.Attachment{attachment}
		util.API.PostMessage(ev.Channel, "", errorParams)
		return
	}

	if len(filename) == 0 {
		errorParams := slack.PostMessageParameters{}
		attachment := slack.Attachment{
			Color: gui.Red,
			Title: "Ошибка!",
			Text:  fmt.Sprintf("Файла с номером %s нет в списке!", key),
		}
		errorParams.Attachments = []slack.Attachment{attachment}
		util.API.PostMessage(ev.Channel, "", errorParams)
		return
	}

	delete(util.Files, key)
	util.UpdateBackupFile()

	params := slack.PostMessageParameters{}
	attachment := slack.Attachment{
		Color: gui.Blue,
		Text:  fmt.Sprintf("Файл `%s` удален из списка", filename),
	}
	params.Attachments = []slack.Attachment{attachment}
	util.API.PostMessage(ev.Channel, "", params)
}
