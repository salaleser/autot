package command

import (
	"github.com/sbstjn/hanu"
	"salaleser.ru/autot/util"
)

var Ver = func(conv hanu.ConversationInterface) {
	message := "Участвовать в разработке можно на гитхабе https://github.com/salaleser/autot"
	conv.Reply("Текущая версия бота: *%s*\n%s", util.Ver, message)
}
