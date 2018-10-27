package poster

import (
	"github.com/nlopes/slack"
	"salaleser.ru/autot/gui"
	"salaleser.ru/autot/util"
)

var params = slack.PostMessageParameters{}

func PostError(channel string, title string, text string) {
	attachment := slack.Attachment{
		Color: gui.Red,
		Title: title,
		Text:  text,
	}
	params.Attachments = []slack.Attachment{attachment}
	util.API.PostMessage(channel, "", params)
}

func PostWarning(channel string, title string, text string, footer string) {
	attachment := slack.Attachment{
		Color:  gui.Orange,
		Title:  title,
		Text:   text,
		Footer: footer,
	}
	params.Attachments = []slack.Attachment{attachment}
	util.API.PostMessage(channel, "", params)
}

func Post(channel string, title string, text string, footer string) {
	attachment := slack.Attachment{
		Color:  gui.Green,
		Title:  title,
		Text:   text,
		Footer: footer,
	}
	params.Attachments = []slack.Attachment{attachment}
	util.API.PostMessage(channel, "", params)
}
