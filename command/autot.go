package command

import (
	"time"

	"github.com/sbstjn/hanu"
	"salaleser.ru/autot/util"
)

// Autot сореджит функцию, которая остановит службу, запакует и скопирует шаблоны и запустит службу
var Autot = func(conv hanu.ConversationInterface) {
	Stop(conv)
	time.Sleep(time.Second)
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
