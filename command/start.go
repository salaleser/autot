package command

import (
	"github.com/sbstjn/hanu"
	"salaleser.ru/autot/util"
)

// Start содержит функцию, которая запустит службу
var Start = func(conv hanu.ConversationInterface) {
	var text string
	switch util.Status {
	case util.StatusStopped:
		go util.Execute("start")
		text = "Запускаю…"
	case util.StatusStartPending:
		text = "Терпение, служба уже запускается!"
	case util.StatusStopPending:
		text = "Подождите, служба еще не остановлена!"
	case util.StatusRunning:
		text = "Служба уже запущена!"
	}
	conv.Reply(text)
}
