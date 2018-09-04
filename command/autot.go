package command

import (
	"time"

	"github.com/sbstjn/hanu"
	"salaleser.ru/autot/util"
)

// Autot содержит функцию, которая поочередно запустит команды на остановку службы, копирования и
// упаковку шаблонов из списка отправляемых в специальную папку и запуска службы
var Autot = func(conv hanu.ConversationInterface) {
	Stop(conv)
	var count int
	for util.Status != util.StatusStopped {
		time.Sleep(time.Second)
		if count > 60 {
			conv.Reply("Превышено время ожидания. Операция отменена")
			return
		}
		count++
	}
	Pull(conv)
	Start(conv)
}
