package command

import (
	"time"

	"github.com/sbstjn/hanu"
	"salaleser.ru/autot/util"
)

// Autot сореджит функцию, которая остановит службу, запакует и скопирует шаблоны и запустит службу
var Autot = func(conv hanu.ConversationInterface) {
	Stop(conv)

	var count int
	const timeout = 180
	time.Sleep(time.Second)
	for util.Status != util.StatusStopped {
		time.Sleep(time.Second)
		if count > timeout {
			conv.Reply("Превышено время ожидания (%d с). Служба останавливается слишком долго. "+
				"Попробуйте запустить отправку командой `!pull` вручную после остановки службы, "+
				"или перезапустите команду `!autot` немного позже", timeout)
			return
		}
		count++
	}
	Pull(conv)
	Ping(conv)
	Start(conv)
}
