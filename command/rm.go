package command

import (
	"github.com/sbstjn/hanu"
	"salaleser.ru/autot/util"
)

var Rm = func(conv hanu.ConversationInterface) {
	key, err := conv.String("номер")
	if err != nil {
		conv.Reply("```Ошибка!\n%s```", err)
		return
	}
	filename := util.Files[key]

	delete(util.Files, key)

	util.SaveFileList()

	conv.Reply("Файл `%s` удален из списка", filename)
}
