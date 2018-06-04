// ...
// Тут будет описание
// ...

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/sbstjn/hanu"
)

const (
	statusStopped      = "1  STOPPED"
	statusStartPending = "2  START_PENDING"
	statusStopPending  = "3  STOP_PENDING"
	statusRunning      = "4  RUNNING"
)

var (
	status      int
	lastStarted int64
	lastStopped int64
)

type config struct {
	Token   string `json:"token"`
	Version string `json:"version"`
}

func runCommand(s string) string {
	out, err := exec.Command("sc", "\\\\Kmisdevserv", s, "IBM Domino Server (DDOMINODATA)").Output()
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

	log.Println("Connecting...")
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

	log.Println("Connected!")

	// Вечный цикл проверки статуса.
	// 2 раза в секунду опрашивает состояние службы утилитой sc.exe, сравнивает ее
	// вывод с константами и перезаписывает переменную status
	go func() {
		for {
			out := runCommand("query")
			if strings.Contains(out, statusRunning) {
				status = 4
				if lastStarted < lastStopped {
					lastStarted = time.Now().UnixNano()
					fmt.Println("Служба запущена.")
				}
			} else if strings.Contains(out, statusStopped) {
				// FIXME: может произойти ситуация, когда служба служба успеет остановиться между
				// итерациями этого цикла, тогда бот не сможет оповестить об остановке службы
				status = 1
				if lastStarted > lastStopped {
					lastStopped = time.Now().UnixNano()
					fmt.Println("Служба останавливается!")
				}
			} else if strings.Contains(out, statusStartPending) {
				status = 2
			} else if strings.Contains(out, statusStopPending) {
				status = 3
			}
			time.Sleep(500 * time.Millisecond)
		}
	}()

	// slack.Command("shout <word>", func(conv hanu.ConversationInterface) {
	// 	str, _ := conv.String("word")
	// 	conv.Reply(strings.ToUpper(str))
	// })

	slack.Command("query", func(conv hanu.ConversationInterface) {
		rndSource := rand.NewSource(time.Now().UnixNano())
		rndRand := rand.New(rndSource)
		rndInt := rndRand.Intn(100)
		var msg string
		switch status {
		case 4:
			if rndInt > 90 {
				msg = "лотусист спит, служба идет..."
			} else {
				msg = "служба работает."
			}
		case 1:
			if rndInt > 90 {
				msg = "служба лежит!"
			} else {
				msg = "служба остановлена!"
			}
		case 2:
			msg = "служба запускается..."
		case 3:
			msg = "служба останавливается..."
		}
		msg += "\nlastStarted = " + strconv.FormatInt(lastStarted, 10) +
			" / lastStopped = " + strconv.FormatInt(lastStopped, 10)
		conv.Reply("%s", msg)
	})

	slack.Command("stop", func(conv hanu.ConversationInterface) {
		if status == 4 {
			runCommand("stop")
			conv.Reply("%s", "останавливаю...")
		} else {
			conv.Reply("%s", "сервис уже останавливается или остановлен.")
		}
	})

	slack.Command("start", func(conv hanu.ConversationInterface) {
		if status != 4 {
			runCommand("start")
			conv.Reply("%s", "запускаю...")
		} else {
			conv.Reply("%s", "сервис уже запущен.")
		}
	})

	slack.Command("version", func(conv hanu.ConversationInterface) {
		conv.Reply("`%s`", response.Version)
	})

	slack.Command("help", func(conv hanu.ConversationInterface) {
		conv.Reply("`query` — проверка состояния службы;\n" +
			"`start` — запустить службу;\n" +
			"`stop` — остановить службу.")
	})

	slack.Listen()
}
