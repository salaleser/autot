package command

import (
	"github.com/sbstjn/hanu"
	"salaleser.ru/autot/util"
)

var Query = func(conv hanu.ConversationInterface) {
	isTimeToJoke := util.RandomNumber(100) > 95
	var text string
	switch util.Status {
	case 1:
		if isTimeToJoke {
			text = "Служба лежит!"
		} else {
			text = "Служба остановлена!"
		}
	case 2:
		text = "Служба запускается…"
	case 3:
		text = "Служба останавливается…"
	case 4:
		if isTimeToJoke {
			text = "Лотусист спит, служба идет…"
		} else {
			text = "Служба работает"
		}
	}
	conv.Reply(text)
}
