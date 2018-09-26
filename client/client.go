package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
	"github.com/sbstjn/hanu"
	"salaleser.ru/autot/gui"
	"salaleser.ru/autot/loader"
	"salaleser.ru/autot/util"
)

var (
	statuses = []string{
		"0  ВАНИЛЬНЫЙ_СТАТУС", // этот элемент никогда не используется
		"1  STOPPED",
		"2  START_PENDING",
		"3  STOP_PENDING",
		"4  RUNNING",
	}

	alerts = []string{
		"Ванильный статус", // этот элемент никогда не используется
		"*СЛУЖБА ОСТАНОВЛЕНА*",
		"*СЛУЖБА ЗАПУСКАЕТСЯ*",
		"*СЛУЖБА ОСТАНАВЛИВАЕТСЯ*",
		"*СЛУЖБА ЗАПУЩЕНА*",
	}

	colors = []string{
		"Ванильный цвет", // этот элемент никогда не используется
		gui.Grey,
		gui.Yellow,
		gui.Red,
		gui.Green,
	}

	timeLog = []int64{
		0, // этот элемент никогда не используется
		0, // lastStartPending int64
		0, // lastStopPending  int64
		0, // lastStopped      int64
		0, // lastStarted      int64
	}
)

// Connect пытается соединиться с сервером слэка
func Connect(token string) {
	gui.SetTitle("Autot Server " + util.Ver)

Reconnect:
	log.Println("Подключение...")
	bot, err := hanu.New(token)
	if err != nil {
		log.Println(err)
		fmt.Print("Следующая попытка переподключения через... ")
		for i := 30; i > 0; i-- {
			time.Sleep(time.Second)
			fmt.Print(i, " ")
		}
		fmt.Println("0")
		goto Reconnect
	}
	log.Println("Подключен!")

	util.API = slack.New(token) // Второй бот (временно их два одновременно)

	util.ReadFileIntoMap(util.FilenameAliasList, util.Aliases)
	util.ReadFileIntoMap(util.FilenameBackup, util.Files)

	log.Println("Загружаю команды...")
	commands := loader.LoadCommands()
	var cmd hanu.Command
	for i := 0; i < len(commands); i++ {
		cmd = commands[i]
		bot.Register(cmd)
	}
	log.Println(len(commands), "команд загружено")

	bot.Listen()
}

// Connect2 подключает второго бота (более продвинутую библиотеку), который скоро станет основным
func Connect2(token string) {
	http.HandleFunc("/events-endpoint", func(w http.ResponseWriter, r *http.Request) {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		body := buf.String()
		event := slackevents.OptionVerifyToken(&slackevents.TokenComparator{token})
		eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), event)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		if eventsAPIEvent.Type == slackevents.URLVerification {
			log.Println("eventsAPIEvent.Type == slackevents.URLVerification")
			var r *slackevents.ChallengeResponse
			err := json.Unmarshal([]byte(body), &r)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			w.Header().Set("Content-Type", "text")
			w.Write([]byte(r.Challenge))
		}
		if eventsAPIEvent.Type == slackevents.CallbackEvent {
			log.Println("eventsAPIEvent.Type == slackevents.CallbackEvent")
			postParams := slack.PostMessageParameters{}
			innerEvent := eventsAPIEvent.InnerEvent
			switch ev := innerEvent.Data.(type) {
			case *slackevents.AppMentionEvent:
				util.API.PostMessage(ev.Channel, "Yes, hello.", postParams)
			case *slackevents.MessageEvent:
				if ev.Text == "333" {
					util.API.PostMessage(ev.Channel, "333?", postParams)
				}
			}
		}
	})
	fmt.Println("[nlopes] Server listening")
	http.ListenAndServe(":3000", nil)
}

// Loop запускает вечный цикл опроса состояние службы утилитой sc.exe, сравнивает ее
// вывод с константами и перезаписывает переменную Status
func Loop() {
	for {
		out := util.Execute("query")
		if strings.Contains(out, statuses[util.StatusRunning]) {
			util.Status = util.StatusRunning
		} else if strings.Contains(out, statuses[util.StatusStopped]) {
			util.Status = util.StatusStopped
		} else if strings.Contains(out, statuses[util.StatusStartPending]) {
			util.Status = util.StatusStartPending
		} else if strings.Contains(out, statuses[util.StatusStopPending]) {
			util.Status = util.StatusStopPending
		}
		process(util.Status)
		cd := time.Duration(util.Cooldown) * time.Millisecond
		time.Sleep(cd)
	}
}

func process(s int) {
	n := s + 1
	if s == util.StatusRunning {
		n = 1
	}
	if timeLog[n] >= timeLog[s] {
		util.Beep(util.Sounds[s])
		timeLog[s] = time.Now().Unix()
		gui.Change(s)

		alertChannel, err := util.GetAlertChannel()
		if err == nil {
			params := slack.PostMessageParameters{}
			attachment := slack.Attachment{
				Color:  colors[s],
				Text:   alerts[s],
				Footer: "Оповещение об изменении статуса службы",
			}
			params.Attachments = []slack.Attachment{attachment}
			params.AsUser = true

			util.API.PostMessage(alertChannel.ID, "", params)
		}
	}
}
