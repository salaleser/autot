// ...
// TODO Тут будет описание
//

// [SC] ControlService: ошибка: 1061:
// Служба в настоящее время не может принимать команды.

// [SC] ControlService: ошибка: 1062:
// Служба не запущена.

// [SC] StartService: ошибка: 1056:
// Одна копия службы уже запущена.

// FIXME На данный момент бот не умеет понимать префикс, сейчас префикс захардкожен в сами команды
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
	countdown = 0 // количество секунд до остановки службы
)

var (
	status           int8 // 1 -- остановлена, 2 -- запускается, 3 -- останавливается, 4 -- запущена
	lastStartPending int64
	lastStopPending  int64
	lastStopped      int64
	lastStarted      int64

	statuses = []string{
		"0  ВАНИЛЬНЫЙ_СТАТУС",
		"1  STOPPED",
		"2  START_PENDING",
		"3  STOP_PENDING",
		"4  RUNNING",
	}

	hook hanu.ConversationInterface // FIXME это костыль, смотри комментарий к команде @hook
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

// Вечный цикл проверки статуса.
// 2 раза в секунду опрашивает состояние службы утилитой sc.exe, сравнивает ее
// вывод с константами и перезаписывает переменную status
func loop() {
	// TODO первое сообщение должно быть
	// что-то вроде "на момент запуска бота служба была запущена"
	for {
		out := execute("query")
		if strings.Contains(out, statuses[4]) { // TODO связать "statuses[x]" и "status = x"
			status = 4
			if lastStopped >= lastStarted {
				lastStarted = time.Now().Unix()
				log.Println("Служба запущена.")
				if hook != nil {
					hook.Reply("%s", "*СЛУЖБА ЗАПУЩЕНА*")
				}
			}
		} else if strings.Contains(out, statuses[1]) {
			status = 1
			if lastStartPending >= lastStopped {
				lastStopped = time.Now().Unix()
				log.Println("Служба остановлена.")
				if hook != nil {
					hook.Reply("%s", "*СЛУЖБА ОСТАНОВЛЕНА*")
				}
			}
		} else if strings.Contains(out, statuses[2]) {
			status = 2
			if lastStopPending >= lastStartPending {
				lastStartPending = time.Now().Unix()
				log.Println("Служба запускается.")
				if hook != nil {
					hook.Reply("%s", "*СЛУЖБА ЗАПУСКАЕТСЯ*")
				}
			}
		} else if strings.Contains(out, statuses[3]) {
			status = 3
			if lastStarted >= lastStopPending {
				lastStopPending = time.Now().Unix()
				log.Println("Служба останавливается!")
				if hook != nil {
					hook.Reply("%s", "*СЛУЖБА ОСТАНАВЛИВАЕТСЯ*")
				}
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func main() {
	f, err := os.Open("config.json")
	if err != nil {
		log.Fatal("Конфигурационный файл не найден!")
	}

	body := make([]byte, 86) //FIXME hardcode!

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
	log.Println("Connecting...")
	slack, err := hanu.New(response.Token)
	if err != nil {
		log.Println(err)
		fmt.Print("Reconnecting in... ")
		for i := 30; i > 0; i-- {
			time.Sleep(time.Second)
			fmt.Print(i, " ")
		}
		fmt.Println("0")
		goto Reconnect // goto нормальная тема
	}

	log.Println("Connected!")

	opStatus := make(chan bool) // канал для отмены остановки

	go loop() // Андрею не понравится, но по-моему так читается легче

	var cmd hanu.Command
	cmd = hanu.NewCommand("!query", "проверяет состояние службы", // FIXME hardcode
		func(conv hanu.ConversationInterface) {
			isTimeToJoke := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100) > 95
			var text string
			switch status {
			case 1:
				if isTimeToJoke {
					text = "Служба лежит!"
				} else {
					text = "Служба остановлена!"
				}
			case 2:
				text = "Служба запускается…"
			case 3:
				text = "Служба останавливается…"
			case 4:
				if isTimeToJoke {
					text = "Лотусист спит, служба идет…"
				} else {
					text = "Служба работает."
				}
			}
			conv.Reply("%s", text)
		})
	slack.Register(cmd)

	cmd = hanu.NewCommand("!stop", "останавливает службу", // FIXME hardcode
		func(conv hanu.ConversationInterface) {
			var text string
			var cancelled bool
			switch status {
			case 1:
				text = "Служба уже остановлена!"
			case 2:
				text = "Подождите, служба еще не запущена!"
			case 3:
				text = "Терпение, служба уже останавливается!"
			case 4:
				// TODO добавить отправку предупреждения в скайп
				for i := countdown; i > 0; i-- {
					conv.Reply("%s", "Служба будет остановлена через "+
						strconv.Itoa(i)+" сек. (`!cancel` — отмена)")
					select {
					case <-opStatus:
						conv.Reply("%s", "Остановка службы отменена.")
						// TODO добавить отправку предупреждения в скайп
						cancelled = false // TODO вроде избыточная переменная, лучше придумать надо
						break
					default:
						time.Sleep(time.Second)
					}
				}
				if !cancelled {
					go execute("stop")
					text = "Останавливаю…"
				}
			}
			conv.Reply("%s", text)
		})
	slack.Register(cmd)

	cmd = hanu.NewCommand("!start", "запускает службу", // FIXME hardcode
		func(conv hanu.ConversationInterface) {
			var text string
			switch status {
			case 1:
				go execute("start")
				text = "Запускаю…"
			case 2:
				text = "Терпение, служба уже запускается!"
			case 3:
				text = "Подождите, служба еще не остановлена!"
			case 4:
				text = "Служба уже запущена!"
			}
			conv.Reply("%s", text)
		})
	slack.Register(cmd)

	cmd = hanu.NewCommand("!cancel", "отменяет запланированную остановку службы",
		func(conv hanu.ConversationInterface) {
			// TODO придумать красивый способ определять планируется ли остановка
			if status == 4 {
				opStatus <- true
				fmt.Println("opStatus <- true") //debug
			} else {
				conv.Reply("%s", "Службу не планировалось останавливать.")
			}
		})
	slack.Register(cmd)

	cmd = hanu.NewCommand("!version", "версия", // FIXME hardcode
		func(conv hanu.ConversationInterface) {
			conv.Reply("`%s`", response.Version)
		})
	slack.Register(cmd)

	// Эта команда нужна только для присвоения переменной hook экземпляра Conversation.
	// TODO Я пока не смог найти способ как публиковать сообщения ботом в произвольный канал.
	cmd = hanu.NewCommand("!hook", "включает оповещение об изменении состояния Службы в этот канал",
		func(conv hanu.ConversationInterface) {
			if hook == nil {
				hook = conv
				hook.Reply("%s", "Канал зохвачен!"+
					" (Оповещения об изменении состояния Службы будут приходить сюда)")
			}
		})
	slack.Register(cmd)

	cmd = hanu.NewCommand("!db <файлы,через,запятую>", "обновляет список отправляемых баз данных",
		func(conv hanu.ConversationInterface) {
			conv.Reply("%s", "Список баз данных обновляется…")
			s, _ := conv.String("файлы,через,запятую")
			files := strings.Split(s, ",")
			// TODO обработать файлы
			conv.Reply("%s", "Список баз данных обновлен."+
				" Дайте команду `!send` для отправки их в ОП.")
			text := "Список баз данных, готовых к отправке:```"
			for i := 1; i <= len(files); i++ {
				n := strconv.FormatInt(int64(i), 10)
				text += n + ". " + files[i-1] + "\n"
			}
			conv.Reply("%s", text+"```")
		})
	slack.Register(cmd)

	cmd = hanu.NewCommand("!send ", "копирует указанные в списке !db базы данных в ОП",
		func(conv hanu.ConversationInterface) {
			conv.Reply("%s", "Начинаю копирование…")
			// TODO Добавить копирование файлов
			conv.Reply("%s", "Копирование файлов завершено.")
		})
	slack.Register(cmd)

	slack.Listen()
}
