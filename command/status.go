package command

import (
	"github.com/sbstjn/hanu"
	"salaleser.ru/autot/util"
)

var Status = func(conv hanu.ConversationInterface) {
	var text string

	if len(util.Files) == 0 {
		text = "Список отправляемых файлов пуст"
		conv.Reply("```%s```\n(`!add <имя_файла>` — добавить файл)", text)
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
	conv.Reply("```%s```\n(`!clear` — очистить список, `!rm <номер_строки>` — удалить файл)", text)
}
