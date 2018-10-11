package client

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"salaleser.ru/autot/gui"

	"github.com/nlopes/slack"
	"salaleser.ru/autot/command"
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

type handler func(c *slack.Client, rtm *slack.RTM, ev *slack.MessageEvent, data []string)

// Run runs Slack Bot
func Run(token string) error {
	gui.SetTitle("Autot Server " + util.Ver)
	util.ReadFileIntoMap(util.FilenameAliasList, util.Aliases)
	util.ReadFileIntoMap(util.FilenameBackup, util.Files)
	util.API = slack.New(token) // Второй бот (временно их два одновременно)
	rtm := util.API.NewRTM()
	err := make(chan error)
	go serveEvents(util.API, rtm, err)
	go rtm.ManageConnection()
	return <-err
}

func serveEvents(c *slack.Client, rtm *slack.RTM, err chan error) {
	var currentUser string
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				currentUser = fmt.Sprintf("<@%s>", ev.Info.User.ID)
				log.Printf("Подключен nlopes/slack как %q", currentUser)
			case *slack.HelloEvent:

			case *slack.InvalidAuthEvent:
				log.Printf("InvalidAuthEvent %+v\n", ev)
				err <- errors.New("token was not provided")
			case *slack.ConnectionErrorEvent:
				log.Printf("ConnectionErrorEvent %+v\n", ev)
				err <- errors.New("connection error")
			case *slack.MessageEvent:
				handleMessageEvent(c, rtm, ev, currentUser)
			}
		}
	}
}

var handlers = map[string]handler{
	"!add":     command.AddHandler,
	"!aliases": command.AliasesHandler,
	"!autot":   command.AutotHandler,
	"!clear":   command.ClearHandler,
	"!config":  command.ConfigHandler,
	"привет":   command.GreetHandler,
	"!help":    command.HelpHandler,
	"!ping":    command.PingHandler,
	"!pull":    command.PullHandler,
	// "!push":          command.PushHandler,
	"!query":         command.QueryHandler,
	"!config reload": command.ConfigReloadHandler,
	"!rm":            command.RmHandler,
	"!start":         command.StartHandler,
	"!status":        command.StatusHandler,
	"!stop":          command.StopHandler,
	"!ver":           command.VerHandler,
	"-":              command.VoteNegativetHandler,
}

func handleMessageEvent(c *slack.Client, rtm *slack.RTM, ev *slack.MessageEvent, u string) {
	cmds := strings.Split(ev.Text, " ")
	var key string
	l := len(cmds)
	switch {
	case l == 1:
		key = cmds[0]
	case l >= 2:
		if cmds[0] == u {
			key = cmds[1]
			cmds = cmds[1:]
		} else {
			key = cmds[0]
		}
	default:
		return
	}
	if f, ok := handlers[key]; ok {
		f(c, rtm, ev, cmds)
		return
	}
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
		time.Sleep(time.Duration(util.Cooldown) * time.Millisecond)
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
