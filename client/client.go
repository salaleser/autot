package client

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"salaleser.ru/autot/command"
	"salaleser.ru/autot/gui"
	"salaleser.ru/autot/poster"

	"github.com/nlopes/slack"
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
)

// type timeLog struct {
// 	lastStartPending int64
// 	lastStopPending  int64
// 	lastStopped      int64
// 	lastStarted      int64
// }

// Run runs Slack Bot
func Run(token string) error {
	gui.SetTitle("Autot Server " + util.Ver)
	util.ReadFileIntoMap(util.FilenameAliasList, util.Aliases)
	util.ReadFileIntoMap(util.FilenameBackup, util.Files)
	util.API = slack.New(token)
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

type handler func(c *slack.Client, rtm *slack.RTM, ev *slack.MessageEvent, data []string)

// Command описывает команду
type Command struct {
	Hanlder handler
	Help    string
}

var handlers = map[string]Command{
	"!add": Command{command.AddHandler,
		"Добавляет файл в список _отправляемых_ файлов." +
			"\nИспользование: `!add имя_файла [core | имена_файлов]`" +
			"\ncore -- заменяется на файлы ядра"},
	"!aliases": Command{command.AliasesHandler,
		"Показывает список алиасов шаблонов"},
	"!autot": Command{command.AutotHandler,
		"Останавливает службу, пакует и копирует шаблоны и запускает службу." +
			"\nИспользование: `!autot [имя_файла]`"},
	"!clear":         Command{command.ClearHandler, "Очищает список _отправляемых_ файлов"},
	"!config":        Command{command.ConfigHandler, "Показывает текущие настройки"},
	"привет":         Command{command.GreetHandler, "Отвечает приветствием"},
	"!ping":          Command{command.PingHandler, "Отправляет начальнику сообщение об _отправленных_ файлах"},
	"!pull":          Command{command.PullHandler, "_Отправляет_ файлы"},
	"!push":          Command{command.PushHandler, "_Ставит_ подписанные шаблоны"},
	"!query":         Command{command.QueryHandler, "Показывает текущее состояние службы"},
	"!config reload": Command{command.ConfigReloadHandler, "Перезагружает настройки из файла в память"},
	"!rm":            Command{command.RmHandler, "Удаляет указанный по номеру файл из списка _отправляемых_ файлов"},
	"!start":         Command{command.StartHandler, "Запускает службу"},
	"!status":        Command{command.StatusHandler, "Показывает список _отправляемых_ файлов"},
	"!stop":          Command{command.StopHandler, "Останавливает службу"},
	"!ver":           Command{command.VerHandler, "Показывает краткую информацию о боте"},
	"-":              Command{command.VoteNegativetHandler, "Отменяет запланированную остановку службы"},
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
	if key == "!help" {
		go handleHelp(c, rtm, ev, cmds)
	}
	if f, ok := handlers[key]; ok {
		go f.Hanlder(c, rtm, ev, cmds)
		return
	}

	if len(key) > 0 && key[0] == '!' && key != "!help" {
		poster.PostWarning(ev.Channel, fmt.Sprintf("Команда «%s» не поддерживается", key),
			"", "`!help` — список поддерживаемых команд")
	}
}

func handleHelp(c *slack.Client, rtm *slack.RTM, ev *slack.MessageEvent, args []string) {
	if len(args) == 1 {
		var text string
		for name, command := range handlers {
			text += "`" + name + "` — " + command.Help + "\n"
		}
		poster.Post(ev.Channel, "Список зарегистрированных команд:", text,
			"`!help <имя_команды>` — подробное описание команды")
		return
	}

	if f, ok := handlers[args[1]]; ok {
		if f.Help == "" {
			text := ""
			footer := ""
			isKompotStyle := util.RandomNumber(100) > 80
			if isKompotStyle {
				text = "Помощи «нигде нет». Просто слов нет. Найдем слова – сделаем помощь." +
					" Вы держитесь здесь, вам всего доброго, хорошего настроения и здоровья."
				footer = "© Компот"
			}
			poster.PostWarning(ev.Channel, fmt.Sprintf("Команда «%s» не поддерживается", args[1]),
				text, footer)
			return
		}

		poster.Post(ev.Channel, "Описание команды «"+args[1]+"»:", f.Help, "")
		return
	}

	poster.PostWarning(ev.Channel, fmt.Sprintf("Команда «%s» не поддерживается", args[1]), "", "")
}

// Loop запускает вечный цикл опроса состояние службы утилитой sc.exe, сравнивает ее
// вывод с константами и перезаписывает переменную Status
func Loop() {
	for {
		out := util.Execute("query")
		if strings.Contains(out, statuses[util.StatusRunning]) {
			if util.Status != util.StatusRunning {
				util.Status = util.StatusRunning
				process(util.Status)
			}
		} else if strings.Contains(out, statuses[util.StatusStopped]) {
			if util.Status != util.StatusStopped {
				util.Status = util.StatusStopped
				process(util.Status)
			}
		} else if strings.Contains(out, statuses[util.StatusStartPending]) {
			if util.Status != util.StatusStartPending {
				util.Status = util.StatusStartPending
				process(util.Status)
			}
		} else if strings.Contains(out, statuses[util.StatusStopPending]) {
			if util.Status != util.StatusStopPending {
				util.Status = util.StatusStopPending
				process(util.Status)
			}
		}
		time.Sleep(time.Duration(util.Cooldown) * time.Millisecond)
	}
}

func process(s int) {
	gui.Change(s)
	util.Beep(util.Sounds[s])

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
