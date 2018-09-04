package command

import (
	"time"

	"github.com/sbstjn/hanu"
	"salaleser.ru/autot/util"
)

// Stop сореджит функцию, которая остановит службу
var Stop = func(conv hanu.ConversationInterface) {
	util.Conv = conv
	switch util.Status {
	case util.StatusStopped:
		conv.Reply("Служба уже остановлена!")
	case util.StatusStartPending:
		conv.Reply("Подождите, служба еще не запущена!")
	case util.StatusStopPending:
		conv.Reply("Проявите терпение, служба уже останавливается!")
	case util.StatusRunning:
		cd := util.Countdown
		util.OpStatus = make(chan bool)
		util.Beep(util.Sounds[5])
		conv.Reply("*ВНИМАНИЕ!*\nСлужба будет остановлена через %d секунд!\n(`-` — отмена)", cd)
		for i := cd; i > 0; i-- {
			select {
			case <-util.OpStatus:
				conv.Reply("```Остановка службы отменена```")
				return
			default:
				time.Sleep(time.Second)
				if i%3 == 0 {
					util.Beep(util.Sounds[6])
				}
			}
		}
		conv.Reply("Останавливаю…")
		util.Execute("stop")
	}
}
