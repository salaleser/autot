package command

import (
	"github.com/sbstjn/hanu"
	"salaleser.ru/autot/util"
)

var VoteNegative = func(conv hanu.ConversationInterface) {
	if util.Status == 4 {
		util.OpStatus <- true
	} else {
		conv.Reply("```Службу не планировалось останавливать```")
	}
}
