package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/sbstjn/hanu"
)

type config struct {
	Token   string `json:"token"`
	Version string `json:"version"`
}

func main() {
	f, err := os.Open("config.json")
	if err != nil {
		log.Fatal("Конфигурационный файл не найден!")
	}

	body := make([]byte, 68) //FIXME: hardcode!

	n1, err := f.Read(body)
	if err != nil {
		fmt.Printf("%d bytes: %s\n", n1, string(body))
	}

	var response config
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("Failed to unmarshal JSON: %s", body)
	}

Reconnect:
	slack, err := hanu.New(response.Token)
	if err != nil {
		log.Println(err)
		fmt.Print("Reconnecting... ")
		for i := 30; i > 0; i-- {
			time.Sleep(time.Second)
			fmt.Print(i)
			fmt.Print(" ")
		}
		fmt.Println()
		// log.Fatal(err)
		goto Reconnect
	}

	fmt.Println("Connected!")

	slack.Command("shout <word>", func(conv hanu.ConversationInterface) {
		str, _ := conv.String("word")
		conv.Reply(strings.ToUpper(str))
	})

	slack.Command("<word>", func(conv hanu.ConversationInterface) {
		str, _ := conv.String("word")
		fmt.Println(str)
	})

	slack.Command("version", func(conv hanu.ConversationInterface) {
		conv.Reply("Thanks for asking! I'm running `%s`", response.Version)
	})

	slack.Listen()
}
