package command

import (
	"strconv"
	"time"

	"github.com/nlopes/slack"

	"salaleser.ru/autot/gui"
	"salaleser.ru/autot/util"
)

// StopHandler сореджит функцию, которая остановит службу
func StopHandler(c *slack.Client, rtm *slack.RTM, ev *slack.MessageEvent, data []string) {
	switch util.Status {
	case util.StatusStopped:
		warningParams := slack.PostMessageParameters{}
		attachment := slack.Attachment{
			Color: gui.Orange,
			Title: "Служба уже остановлена!",
		}
		warningParams.Attachments = []slack.Attachment{attachment}
		util.API.PostMessage(ev.Channel, "", warningParams)
	case util.StatusStartPending:
		warningParams := slack.PostMessageParameters{}
		attachment := slack.Attachment{
			Color: gui.Orange,
			Title: "Подождите, служба еще не запущена!",
		}
		warningParams.Attachments = []slack.Attachment{attachment}
		util.API.PostMessage(ev.Channel, "", warningParams)
	case util.StatusStopPending:
		warningParams := slack.PostMessageParameters{}
		attachment := slack.Attachment{
			Color: gui.Orange,
			Title: "Проявите терпение, служба уже останавливается!",
		}
		warningParams.Attachments = []slack.Attachment{attachment}
		util.API.PostMessage(ev.Channel, "", warningParams)
	case util.StatusRunning:
		cd := util.Countdown
		util.OpStatus = make(chan bool)
		util.Beep(util.Sounds[5])

		alertChannel, err := util.GetAlertChannel()
		if err == nil {
			params := slack.PostMessageParameters{}
			attachment := slack.Attachment{
				Color: gui.Red,
				Text: "*ВНИМАНИЕ!*\nСлужба будет остановлена через " +
					strconv.Itoa(cd) + " секунд!",
				Footer: "(`-` — отмена)",
			}
			params.Attachments = []slack.Attachment{attachment}
			params.AsUser = true
			util.API.PostMessage(alertChannel.ID, "", params)
		}

		for i := cd; i > 0; i-- {
			select {
			case <-util.OpStatus:
				warningParams := slack.PostMessageParameters{}
				attachment := slack.Attachment{
					Color: gui.Orange,
					Title: "Остановка службы отменена",
				}
				warningParams.Attachments = []slack.Attachment{attachment}
				util.API.PostMessage(ev.Channel, "", warningParams)
				return
			default:
				time.Sleep(time.Second)
				if i%3 == 0 {
					util.Beep(util.Sounds[6])
				}
			}
		}
		attachment := slack.Attachment{
			Color: gui.Green,
			Title: "Останавливаю…",
		}
		params.Attachments = []slack.Attachment{attachment}
		util.Execute("stop")
	}
}
