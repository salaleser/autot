package command

import (
	"fmt"
	"strings"

	"github.com/nlopes/slack"
	"salaleser.ru/autot/gui"
	"salaleser.ru/autot/util"
)

// ConfigHandler меняет значение указанного ключа в памяти программы (не на диске)
func ConfigHandler(c *slack.Client, rtm *slack.RTM, ev *slack.MessageEvent, data []string) {
	a := strings.Split(ev.Msg.Text, " ")

	if len(a) < 3 {
		var text string
		columnWidth := 20
		for key, value := range util.Config {
			spaces := columnWidth - len(key)
			if spaces < 1 {
				spaces = 1
			}
			text += key + strings.Repeat(" ", spaces) + value + "\n"
		}

		params := slack.PostMessageParameters{}
		attachment := slack.Attachment{
			Color: gui.Green,
			Title: "Текущие настройки:",
			Text:  "```" + text + "```",
			Footer: "*!config <key> <value>* — изменить настройки," +
				"*!config reload* — загрузить настройки из файла",
		}
		params.Attachments = []slack.Attachment{attachment}
		util.API.PostMessage(ev.Channel, "", params)
		return
	}

	key := a[1]
	value := a[2]

	_, ok := util.Config[key]
	if !ok {
		errorParams := slack.PostMessageParameters{}
		attachment := slack.Attachment{
			Color: gui.Red,
			Title: "Ошибка!",
			Text:  "Нет такого ключа",
		}
		errorParams.Attachments = []slack.Attachment{attachment}
		util.API.PostMessage(ev.Channel, "", errorParams)
		return
	}
	util.Config[key] = value
	util.ReloadConfig()

	params := slack.PostMessageParameters{}
	attachment := slack.Attachment{
		Color: gui.Green,
		Text:  fmt.Sprintf("Значение ключа `%s` изменено на `%s`", key, value),
	}
	params.Attachments = []slack.Attachment{attachment}
	util.API.PostMessage(ev.Channel, "", params)
}
