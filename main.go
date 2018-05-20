package main

import (
	"log"
	"strings"

	"github.com/sbstjn/hanu"
)

func main() {
	slack, err := hanu.New("xoxb-360859824437-368020146983-YBCW6cEWKP0O1GTUNS2KDepr")
	if err != nil {
		log.Fatal(err)
	}

	Version := "0.0.1"

	slack.Command("shout <word>", func(conv hanu.ConversationInterface) {
		str, _ := conv.String("word")
		conv.Reply(strings.ToUpper(str))
	})

	slack.Command("whisper <word>", func(conv hanu.ConversationInterface) {
		str, _ := conv.String("word")
		conv.Reply(strings.ToLower(str))
	})

	slack.Command("version", func(conv hanu.ConversationInterface) {
		conv.Reply("Thanks for asking! I'm running `%s`", Version)
	})

	slack.Listen()
}
