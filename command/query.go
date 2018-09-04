package command

import (
	"github.com/sbstjn/hanu"
	"salaleser.ru/autot/util"
)

// Query содержит функцию, которая вернет сообщение о текущем состоянии службы в текущий канал
var Query = func(conv hanu.ConversationInterface) {
	isTimeToJoke := util.RandomNumber(100) > 95
	var text string
	switch util.Status {
	case util.StatusStopped:
		if isTimeToJoke {
			text = "Служба лежит!"
		} else {
			text = "Служба остановлена!"
		}
	case util.StatusStartPending:
		text = "Служба запускается…"
	case util.StatusStopPending:
		text = "Служба останавливается…"
	case util.StatusRunning:
		if isTimeToJoke {
			text = "Лотусист спит, служба идет…"
		} else {
			text = "Служба работает"
		}
	}
	conv.Reply(text)
}
