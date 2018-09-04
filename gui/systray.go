package gui

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/getlantern/systray"
	"salaleser.ru/autot/icon"
	"salaleser.ru/autot/util"
)

var (
	item *systray.MenuItem

	title = "Autot Client " + util.Ver

	tooltips = []string{
		"Ванильный статус", // этот элемент никогда не используется
		"Служба остановлена",
		"Служба запускается…",
		"Служба останавливается!",
		"Служба запущена",
	}

	titles = []string{
		"Ванильный статус", // этот элемент никогда не используется
		"Запустить",
		"Служба запускается…",     // копия из tooltips
		"Служба останавливается!", // копия из tooltips
		"Остановить",
	}

	iconsArray = [][]byte{
		icons.Grey, // этот элемент никогда не используется
		icons.Red,
		icons.Yellow,
		icons.Red,
		icons.Green,
	}
)

func SetTitle(_title string) {
	title = _title
}

func Change(status int8) {
	log.Println(tooltips[status])
	item.SetTitle(titles[status])
	systray.SetIcon(iconsArray[status])
	systray.SetTooltip(tooltips[status])
}

func OnReady() {
	systray.SetIcon(iconsArray[2])
	systray.AddMenuItem(title, "Автотправитель")
	systray.AddSeparator()
	systray.AddMenuItem(util.Service, "service")
	systray.AddSeparator()
	item = systray.AddMenuItem("……………………", "***")
	mQuit := systray.AddMenuItem("Выход", "")
	for {
		select {
		case <-item.ClickedCh:
			if util.Status == 4 {
				go util.Execute("stop")
			} else if util.Status == 1 {
				go util.Execute("start")
			}
		case <-mQuit.ClickedCh:
			systray.Quit()
			os.Exit(0)
		}
	}
}

func OnExit() {
	now := time.Now()
	filename := fmt.Sprintf(`on_exit_%d.txt`, now.UnixNano())
	data := []byte(now.String())
	ioutil.WriteFile(filename, data, 0644)
}
