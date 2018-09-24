package main

// ...
// TODO Тут будет описание
// сохранять данные бота в файле на случай вылета
//

// TODO На данный момент бот не умеет понимать префикс, сейчас префикс захардкожен в сами команды
// ...

import (
	"log"
	"os"

	// "github.com/google/gxui/drivers/gl"

	"github.com/getlantern/systray"
	"salaleser.ru/autot/client"
	"salaleser.ru/autot/gui"
	"salaleser.ru/autot/util"
)

var ()

func main() {
	// go gl.StartDriver(appMain)
	util.ReadFileIntoMap(util.FilenameConfig, util.Config)
	util.ReloadConfig()

	if len(os.Args) > 1 {
		arg := os.Args[1]
		if arg[:5] != "xoxb-" { // TODO добавить проверку паттерном регекспа
			log.Fatal("Неправильный токен! Попробуйте указать токен slack-бота заново.")
		}
		go client.Connect(arg)
		go client.Connect2(arg)
	}

	go systray.Run(gui.OnReady, gui.OnExit)

	client.Loop()
}
