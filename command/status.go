package command

import (
	"fmt"

	"github.com/nlopes/slack"
	"github.com/sbstjn/hanu"
	"salaleser.ru/autot/util"
)

// Status содержит функцию, которая отобразит список отправляемых файлов
var Status = func(conv hanu.ConversationInterface) {
	var text string

	if len(util.Files) == 0 {
		text = "Список отправляемых файлов пуст"
		const addCommandName = "`!add <имя_файла>`"
		conv.Reply("```%s```\n(%s — добавить файл)", text, addCommandName)
		return
	}

	text = "Список отправляемых файлов:\n"
	for key, value := range util.Files {
		alias, ok := util.Aliases[value]
		if ok {
			alias = " («" + alias + "»)"
		}
		text += key + ". " + value + alias + "\n"
	}

	params := slack.PostMessageParameters{}

	const cmdClear = "`!clear`"
	const cmdRm = "`!rm <номер_строки>`"
	footer := fmt.Sprintf("(%s — очистить список, %s — удалить файл)", cmdClear, cmdRm)
	attachment := slack.Attachment{
		Color:  "e7ff47",
		Text:   text,
		Footer: footer,
	}
	params.Attachments = []slack.Attachment{attachment}
	params.AsUser = true

	util.Api.PostMessage("", "", params)
	conv.Reply("```%s```", text)
}
