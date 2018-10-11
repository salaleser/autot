package command

import (
	"fmt"
	"time"

	"github.com/nlopes/slack"
	"salaleser.ru/autot/gui"
	"salaleser.ru/autot/util"
)

// AutotHandler сореджит функцию, которая остановит службу, запакует и скопирует шаблоны и запустит службу
func AutotHandler(c *slack.Client, rtm *slack.RTM, ev *slack.MessageEvent, data []string) {
	StopHandler(c, rtm, ev, data)

	var count int
	const timeout = 180
	time.Sleep(time.Second)
	for util.Status != util.StatusStopped {
		time.Sleep(time.Second)
		if count > timeout {
			errorParams := slack.PostMessageParameters{}
			attachment := slack.Attachment{
				Color: gui.Red,
				Title: fmt.Sprintf("Превышено время ожидания (%d с). Служба останавливается слишком долго. "+
					"Попробуйте запустить отправку командой `!pull` вручную после остановки службы, "+
					"или перезапустите команду `!autot` немного позже", timeout),
			}
			errorParams.Attachments = []slack.Attachment{attachment}
			util.API.PostMessage(ev.Channel, "", errorParams)
			return
		}
		count++
	}

	PullHandler(c, rtm, ev, data)

	PingHandler(c, rtm, ev, data)

	StartHandler(c, rtm, ev, data)
}
