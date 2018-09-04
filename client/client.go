package client

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
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

// Loop --- Вечный цикл проверки статуса.
// 2 раза в секунду опрашивает состояние службы утилитой sc.exe, сравнивает ее
// вывод с константами и перезаписывает переменную Status
func Loop() {
	// TODO первое сообщение должно быть
	// что-то вроде "на момент запуска бота служба была запущена"
	for {
		// TODO replace with switch or something
		out := util.Execute("query")
		if strings.Contains(out, statuses[4]) {
			util.Status = 4
			process(util.Status)
		} else if strings.Contains(out, statuses[1]) {
			util.Status = 1
			process(util.Status)
		} else if strings.Contains(out, statuses[2]) {
			util.Status = 2
			process(util.Status)
		} else if strings.Contains(out, statuses[3]) {
			util.Status = 3
			process(util.Status)
		}
		cooldownMillis := time.Duration(util.Cooldown) * time.Millisecond
		time.Sleep(cooldownMillis)
	}
}

func process(status int8) {
	var n int8
	if status == 4 {
		n = 1
	} else {
		n = status + 1
	}
	if timeLog[n] >= timeLog[status] {
		beep(util.Sounds[status])
		timeLog[status] = time.Now().Unix()
		gui.Change(status)
		if util.Conv != nil {
			util.Conv.Reply(alerts[status])
		}
	}
}

func beep(filename string) {
	f, err := os.Open("sounds\\" + filename)
	if err != nil {
		log.Printf("Не удалось найти звуковой файл (%s)", err)
		return
	}

	s, format, _ := wav.Decode(f)
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	speaker.Play(s)
}
