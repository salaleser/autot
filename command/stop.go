package command

import (
	"strconv"
	"time"

	"github.com/nlopes/slack"

	"salaleser.ru/autot/poster"
	"salaleser.ru/autot/util"
)

// StopHandler останавливает службу
func StopHandler(c *slack.Client, rtm *slack.RTM, ev *slack.MessageEvent, data []string) {
	switch util.Status {
	case util.StatusStopped:
		poster.PostWarning(ev.Channel, "Служба уже остановлена!", "", "")
	case util.StatusStartPending:
		poster.PostWarning(ev.Channel, "Подождите, служба еще не запущена!", "", "")
	case util.StatusStopPending:
		poster.PostWarning(ev.Channel, "Проявите терпение, служба уже останавливается!", "", "")
	case util.StatusRunning:
		cd := util.Countdown
		util.OpStatus = make(chan bool)
		util.Beep(util.Sounds[5])

		alertChannel, err := util.GetAlertChannel()
		if err == nil {
			poster.PostWarning(alertChannel.ID, "ВНИМАНИЕ!",
				"Служба будет остановлена через "+strconv.Itoa(cd)+" секунд!", "`-` — отмена")
		}

		for i := cd; i > 0; i-- {
			select {
			case <-util.OpStatus:
				poster.Post(ev.Channel, "Остановка службы отменена", "", "")
				return
			default:
				time.Sleep(time.Second)
				if i%3 == 0 {
					util.Beep(util.Sounds[6])
				}
			}
		}
		util.Execute("stop")
	}
}
