package command

import (
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
	const clearCommandName = "`!clear`"
	const rmCommandName = "`!rm <номер_строки>`"
	conv.Reply("```%s```\n(%s — очистить список, %s — удалить файл)",
		text, clearCommandName, rmCommandName)
}
