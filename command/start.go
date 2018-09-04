package command

import (
	"github.com/sbstjn/hanu"
	"salaleser.ru/autot/util"
)

var Start = func(conv hanu.ConversationInterface) {
	var text string
	switch util.Status {
	case 1:
		go util.Execute("start")
		text = "Запускаю…"
	case 2:
		text = "Терпение, служба уже запускается!"
	case 3:
		text = "Подождите, служба еще не остановлена!"
	case 4:
		text = "Служба уже запущена!"
	}
	conv.Reply(text)
}
