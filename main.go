package main

// ...
// TODO Тут будет описание

// TODO Сортировать список файлов по ключу
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

func main() {
	// go gl.StartDriver(appMain) // GUI возможно будет подключен позже
	util.ReadFileIntoMap(util.FilenameConfig, util.Config)
	util.ReloadConfig()

	if len(os.Args) > 1 {
		arg := os.Args[1]
		if arg[:5] != "xoxb-" { // TODO добавить проверку паттерном регекспа
			log.Fatal("Неправильный токен! Попробуйте указать токен slack-бота заново.")
		}
		go client.Run(arg)
	}

	go systray.Run(gui.OnReady, gui.OnExit)

	client.Loop()
}
