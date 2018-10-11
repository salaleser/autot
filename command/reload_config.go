package command

import (
	"github.com/nlopes/slack"
	"salaleser.ru/autot/util"
)

// ConfigReloadHandler содержит функцию, которая считает настройки из файла и перезапишет ими текущие
func ConfigReloadHandler(c *slack.Client, rtm *slack.RTM, ev *slack.MessageEvent, data []string) {
	util.ReadFileIntoMap("config", util.Config)
	util.ReloadConfig()
}
