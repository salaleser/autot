package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/sbstjn/hanu"
)

const (
	statusRunning = "4  RUNNING"
	process       = "sc"
	server        = "\\\\Kmisdevserv"
	service       = "IBM Domino Server (DDOMINODATA)"
)

var (
	running     bool
	lastStarted int64
	lastStopped int64
)

type config struct {
	Token   string `json:"token"`
	Version string `json:"version"`
}

func loopQuery() {
	for {
		out := runCommand("query")
		running = strings.Contains(out, statusRunning)
		if running {
			if lastStarted < lastStopped {
				lastStarted = time.Now().UnixNano()
				fmt.Println("Служба запущена.")
			}
		} else {
			if lastStarted > lastStopped {
				lastStopped = time.Now().UnixNano()
				fmt.Println("Служба остановлена!")
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func runCommand(s string) string {
	out, err := exec.Command(process, server, s, service).Output()
	if err != nil {
		log.Fatal(err)
	}

	return string(out)
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
		fmt.Print("Reconnecting in... ")
		for i := 30; i > 0; i-- {
			time.Sleep(time.Second)
			fmt.Print(i, " ")
		}
		fmt.Println("Reconnecting now...")
		goto Reconnect
	}

	fmt.Println("Connected!")

	go loopQuery()

	// slack.Command("shout <word>", func(conv hanu.ConversationInterface) {
	// 	str, _ := conv.String("word")
	// 	conv.Reply(strings.ToUpper(str))
	// })

	slack.Command("query", func(conv hanu.ConversationInterface) {
		var msg string
		if running {
			rndSource := rand.NewSource(time.Now().UnixNano())
			rndRand := rand.New(rndSource)
			rndInt := rndRand.Intn(100)
			if rndInt > 90 {
				msg = "лотусист спит, служба идет..."
			} else {
				msg = "служба бежит..."
			}
		} else {
			msg = "служба лежит!"
		}
		conv.Reply("%s", msg)
	})

	slack.Command("stop", func(conv hanu.ConversationInterface) {
		runCommand("stop")
		conv.Reply("%s", "останавливаю...")
	})

	slack.Command("start", func(conv hanu.ConversationInterface) {
		runCommand("start")
		conv.Reply("%s", "запускаю...")
	})

	slack.Command("version", func(conv hanu.ConversationInterface) {
		conv.Reply("`%s`", response.Version)
	})

	slack.Command("?", func(conv hanu.ConversationInterface) {
		conv.Reply("`query` — проверка состояния службы;\n" +
			"`start` — запустить службу;\n" +
			"`stop` — остановить службу.")
	})

	slack.Listen()
}
