package command

import (
	"strings"

	"github.com/sbstjn/hanu"
	"salaleser.ru/autot/util"
)

var Config = func(conv hanu.ConversationInterface) {
	text := "Текущие настройки:\n"
	columnWidth := 20
	for key, value := range util.Config {
		spaces := columnWidth - len(key)
		if spaces < 1 {
			spaces = 1
		}
		text += key + strings.Repeat(" ", spaces) + value + "\n"
	}
	conv.Reply("```%s```\nИзменить настройки можно командой `!config <key> <value>`", text)
}

var ConfigReplaceValue = func(conv hanu.ConversationInterface) {
	key, err := conv.String("key")
	if err != nil {
		conv.Reply("```Ошибка!\n%s```", err)
		return
	}

	value, err := conv.String("value")
	if err != nil {
		conv.Reply("```Ошибка!\n%s```", err)
		return
	}

	_, ok := util.Config[key]
	if !ok {
		conv.Reply("```Нет такого ключа```")
		return
	}
	util.Config[key] = value
	util.ReloadConfig()
	conv.Reply("Значение ключа `%s` изменено на `%s`", key, value)
}

var ConfigReload = func(conv hanu.ConversationInterface) {
	util.ReadFileIntoMap("config", util.Config)
	util.ReloadConfig()
}
