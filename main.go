// ...
// Тут будет описание
// ...

package main

import (
	"log"
	"os/exec"
	"strings"
	"time"
)

const (
	statusStopped      = "1  STOPPED"
	statusStartPending = "2  START_PENDING"
	statusStopPending  = "3  STOP_PENDING"
	statusRunning      = "4  RUNNING"
	countdown          = 5 // количество секунда до остановки службы
)

var (
	status           int8 // 1 -- остановлена, 2 -- запускается, 3 -- останавливается, 4 -- запущена
	lastStartPending int64
	lastStopPending  int64
	lastStopped      int64
	lastStarted      int64

	// hook hanu.ConversationInterface
)

type config struct {
	Token   string `json:"token"`
	Version string `json:"version"`
}

func execute(s string) string {
	out, err := exec.Command("sc", "\\\\Kmisdevserv", s, "IBM Domino Server (DDOMINODATA)").Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(out)
}

func discordbot() {
	// discord, err := discordgo.New("Bot " + "authentication token")
	// if err != nil {
	// 	log.Fatal(err)
	// }

}

func main() {
	// f, err := os.Open("config.json")
	// if err != nil {
	// 	log.Fatal("Конфигурационный файл не найден!")
	// }

	// body := make([]byte, 86) //FIXME: hardcode!

	// n1, err := f.Read(body)
	// if err != nil {
	// 	fmt.Printf("%d bytes: %s\n", n1, string(body))
	// }

	// var response config
	// err = json.Unmarshal(body, &response)
	// if err != nil {
	// 	fmt.Printf("Failed to unmarshal JSON: %s", body)
	// }

	// В связи с тем, что slack в России (в США тоже подозреваю проблемы, но во Франции всегда
	// стабильное подключение) подключается через раз, то пока используем discord
	// discordbot()
	// os.Exit(545)

	// Reconnect:
	// log.Println("Connecting...")
	// 	slack, err := hanu.New(response.Token)
	// 	if err != nil {
	// 		log.Println(err)
	// 		fmt.Print("Reconnecting in... ")
	// 		for i := 30; i > 0; i-- {
	// 			time.Sleep(time.Second)
	// 			fmt.Print(i, " ")
	// 		}
	// 		fmt.Println("0")
	// 		goto Reconnect
	// 	}

	// log.Println("Connected!")

	// opStatus := make(chan bool) // канал для отмены остановки

	// Вечный цикл проверки статуса.
	// 2 раза в секунду опрашивает состояние службы утилитой sc.exe, сравнивает ее
	// вывод с константами и перезаписывает переменную status
	func() {
		for {
			out := execute("query")
			if strings.Contains(out, statusRunning) {
				status = 4
				if lastStopped >= lastStarted {
					lastStarted = time.Now().Unix()
					log.Println("Служба запущена.")
					// hook.Reply("`%s`", "Служба запущена.")
				}
			} else if strings.Contains(out, statusStopped) {
				status = 1
				if lastStartPending >= lastStopped {
					lastStopped = time.Now().Unix()
					log.Println("Служба остановлена.")
					// hook.Reply("`%s`", "Служба остановлена.")
				}
			} else if strings.Contains(out, statusStartPending) {
				status = 2
				if lastStopPending >= lastStartPending {
					lastStartPending = time.Now().Unix()
					log.Println("Служба запускается.")
					// hook.Reply("`%s`", "Служба запускается.")
				}
			} else if strings.Contains(out, statusStopPending) {
				status = 3
				if lastStarted >= lastStopPending {
					lastStopPending = time.Now().Unix()
					log.Println("Служба останавливается!")
					// hook.Reply("`%s`", "Служба останавливается!")
				}
			}
			time.Sleep(500 * time.Millisecond)
		}
	}()

	// slack.Command("shout <word>", func(conv hanu.ConversationInterface) {
	// 	str, _ := conv.String("word")
	// 	conv.Reply(strings.ToUpper(str))
	// })

	// var cmd hanu.Command
	// cmd = hanu.NewCommand("query", "проверяет состояние службы",
	// 	func(conv hanu.ConversationInterface) {
	// 		isTimeToJoke := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100) > 95
	// 		var text string
	// 		switch status {
	// 		case 1:
	// 			if isTimeToJoke {
	// 				text = "служба лежит!"
	// 			} else {
	// 				text = "служба остановлена!"
	// 			}
	// 		case 2:
	// 			text = "служба запускается..."
	// 		case 3:
	// 			text = "служба останавливается..."
	// 		case 4:
	// 			if isTimeToJoke {
	// 				text = "лотусист спит, служба идет..."
	// 			} else {
	// 				text = "служба работает."
	// 			}
	// 		}
	// 		conv.Reply("%s", text)
	// 	})
	// slack.Register(cmd)

	// cmd = hanu.NewCommand("stop", "останавливает службу",
	// 	func(conv hanu.ConversationInterface) {
	// 		var text string
	// 		var cancelled bool
	// 		switch status {
	// 		case 1:
	// 			text = "служба уже остановлена!"
	// 		case 2:
	// 			text = "подождите, служба еще не запущена!"
	// 		case 3:
	// 			text = "терпение, служба уже останавливается!"
	// 		case 4:
	// 			// TODO: добавить отправку предупреждения в скайп
	// 			for i := countdown; i < 0; i-- {
	// 				conv.Reply("%s", "служба будет остановлена через "+
	// 					strconv.Itoa(i)+" секунд.\n\n`cancel` — отмена.")
	// 				select {
	// 				case <-opStatus:
	// 					conv.Reply("%s", "остановка отменена.")
	// 					// TODO: добавить отправку предупреждения в скайп
	// 					cancelled = false // TODO: вроде избыточная переменная, лучше придумать надо
	// 					break
	// 				default:
	// 					time.Sleep(time.Second)
	// 				}
	// 			}
	// 			if !cancelled {
	// 				go execute("stop")
	// 				text = "останавливаю..."
	// 			}
	// 		}
	// 		conv.Reply("%s", text)
	// 	})
	// slack.Register(cmd)

	// cmd = hanu.NewCommand("start", "запускает службу",
	// 	func(conv hanu.ConversationInterface) {
	// 		var text string
	// 		switch status {
	// 		case 1:
	// 			go execute("start")
	// 			text = "запускаю..."
	// 		case 2:
	// 			text = "терпение, служба уже запускается!"
	// 		case 3:
	// 			text = "подождите, служба еще не остановлена!"
	// 		case 4:
	// 			text = "служба уже запущена!"
	// 		}
	// 		conv.Reply("%s", text)
	// 	})
	// slack.Register(cmd)

	// cmd = hanu.NewCommand("cancel", "отменяет запланированную остановку службы",
	// 	func(conv hanu.ConversationInterface) {
	// 		// TODO: придумать красивый способ определять планируется ли остановка
	// 		if status == 4 {
	// 			opStatus <- true
	// 			fmt.Println("opStatus <- true") //debug
	// 		} else {
	// 			conv.Reply("%s", "службу не планировалось останавливать.")
	// 		}
	// 	})
	// slack.Register(cmd)

	// cmd = hanu.NewCommand("version", "версия",
	// 	func(conv hanu.ConversationInterface) {
	// 		conv.Reply("`%s`", response.Version)
	// 	})
	// slack.Register(cmd)

	// cmd = hanu.NewCommand("<word>", "тест",
	// 	func(conv hanu.ConversationInterface) {
	// 		if hook == nil {
	// 			hook = conv
	// 			fmt.Println("hooked now") //debug
	// 		}
	// 		fmt.Println("already hooked") //debug
	// 	})
	// slack.Register(cmd)

	// slack.Listen()
}
