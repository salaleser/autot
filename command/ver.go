package command

import (
	"fmt"

	"github.com/nlopes/slack"
	"salaleser.ru/autot/gui"
	"salaleser.ru/autot/util"
)

// VerHandler отображает краткую справку о боте
func VerHandler(c *slack.Client, rtm *slack.RTM, ev *slack.MessageEvent, data []string) {
	params := slack.PostMessageParameters{}
	attachment := slack.Attachment{
		Color:     gui.Green,
		ImageURL:  "https://github.com/salaleser/autot", // FIXME
		ThumbURL:  "https://github.com/salaleser/autot", // FIXME
		Title:     fmt.Sprintf("Текущая версия бота: %s", util.Ver),
		TitleLink: "https://github.com/salaleser/autot",
		Text:      "Участвовать в разработке можно на гитхабе https://github.com/salaleser/autot",
	}
	params.Attachments = []slack.Attachment{attachment}
	params.AsUser = true
	util.API.PostMessage(ev.Channel, "", params)
}
