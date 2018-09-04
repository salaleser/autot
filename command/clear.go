package command

import (
	"github.com/sbstjn/hanu"
	"salaleser.ru/autot/util"
)

// Clear содержит функцию, которая очистит список отправляемых файлов и удалит файл-бэкап с диска
var Clear = func(conv hanu.ConversationInterface) {
	util.Files = map[string]string{}
	util.DeleteFileList()
	conv.Reply("Список файлов очищен.")
}
