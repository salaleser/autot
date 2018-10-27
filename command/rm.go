package command

import (
	"fmt"
	"strconv"

	"github.com/nlopes/slack"
	"salaleser.ru/autot/poster"
	"salaleser.ru/autot/util"
)

// RmHandler удаляет элемент из списка отправляемых файлов по его номеру (ключу в мапе)
func RmHandler(c *slack.Client, rtm *slack.RTM, ev *slack.MessageEvent, data []string) {
	if len(data) < 2 {
		poster.PostError(ev.Channel, "Ошибка!", "Не указан номер")
		return
	}
	key := data[1]
	filename := util.Files[key]

	if _, err := strconv.Atoi(key); err != nil {
		poster.PostError(ev.Channel, "Ошибка!", fmt.Sprintf("%q не является числом!", key))
		return
	}

	if len(filename) == 0 {
		poster.PostError(ev.Channel, "Ошибка!",
			fmt.Sprintf("Файла с номером %s нет в списке!", key))
		return
	}

	delete(util.Files, key)
	util.UpdateBackupFile()

	poster.Post(ev.Channel, "", fmt.Sprintf("Файл `%s` удален из списка", filename), "")
}
