package command

import (
	"github.com/sbstjn/hanu"
	"salaleser.ru/autot/util"
)

// VoteNegative содержит функцию, которая отменяет остановку службы
var VoteNegative = func(conv hanu.ConversationInterface) {
	if util.Status == util.StatusRunning {
		util.OpStatus <- true
	} else {
		conv.Reply("```Службу не планировалось останавливать```")
	}
}
