package command

import (
	"github.com/nlopes/slack"
	"salaleser.ru/autot/util"
)

// ConfigReloadHandler считывает настройки из файла и перезаписывает ими текущие в памяти
func ConfigReloadHandler(c *slack.Client, rtm *slack.RTM, ev *slack.MessageEvent, data []string) {
	util.ReadFileIntoMap("config", util.Config)
	util.ReloadConfig()
}
