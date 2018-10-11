package command

import (
	"github.com/nlopes/slack"
)

// HelpHandler показывает подсказки к командам
func HelpHandler(c *slack.Client, rtm *slack.RTM, ev *slack.MessageEvent, data []string) {
	msg := "Помощи «нигде нет». Просто слов нет. Найдем слова – сделаем помощь." +
		" Вы держитесь здесь, вам всего доброго, хорошего настроения и здоровья."
	rtm.SendMessage(rtm.NewOutgoingMessage("Помощь на подходе (в разработке)\n"+msg, ev.Channel))
}
