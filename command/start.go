package command

import (
	"github.com/nlopes/slack"
	"salaleser.ru/autot/poster"
	"salaleser.ru/autot/util"
)

// StartHandler запускает службу
func StartHandler(c *slack.Client, rtm *slack.RTM, ev *slack.MessageEvent, data []string) {
	switch util.Status {
	case util.StatusStopped:
		go util.Execute("start")
	case util.StatusStartPending:
		poster.PostWarning(ev.Channel, "", "Терпение, служба уже запускается!", "")
	case util.StatusStopPending:
		poster.PostWarning(ev.Channel, "", "Подождите, служба еще не остановлена!", "")
	case util.StatusRunning:
		poster.PostWarning(ev.Channel, "", "Служба уже запущена!", "")
	}
}
