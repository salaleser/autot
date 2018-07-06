// ...
// TODO Тут будет описание
//

// Ошибки:
// [SC] ControlService: ошибка: 1061:
// Служба в настоящее время не может принимать команды.

// [SC] ControlService: ошибка: 1062:
// Служба не запущена.

// [SC] StartService: ошибка: 1056:
// Одна копия службы уже запущена.

// [SC] OpenSCManager: ошибка: 1722:
// Сервер RPC недоступен.

// TODO На данный момент бот не умеет понимать префикс, сейчас префикс захардкожен в сами команды
// ...

package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/getlantern/systray"
	"github.com/sbstjn/hanu"
	"salaleser.ru/autot/icon"
)

const (
	ver        = "0.5"
	countdown  = 6 // количество секунд до остановки службы
	sourcePath = "\\\\Kmisdevserv\\dominodata\\"
	destPath   = "\\\\KSERVER\\!Common\\КМИС ОП"
)

var (
	status   int8 // 1 -- остановлена, 2 -- запускается, 3 -- останавливается, 4 -- запущена
	statuses = []string{
		"0  ВАНИЛЬНЫЙ_СТАТУС", // этот элемент никогда не используется
		"1  STOPPED",
		"2  START_PENDING",
		"3  STOP_PENDING",
		"4  RUNNING",
	}
	tooltips = []string{
		"Ванильный статус", // этот элемент никогда не используется
		"Служба остановлена.",
		"Служба запускается…",
		"Служба останавливается!",
		"Служба запущена.",
	}

	lastStartPending int64
	lastStopPending  int64
	lastStopped      int64
	lastStarted      int64

	about = "Autot (lite) " + ver

	cooldown = 500

	files []string

	hook hanu.ConversationInterface // FIXME это костыль, смотри комментарий к команде @hook
	item *systray.MenuItem
)

func execute(s string) string {
	o, err := exec.Command("sc", "\\\\Kmisdevserv", s, "IBM Domino Server (DDOMINODATA)").Output()
	if err != nil {
		log.Println(err)
	}
	return string(o)
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
				log.Println(tooltips[status])
				item.SetTitle("Остановить")
				// TODO звуковой сигнал
				systray.SetIcon(icons.Green)
				systray.SetTooltip(tooltips[status])
				if hook != nil {
					hook.Reply("*СЛУЖБА ЗАПУЩЕНА*")
				}
			}
		} else if strings.Contains(out, statuses[1]) {
			status = 1
			if lastStartPending >= lastStopped {
				lastStopped = time.Now().Unix()
				log.Println(tooltips[status])
				item.SetTitle("Запустить")
				// TODO звуковой сигнал
				systray.SetIcon(icons.Red)
				systray.SetTooltip(tooltips[status])
				if hook != nil {
					hook.Reply("*СЛУЖБА ОСТАНОВЛЕНА*")
				}
			}
		} else if strings.Contains(out, statuses[2]) {
			status = 2
			if lastStopPending >= lastStartPending {
				lastStartPending = time.Now().Unix()
				log.Println(tooltips[status])
				item.SetTitle(tooltips[status])
				// TODO звуковой сигнал
				systray.SetIcon(icons.Yellow)
				systray.SetTooltip(tooltips[status])
				if hook != nil {
					hook.Reply("*СЛУЖБА ЗАПУСКАЕТСЯ*")
				}
			}
		} else if strings.Contains(out, statuses[3]) {
			status = 3
			if lastStarted >= lastStopPending {
				lastStopPending = time.Now().Unix()
				log.Println(tooltips[status])
				item.SetTitle(tooltips[status])
				// TODO звуковой сигнал
				systray.SetIcon(icons.Red)
				systray.SetTooltip(tooltips[status])
				if hook != nil {
					hook.Reply("*СЛУЖБА ОСТАНАВЛИВАЕТСЯ*")
				}
			}
		}
		time.Sleep(time.Duration(cooldown) * time.Millisecond)
	}
}

