package command

import (
	"log"
	"strconv"

	"github.com/nlopes/slack"
	"salaleser.ru/autot/poster"
	"salaleser.ru/autot/util"
)

// StatusHandler показывает список отправляемых файлов
func StatusHandler(c *slack.Client, rtm *slack.RTM, ev *slack.MessageEvent, data []string) {
	if len(util.Files) == 0 {
		poster.PostWarning(ev.Channel, "Список отправляемых файлов пуст", "",
			"`!add <имена_файлов_через_пробелы>` — добавить один или несколько файлов")
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

	poster.Post(ev.Channel, "Список отправляемых файлов:",
		"```"+text+"```", "`!clear` — очистить список, `!rm <номер_строки>` — удалить файл")
}
