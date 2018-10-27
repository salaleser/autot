package command

import (
	"fmt"
	"strings"

	"github.com/nlopes/slack"
	"salaleser.ru/autot/poster"
	"salaleser.ru/autot/util"
)

// ConfigHandler меняет значение указанного ключа в памяти программы (не на диске)
func ConfigHandler(c *slack.Client, rtm *slack.RTM, ev *slack.MessageEvent, data []string) {
	if len(data) < 3 {
		var text string
		columnWidth := 20
		for key, value := range util.Config {
			spaces := columnWidth - len(key)
			if spaces < 1 {
				spaces = 1
			}
			text += key + strings.Repeat(" ", spaces) + value + "\n"
		}

		poster.Post(ev.Channel, "Текущие настройки:", "```"+text+"```", "`!config <key> <value>` "+
			"— заменить значение ключ, `!config reload` — загрузить настройки из файла")
		return
	}

	key := data[1]
	value := data[2]

	_, ok := util.Config[key]
	if !ok {
		poster.PostError(ev.Channel, "Ошибка!", "Нет такого ключа")
		return
	}
	util.Config[key] = value
	util.ReloadConfig()

	poster.Post(ev.Channel, "", fmt.Sprintf("Значение ключа `%s` изменено на `%s`", key, value), "")
}
