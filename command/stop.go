package command

import (
	"time"

	"github.com/sbstjn/hanu"
	"salaleser.ru/autot/util"
)

// Stop stops the service
var Stop = func(conv hanu.ConversationInterface) {
	util.Conv = conv

	var text string
	switch util.Status {
	case 1:
		text = "Служба уже остановлена!"
	case 2:
		text = "Подождите, служба еще не запущена!"
	case 3:
		text = "Проявите терпение, служба уже останавливается!"
	case 4:
		countdown := util.Countdown
		util.OpStatus = make(chan bool)
		conv.Reply("*ВНИМАНИЕ!*\nСлужба будет остановлена через %d с!\n(`-` — отмена)", countdown)
		for i := countdown; i > 0; i-- {
			select {
			case <-util.OpStatus:
				conv.Reply("```Остановка службы отменена```")
				return
			default:
				time.Sleep(time.Second)
			}
		}
		go util.Execute("stop")
		text = "Останавливаю…"
	}
	conv.Reply(text)
}
