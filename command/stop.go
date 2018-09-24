package command

import (
	"strconv"
	"time"

	"github.com/nlopes/slack"

	"github.com/sbstjn/hanu"
	"salaleser.ru/autot/gui"
	"salaleser.ru/autot/util"
)

// Stop сореджит функцию, которая остановит службу
var Stop = func(conv hanu.ConversationInterface) {
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

		alertChannel, err := util.GetAlertChannel()
		if err == nil {
			params := slack.PostMessageParameters{}
			attachment := slack.Attachment{
				Color: gui.Red,
				Text: "*ВНИМАНИЕ!*\nСлужба будет остановлена через " + strconv.Itoa(cd) +
					" секунд!",
				Footer: "(`-` — отмена)",
			}
			params.Attachments = []slack.Attachment{attachment}
			params.AsUser = true
			util.Api.PostMessage(alertChannel.ID, "", params)
		}

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
