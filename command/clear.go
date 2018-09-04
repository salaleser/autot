package command

import (
	"github.com/sbstjn/hanu"
	"salaleser.ru/autot/util"
)

var Clear = func(conv hanu.ConversationInterface) {
	util.Files = map[string]string{}
	util.DeleteFileList()
	conv.Reply("Список файлов очищен.")
}
