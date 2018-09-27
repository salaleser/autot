package command

import (
	"fmt"
	"log"
	"strconv"

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

	// Здесь нужно не количество элементов, а номер последнего ключа (пропуская несуществующие)
	var lk int
	for k := range util.Files {
		ik, err := strconv.Atoi(k)
		if err != nil {
			log.Println(err)
		}
		if ik > lk {
			lk = ik
		}
	}

	text = "Список отправляемых файлов:\n"
	for i := 1; i <= lk; i++ { // цикл для сортировки
		for key, value := range util.Files {
			if key == strconv.Itoa(i) {
				alias, ok := util.Aliases[value]
				if ok {
					alias = " («" + alias + "»)"
				}
				text += key + ". " + value + alias + "\n"
			}
		}
	}

	// params := slack.PostMessageParameters{}

	const cmdClear = "`!clear`"
	const cmdRm = "`!rm <номер_строки>`"
	footer := fmt.Sprintf("(%s — очистить список, %s — удалить файл)", cmdClear, cmdRm)
	// attachment := slack.Attachment{
	// 	Color:  gui.Green,
	// 	Text:   text,
	// 	Footer: footer,
	// }
	// params.Attachments = []slack.Attachment{attachment}
	// params.AsUser = true

	// util.API.PostMessage("", "", params)
	conv.Reply("```%s```\n%s", text, footer)
}