func onReady() {
	systray.SetIcon(icons.Yellow)
	systray.SetTitle("Autot")
	systray.SetTooltip("Автотправитель")
	go func() {
		systray.AddMenuItem(about, "Автотправитель")
		systray.AddSeparator()
		item = systray.AddMenuItem("…", "")
		mQuit := systray.AddMenuItem("Выход", "")
		for {
			select {
			case <-item.ClickedCh:
				if status == 4 {
					go execute("stop")
				} else if status == 1 {
					go execute("start")
				}
			case <-mQuit.ClickedCh:
				systray.Quit()
				os.Exit(0)
			}
		}
	}()
}

func main() {
	go loop()

	onExit := func() {
		// Пример лога. Может быть добавлю что-то подобное
		// now := time.Now()
		// ioutil.WriteFile(fmt.Sprintf(`on_exit_%d.txt`, now.UnixNano()), []byte(now.String()), 0644)
	}

	if len(os.Args) > 1 {
		arg := os.Args[1]
		if arg[:5] != "xoxb-" { // TODO добавить проверку паттерном регекспа
			log.Fatal("Неправильный токен! Попробуйте указать токен slack-бота заново.")
		}
		go client(arg)
	}
	systray.Run(onReady, onExit)
}

func client(token string) {
	about = "Autot " + ver
Reconnect:
	log.Println("Connecting...")
	slack, err := hanu.New(token)
	if err != nil {
		log.Println(err)
		fmt.Print("Reconnecting in... ")
		for i := 30; i > 0; i-- {
			time.Sleep(time.Second)
			fmt.Print(i, " ")
		}
		fmt.Println("0")
		goto Reconnect
	}
	log.Println("Connected!")

	opStatus := make(chan bool) // канал для отмены остановки

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
			conv.Reply(text)
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
				conv.Reply("*ВНИМАНИЕ!*\nСлужба будет остановлена через " +
					strconv.Itoa(countdown) + " секунд!")
				for i := countdown; i > 0; i-- {
					select {
					case <-opStatus:
						conv.Reply("Остановка службы отменена.")
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
			conv.Reply(text)
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
			conv.Reply(text)
		})
	slack.Register(cmd)

	cmd = hanu.NewCommand("!cancel", "отменяет запланированную остановку службы",
		func(conv hanu.ConversationInterface) {
			// TODO придумать красивый способ определять планируется ли остановка
			if status == 4 {
				opStatus <- true
				fmt.Println("opStatus <- true") //debug
			} else {
				conv.Reply("Службу не планировалось останавливать.")
			}
		})
	slack.Register(cmd)

	cmd = hanu.NewCommand("!version", "версия", // FIXME hardcode
		func(conv hanu.ConversationInterface) {
			conv.Reply(ver)
		})
	slack.Register(cmd)

	// Эта команда нужна только для присвоения переменной hook экземпляра Conversation.
	// TODO Я пока не смог найти способ как публиковать сообщения ботом в произвольный канал.
	cmd = hanu.NewCommand("!hook", "включает оповещение об изменении состояния Службы в этот канал",
		func(conv hanu.ConversationInterface) {
			if hook == nil {
				hook = conv
				hook.Reply("Оповещения об изменении состояния Службы будут приходить в этот канал.")
			}
		})
	slack.Register(cmd)

	cmd = hanu.NewCommand("!add <файлы,через,запятую>", "обновляет список отправляемых файлов",
		func(conv hanu.ConversationInterface) {
			s, err := conv.String("файлы,через,запятую")
			if err != nil {
				conv.Reply("```Ошибка!\n" + err.Error() + "```")
				return
			}
			newFiles := strings.Split(s, ",")
			files = append(files, newFiles...)
			conv.Reply("Список файлов обновлен. Дайте команду `!push` для отправки их в ОП.")
		})
	slack.Register(cmd)

	cmd = hanu.NewCommand("!rm <номер>", "удаляет файл из списка отправляемых по номеру",
		func(conv hanu.ConversationInterface) {
			s, err := conv.String("номер")
			if err != nil {
				conv.Reply("```Ошибка!\n" + err.Error() + "```")
				return
			}
			n, err := strconv.Atoi(s)
			if err != nil {
				conv.Reply("```Ошибка!\n" + err.Error() + "```")
				return
			}
			n--
			if n < 0 || n > len(files) {
				conv.Reply("```Ошибка!\nИндекс вне массива.```")
				return
			}
			f := files[n]

			// FIXME Не очень красивое решение
			var newFiles []string
			for _, file := range files {
				if file != f {
					newFiles = append(newFiles, file)
				}
			}
			files = newFiles

			conv.Reply("Файл `" + f + "` удален из списка.")
		})
	slack.Register(cmd)

	cmd = hanu.NewCommand("!clear", "обнуляет список файлов",
		func(conv hanu.ConversationInterface) {
			files = []string{}
			conv.Reply("Список файлов очищен.")
		})
	slack.Register(cmd)

	cmd = hanu.NewCommand("!status", "показывает список отправляемых файлов",
		func(conv hanu.ConversationInterface) {
			text := "Список отправляемых файлов:\n"
			if len(files) > 0 {
				for i := 1; i <= len(files); i++ {
					n := strconv.FormatInt(int64(i), 10)
					text += n + ". " + files[i-1] + "\n"
				}
			}
			conv.Reply("```" + text + "```")
			conv.Reply("(`!clear` — очистить список, `!rm` — удалить файл)")
		})
	slack.Register(cmd)

	cmd = hanu.NewCommand("!push", "_отправляет_ подготовленные файлы",
		func(conv hanu.ConversationInterface) {
			if len(files) == 0 {
				conv.Reply("Список файлов пустой! Воспользуйтесь сначала командой `!add`.")
				return
			}
			y, m, d := time.Now().Date()
			date := strconv.Itoa(y) + "-" + m.String() + "-" + strconv.Itoa(d)
			fileName := "Templates_" + date + "_KMIS.zip"
			if err := zipFiles(fileName, files); err != nil {
				log.Fatal(err)
			}
			_, err := exec.Command("xcopy", fileName, destPath, "/Y").Output()
			if err != nil {
				log.Fatal(err)
			}
			conv.Reply("Архив с шаблонами в КМИС ОП (`" + fileName + "`).")
		})
	slack.Register(cmd)

	cmd = hanu.NewCommand("!pull", "ставит подписанные шаблоны (в разработке)",
		func(conv hanu.ConversationInterface) {
			conv.Reply("Начинаю распаковку…")
			signedTemplates, err := Unzip(destPath+
				"\\Подписанные\\"+"Templates_388_2018-06-21_KMIS.7z", "temp")
			if err != nil {
				conv.Reply("```Ошибка!\n" + err.Error() + "```")
				return
			}
			conv.Reply("Распаковка файлов завершена.")

			conv.Reply("Начинаю копирование…")
			for _, template := range signedTemplates {
				_, err := exec.Command("xcopy", "temp\\"+template, "sourcePath", "/Y").Output()
				if err != nil {
					conv.Reply("```Ошибка!\n" + err.Error() + "```")
					return
				}
			}
			conv.Reply("Установка подписанных шаблонов завершена. //тест")
		})
	slack.Register(cmd)

	slack.Listen()
}

func zipFiles(filename string, files []string) error {
	newfile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newfile.Close()

	zipWriter := zip.NewWriter(newfile)
	defer zipWriter.Close()

	for _, file := range files {

		zipfile, err := os.Open(sourcePath + file)
		if err != nil {
			return err
		}
		defer zipfile.Close()

		info, err := zipfile.Stat()
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// Change to deflate to gain better compression
		// see http://golang.org/pkg/archive/zip/#pkg-constants
		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, zipfile)
		if err != nil {
			return err
		}
	}
	return nil
}

// Unzip will decompress a zip archive, moving all files and folders
// within the zip file (parameter 1) to an output directory (parameter 2).
func Unzip(src string, dest string) ([]string, error) {

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}
		defer rc.Close()

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)
		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {

			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)

		} else {

			// Make File
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return filenames, err
			}

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return filenames, err
			}

			_, err = io.Copy(outFile, rc)

			// Close the file without defer to close before next iteration of loop
			outFile.Close()

			if err != nil {
				return filenames, err
			}

		}
	}
	return filenames, nil
}
