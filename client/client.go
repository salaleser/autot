package client

import (
	"fmt"
	"log"
	"strings"
	"time"

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

	util.ReadFileIntoMap(util.FilenameAliasList, util.Aliases)
	util.ReadFileIntoMap(util.FilenameUserList, util.Users)
	util.ReadFileIntoMap(util.FilenameFileList, util.Files)

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

// Loop запускает вечный цикл опроса состояние службы утилитой sc.exe, сравнивает ее
// вывод с константами и перезаписывает переменную Status
func Loop() {
	for {
		out := util.Execute("query")
		if strings.Contains(out, statuses[util.StatusRunning]) {
			util.Status = 4
			process(util.Status)
		} else if strings.Contains(out, statuses[util.StatusStopped]) {
			util.Status = 1
			process(util.Status)
		} else if strings.Contains(out, statuses[util.StatusStartPending]) {
			util.Status = 2
			process(util.Status)
		} else if strings.Contains(out, statuses[util.StatusStopPending]) {
			util.Status = 3
			process(util.Status)
		}
		cooldownMillis := time.Duration(util.Cooldown) * time.Millisecond
		time.Sleep(cooldownMillis)
	}
}

func process(status int) {
	var n int
	if status == 4 {
		n = 1
	} else {
		n = status + 1
	}
	if timeLog[n] >= timeLog[status] {
		util.Beep(util.Sounds[status])
		timeLog[status] = time.Now().Unix()
		gui.Change(status)
		if util.Conv != nil {
			util.Conv.Reply(alerts[status])
		}
	}
}
