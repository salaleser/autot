package command

import (
	"strings"

	"github.com/nlopes/slack"
	"salaleser.ru/autot/util"
)

// AliasesHandler содержит функцию, которая отображает известные алиасы шаблонов.
// Эти алиасы содержатся в файле "aliases.list"
func AliasesHandler(c *slack.Client, rtm *slack.RTM, ev *slack.MessageEvent, data []string) {
	var text string
	columnWidth := 24
	for filename, alias := range util.Aliases {
		spaces := columnWidth - len(filename)
		if spaces < 1 {
			spaces = 1
		}
		text += filename + strings.Repeat(" ", spaces) + alias + "\n"
	}
	params := slack.PostMessageParameters{}
	attachment := slack.Attachment{
		Title: "Список алиасов шаблонов:",
		Text:  "```" + text + "```",
	}
	params.Attachments = []slack.Attachment{attachment}
	util.API.PostMessage(ev.Channel, "", params)
}
