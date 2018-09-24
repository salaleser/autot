package command

import (
	"log"

	"github.com/nlopes/slack"
	"github.com/sbstjn/hanu"
	"salaleser.ru/autot/util"
)

// Ping отправляет начальнику сообщение о файле с шаблонами
var Ping = func(conv hanu.ConversationInterface) {
	user, err := util.Api.GetUserByEmail("pravednik@rkmail.ru")
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

	util.Api.PostMessage(user.ID, "", params)
}
