package command

import (
	"github.com/nlopes/slack"
	"salaleser.ru/autot/poster"
	"salaleser.ru/autot/util"
)

// QueryHandler сообщает о текущем состоянии службы
func QueryHandler(c *slack.Client, rtm *slack.RTM, ev *slack.MessageEvent, data []string) {
	isTimeToJoke := util.RandomNumber(100) > 95
	var text string
	switch util.Status {
	case util.StatusStopped:
		if isTimeToJoke {
			text = "Служба лежит!"
		} else {
			text = "Служба остановлена!"
		}
	case util.StatusStartPending:
		text = "Служба запускается…"
	case util.StatusStopPending:
		text = "Служба останавливается…"
	case util.StatusRunning:
		if isTimeToJoke {
			text = "Лотусист спит, служба идет…"
		} else {
			text = "Служба работает"
		}
	}

	poster.Post(ev.Channel, "", text, "")
}
