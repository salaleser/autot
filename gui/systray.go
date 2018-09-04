package gui

import (
	"log"
	"os"

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

// SetTitle устанавливает заголовок в верхнем пункте контекстного меню бота в системном трее винды
func SetTitle(_title string) {
	title = _title
}

// Change обновляет контекстное меню
func Change(status int) {
	log.Println(tooltips[status])
	item.SetTitle(titles[status])
	systray.SetIcon(iconsArray[status])
	systray.SetTooltip(tooltips[status])
}

// OnReady запускает иконку в системном трее
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
			if util.Status == util.StatusRunning {
				go util.Execute("stop")
			} else if util.Status == util.StatusStopped {
				go util.Execute("start")
			}
		case <-mQuit.ClickedCh:
			systray.Quit()
			os.Exit(0)
		}
	}
}

// OnExit пока ничего не делает
func OnExit() {
	// now := time.Now()
	// filename := fmt.Sprintf(`on_exit_%d.txt`, now.UnixNano())
	// data := []byte(now.String())
	// ioutil.WriteFile(filename, data, 0644)
}
