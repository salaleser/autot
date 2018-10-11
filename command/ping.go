package command

import (
	"log"

	"github.com/nlopes/slack"
	"salaleser.ru/autot/util"
)

// PingHandler отправляет начальнику сообщение о файле с шаблонами
func PingHandler(c *slack.Client, rtm *slack.RTM, ev *slack.MessageEvent, data []string) {
	user, err := util.API.GetUserByEmail("pravednik@rkmail.ru")
	if err != nil {
		log.Printf("%s\n", err)
		return
	}

	params := slack.PostMessageParameters{}
	attachment := slack.Attachment{
		Pretext: "Шаблоны упакованы в архив здесь:",
		Text:    "`" + util.ArcFullName + "`",
	}
	params.Attachments = []slack.Attachment{attachment}

	util.API.PostMessage(user.ID, "", params)
}
