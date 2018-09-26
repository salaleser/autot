package command

import (
	"strconv"

	"github.com/sbstjn/hanu"
	"salaleser.ru/autot/util"
)

// Rm удаляет элемент из списка отправляемых файлов по его номеру (ключу в мапе)
var Rm = func(conv hanu.ConversationInterface) {
	conv.Reply("Команда работает нестабильно и временно отключена")
	return

	key, err := conv.String("номер")
	if err != nil {
		conv.Reply("```Ошибка!\n%s```", err)
		return
	}
	filename := util.Files[key]

	if _, err := strconv.Atoi(key); err != nil {
		conv.Reply("%q не является числом!", key)
		return
	}

	if len(filename) == 0 {
		conv.Reply("Файла с номером %s нет в списке!", key)
		return
	}

	delete(util.Files, key)
	util.UpdateBackupFile()

	conv.Reply("Файл `%s` удален из списка", filename)
}
