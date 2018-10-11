package command

import (
	"log"
	"strconv"

	"github.com/nlopes/slack"
	"salaleser.ru/autot/gui"
	"salaleser.ru/autot/util"
)

// StatusHandler показывает список отправляемых файлов
func StatusHandler(c *slack.Client, rtm *slack.RTM, ev *slack.MessageEvent, data []string) {
	if len(util.Files) == 0 {
		params := slack.PostMessageParameters{}
		attachment := slack.Attachment{
			Color:  gui.Orange,
			Title:  "Список отправляемых файлов пуст",
			Footer: "*!add <имена_файлов_через_пробелы>* — добавить файл",
		}
		params.Attachments = []slack.Attachment{attachment}
		params.AsUser = true
		util.API.PostMessage(ev.Channel, "", params)
		return
	}

	// Здесь нужно не количество элементов, а номер последнего ключа (пропуская несуществующие)
	var lk int
	for k := range util.Files {
		ik, err := strconv.Atoi(k)
		if err != nil {
			log.Println(err)
		}
		if ik > lk {
			lk = ik
		}
	}

	var text string
	for i := 1; i <= lk; i++ { // цикл для сортировки
		for key, value := range util.Files {
			if key == strconv.Itoa(i) {
				alias, ok := util.Aliases[value]
				if ok {
					alias = " («" + alias + "»)"
				}
				text += key + ". " + value + alias + "\n"
			}
		}
	}

	params := slack.PostMessageParameters{}
	attachment := slack.Attachment{
		Color:  gui.Green,
		Title:  "Список отправляемых файлов:",
		Text:   "```" + text + "```",
		Footer: "`!clear` — очистить список, `!rm <номер_строки>` — удалить файл",
	}
	params.Attachments = []slack.Attachment{attachment}
	params.AsUser = true
	util.API.PostMessage(ev.Channel, "", params)
}
