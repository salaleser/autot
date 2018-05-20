package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
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

	body := make([]byte, 86) //FIXME: hardcode!

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

	slack.Command("q", func(conv hanu.ConversationInterface) {
		out, err := exec.Command("sc").Output()
		if err != nil {
			log.Println(err)
		}
		conv.Reply("%s", out)
	})

	slack.Command("version", func(conv hanu.ConversationInterface) {
		conv.Reply("Thanks for asking! I'm running `%s`", response.Version)
	})

	slack.Command("help", func(conv hanu.ConversationInterface) {
		conv.Reply("`q` — тест", response.Version)
	})

	slack.Listen()
}
