package command

import (
	"fmt"
	"time"

	"github.com/nlopes/slack"
	"salaleser.ru/autot/poster"
	"salaleser.ru/autot/util"
)

// AutotHandler останавливает службу, пакует и копирует шаблоны и запускает службу
func AutotHandler(c *slack.Client, rtm *slack.RTM, ev *slack.MessageEvent, data []string) {
	if len(data) > 1 {
		ClearHandler(c, rtm, ev, data)
		AddHandler(c, rtm, ev, data)
	}

	StopHandler(c, rtm, ev, data)

	var count int
	const timeout = 180
	time.Sleep(time.Second)
	for util.Status != util.StatusStopped {
		time.Sleep(time.Second)
		if count > timeout {
			text := fmt.Sprintf("Превышено время ожидания (%d с). Служба останавливается слишком "+
				"долго. Попробуйте запустить отправку командой `!pull` вручную после остановки "+
				"службы, или перезапустите команду `!autot` немного позже", timeout)
			if util.Status == util.StatusStopPending { // Если остановка была отменена
				poster.PostError(ev.Channel, "Ошибка при попытке остановки службы!", text)
			}
			return
		}
		count++
	}

	PullHandler(c, rtm, ev, data)

	PingHandler(c, rtm, ev, data)

	StartHandler(c, rtm, ev, data)
}
