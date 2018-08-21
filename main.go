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
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	unarr "github.com/gen2brain/go-unarr"
	"github.com/getlantern/systray"
	"github.com/sbstjn/hanu"
	"salaleser.ru/autot/icon"
)

const (
	// TODO Перенести в конфиг
	ver                = "0.8"
	templateDateFormat = "2006-01-02_15-04"

	// TODO Перенести в конфиг
	server  = "\\\\Kmisdevserv"
	service = "IBM Domino Server (DDOMINODATA)"

	// TODO Перенести в конфиг
	pathTemp   = "temp\\"
	pathData   = "\\\\Kmisdevserv\\dominodata\\"
	pathBackup = "\\\\Kmisdevserv\\dominodata\\backup\\"
	pathKmis   = "\\\\KSERVER\\!Common\\КМИС ОП"
	pathSigned = "\\\\KSERVER\\!Common\\КМИС ОП\\Подписанные\\"
)

var (
	status int8 // 1 -- остановлена, 2 -- запускается, 3 -- останавливается, 4 -- запущена

	statuses = []string{
		"0  ВАНИЛЬНЫЙ_СТАТУС", // этот элемент никогда не используется
		"1  STOPPED",
		"2  START_PENDING",
		"3  STOP_PENDING",
		"4  RUNNING",
	}
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
	alerts = []string{
		"Ванильный статус", // этот элемент никогда не используется
		"*СЛУЖБА ОСТАНОВЛЕНА*",
		"*СЛУЖБА ЗАПУСКАЕТСЯ*",
		"*СЛУЖБА ОСТАНАВЛИВАЕТСЯ*",
		"*СЛУЖБА ЗАПУЩЕНА*",
	}
	iconsArray = [][]byte{
		icons.Yellow, // этот элемент никогда не используется
		icons.Red,
		icons.Yellow,
		icons.Red,
		icons.Green,
	}
	timeLog = []int64{
		0, // этот элемент никогда не используется
		0, // lastStartPending int64
		0, // lastStopPending  int64
		0, // lastStopped      int64
		0, // lastStarted      int64
	}

	about = "Autot Client " + ver

	countdown = 6 // количество секунд до остановки службы
	cooldown  = 500

	item     *systray.MenuItem
	files    []string
	convList []hanu.ConversationInterface

	rpsRock      = ":fist: *КАМЕНЬ*"
	rpsPaper     = ":hand: *БУМАГА*"
	rpsScissors  = ":v: *НОЖНИЦЫ*"
	rpsEnabled   bool
	rpsCountdown = 15
	rpsMessage   = " начал состязание! " + strconv.Itoa(rpsCountdown) +
		" секунд на ответ! (`!r` — камень, `!s` — ножницы, `!p` — бумага)"
	rpsActions = []string{
		" _смотрит сурово..._",
		" _хмурит брови..._",
		" _напрягся..._",
		" _уверенно стоит на ногах..._",
		" _постукивает ногой в такт музыки..._",
		" _затаился в укрытии..._",
		" _готовится к выпаду..._",
		" _внимательно рассматривает свою ладонь..._",
		" _сосредоточенно разглядывает пиксель..._",
		" _сдержанно улыбается..._",
		" _ожидает результатов..._",
		" _задумчиво смотрит в даль..._",
		" _заметно нервничает..._",
		" _разглядывает пятно на полу..._",
		" _выглядит гордым..._",
	}
	players = make(map[string]string)

	voteEnabled bool
	voteChan    = make(chan bool)
	votes       = 3
	voters      []string

	// TODO Перенести в отдельный файл
	aliases = map[string]string{
		"KMIS_main.ntf":         "Основной шаблон КМИС",
		"CRDirectory.ntf":       "Центральный справочник",
		"MKCalendar6.ntf":       "Календарь КМИС",
		"kmis_globkalendar.ntf": "Общий календарь ЛПУ",
		"kmis_RSysDir.ntf":      "Региональная НСИ",
		"kmis_FSysDir.ntf":      "Федеральная НСИ",
		"SystemIntro.ntf":       "Начальная страница КМИС",
		"kmis_mes.ntf":          "Справочник медицинских стандартов",
		"kmis_frmstr.nsf":       "Печатные формы КМИС",
		"MKCurrent2.ntf":        "Истории болезни",
		"MKAmbul2.ntf":          "Амбулаторные карты",
		"MKAmbul2M.ntf":         "Амбулаторные карты (mini)",
		"MKPasport2.ntf":        "Паспортная часть",
		"MKPasport2M.ntf":       "Паспортная часть (mini)",
		"MKArhiv2.ntf":          "Архив документов",
		"CSTrash.ntf":           "Корзина",
		"kmisrir_iemk.ntf":      "Интегрированная ЭМК",
		"kmis_kladr.ntf":        "Классификатор адресов",
		"kmis_udlo.ntf":         "Региональная система лекарственного обеспечения",
		"kmis_labserv.ntf":      "Web-сервисы интеграции",
		"MKAKPasp.ntf":          "Паспорт поликлиники",
		"kmis_globdir.ntf":      "Глобальный справочник КМИС",
		"kmis_Autopsy.ntf":      "Журнал патанатомии",
		"kmis_usl.ntf":          "Услуги для КМИС",
	}

	// FIXME hardcode
	lotusmen = map[string]string{
		"UC7GRMGA2": "Максим Паничев",
		"UAQAL4NPR": "Алексей Салиенко",
		"UA63MNKHR": "Павел Боровинский",
		"UAQC2EAUX": "Андрей Бородулин",
		"UANLUENDP": "Александр Кирпу",
	}
)

func execute(s string) string {
	o, err := exec.Command("sc", server, s, service).Output()
	if err != nil {
		log.Println(err)
	}
	return string(o)
}

func beep() {
	f, _ := os.Open("1.wav")
	s, format, _ := wav.Decode(f)
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	speaker.Play(s)
}

func process() {
	var n int8
	if status == 4 {
		n = 1
	} else {
		n = status + 1
	}
	if timeLog[n] >= timeLog[status] {
		beep()
		timeLog[status] = time.Now().Unix()
		log.Println(tooltips[status])
		item.SetTitle(titles[status])
		systray.SetIcon(iconsArray[status])
		systray.SetTooltip(tooltips[status])
		for _, c := range convList {
			c.Reply(alerts[status])
		}
	}
}

// Вечный цикл проверки статуса.
// 2 раза в секунду опрашивает состояние службы утилитой sc.exe, сравнивает ее
// вывод с константами и перезаписывает переменную status
func loop() {
	// TODO первое сообщение должно быть
	// что-то вроде "на момент запуска бота служба была запущена"
	for {
		// TODO replace with switch or something
		out := execute("query")
		if strings.Contains(out, statuses[4]) {
			status = 4
			process()
		} else if strings.Contains(out, statuses[1]) {
			status = 1
			process()
		} else if strings.Contains(out, statuses[2]) {
			status = 2
			process()
		} else if strings.Contains(out, statuses[3]) {
			status = 3
			process()
		}
		time.Sleep(time.Duration(cooldown) * time.Millisecond)
	}
}

func onReady() {
	systray.SetIcon(iconsArray[2])
	systray.SetTitle("Autot Server")
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
	cmd = hanu.NewCommand("!query",
		"проверяет состояние службы", // FIXME hardcode
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
					text = "Служба работает"
				}
			}
			conv.Reply(text)
		})
	slack.Register(cmd)

	cmd = hanu.NewCommand("!stop",
		"останавливает службу", // FIXME hardcode
		func(conv hanu.ConversationInterface) {
			startPoll(conv)

			var text string
			var cancelled bool
			switch status {
			case 1:
				text = "Служба уже остановлена!"
			case 2:
				text = "Подождите, служба еще не запущена!"
			case 3:
				text = "Проявите терпение, служба уже останавливается!"
			case 4:
				// TODO добавить отправку предупреждения в скайп
				conv.Reply("*ВНИМАНИЕ!*\nСлужба будет остановлена через " +
					strconv.Itoa(countdown) + " секунд!")
				for i := countdown; i > 0; i-- {
					select {
					case <-opStatus:
						conv.Reply("Остановка службы отменена")
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

	cmd = hanu.NewCommand("!start",
		"запускает службу", // FIXME hardcode
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

	cmd = hanu.NewCommand("!cancel",
		"отменяет запланированную остановку службы",
		func(conv hanu.ConversationInterface) {
			// TODO придумать красивый способ определять планируется ли остановка
			if status == 4 {
				opStatus <- true
				fmt.Println("opStatus <- true") //debug
			} else {
				conv.Reply("Службу не планировалось останавливать")
			}
		})
	slack.Register(cmd)

	cmd = hanu.NewCommand("!version", "версия", // FIXME hardcode
		func(conv hanu.ConversationInterface) {
			conv.Reply(ver)
		})
	slack.Register(cmd)

	cmd = hanu.NewCommand("!aliases",
		"показывает список алиасов шаблонов",
		func(conv hanu.ConversationInterface) {
			text := "Список алиасов шаблонов:\n"
			columnWidth := 26
			for filename, alias := range aliases {
				spaces := columnWidth - len(filename)
				if spaces < 1 {
					spaces = 1
				}
				text += filename + strings.Repeat(" ", spaces) + alias + "\n"
			}
			conv.Reply("```" + text + "```")
		})
	slack.Register(cmd)

	// Эта команда нужна только для присвоения переменной hook экземпляра Conversation.
	// TODO Я пока не смог найти способ как публиковать сообщения ботом в произвольный канал.
	cmd = hanu.NewCommand("!hook",
		"включает оповещение об изменении состояния Службы в этот канал (в том числе и в ЛС)",
		func(conv hanu.ConversationInterface) {
			for _, c := range convList {
				if c == conv {
					conv.Reply("Этот канал уже есть в списке оповещения.")
					return
				}
			}

			if len(convList) == 0 {
				convList = append(convList, conv)
				conv.Reply("Этот канал добавлен в список оповещения (в разработке).")
			} else {
				convList[0] = conv
				conv.Reply("Канал из первого элемента списка заменен на этот канал (в разработке)")
			}
		})
	slack.Register(cmd)

	cmd = hanu.NewCommand("!add <файлы,через,запятую,без,пробелов>",
		"обновляет список отправляемых файлов",
		func(conv hanu.ConversationInterface) {
			s, err := conv.String("файлы,через,запятую,без,пробелов")
			if err != nil {
				conv.Reply("```Ошибка при разборе перечня файлов!\n" +
					err.Error() + "```")
				return
			}
			newFilenames := strings.Split(s, ",")
			for _, newFilename := range newFilenames {
				for _, file := range files {
					if file == newFilename {
						files = remove(newFilename)
					}
				}
			}

			allTemplates, err := ioutil.ReadDir(pathData)
			if err != nil {
				conv.Reply("```Ошибка при попытке прочитать файлы из папки" + pathData + "!\n" +
					err.Error() + "```")
				return
			}

			patternTemplateWithoutExtension, err := regexp.Compile("^\\w+$")
			if err != nil {
				conv.Reply("```Ошибка!\n" +
					err.Error() + "```")
				return
			}

			patternDatabaseWithoutExtension, err := regexp.Compile("^kmis_frmstr$")
			if err != nil {
				conv.Reply("```Ошибка!\n" +
					err.Error() + "```")
				return
			}

			patternTemplate, err := regexp.Compile("^\\w+\\.ntf$")
			if err != nil {
				conv.Reply("```Ошибка!\n" +
					err.Error() + "```")
				return
			}

			patternDatabase, err := regexp.Compile("^kmis_frmstr\\.nsf$")
			if err != nil {
				conv.Reply("```Ошибка!\n" +
					err.Error() + "```")
				return
			}

			const templateExtension = ".ntf"
			const databaseExtension = ".nsf"
			var count int
			for _, newFilename := range newFilenames {
				for _, templateFile := range allTemplates {
					if templateFile.IsDir() {
						continue
					}

					if patternDatabaseWithoutExtension.MatchString(newFilename) {
						newFilename = newFilename + databaseExtension
					} else if patternTemplateWithoutExtension.MatchString(newFilename) {
						newFilename = newFilename + templateExtension
					}

					isTemplate := patternTemplate.MatchString(newFilename)
					isDatabase := patternDatabase.MatchString(newFilename)
					if newFilename == templateFile.Name() {
						if isTemplate || isDatabase {
							files = append(files, newFilename)
							count++
						} else {
							conv.Reply("```Файл " + newFilename +
								" не подходит и не был добавлен в список!\n```")
						}
					}
				}
			}
			text := "Все файлы прошли проверку и были добавлены в список, " +
				"дайте команду `!push` для отправки их в `" + pathKmis + "`"
			if count == 0 {
				text = "Ни один файл не прошел проверку (регистр учитывается)"
			} else if count != len(newFilenames) {
				text = "Не все файлы прошли проверку"
			}
			conv.Reply(text)
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
				conv.Reply("```Ошибка выполнения команды `!rm`!\n" +
					"Индекс вне массива.```")
				return
			}
			f := files[n]

			files = remove(f)

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
					alias, ok := aliases[files[i-1]]
					if ok {
						alias = " («" + alias + "»)"
					}
					text += strconv.Itoa(i) + ". " + files[i-1] + alias + "\n"
				}
			}
			conv.Reply("```" + text + "```")
			conv.Reply("(`!clear` — очистить список, `!rm` — удалить файл)")
		})
	slack.Register(cmd)

	cmd = hanu.NewCommand("!push", "_отправляет_ подготовленные файлы",
		func(conv hanu.ConversationInterface) {
			if status != 1 {
				conv.Reply("```Ошибка выполнения команды `!push`!\n" +
					"Нельзя изменять шаблоны пока служба не остановлена.```")
				return
			}

			if len(files) == 0 {
				conv.Reply("```Ошибка!\n" +
					"Список файлов пустой! Воспользуйтесь сначала командой `!add`.```")
				return
			}

			time := time.Now()
			date := time.Format(templateDateFormat)
			fileName := "Templates_" + date + "_KMIS.zip"
			if err := zipFiles(fileName, files); err != nil {
				conv.Reply("```Ошибка при попытке архивировать шаблоны!\n" +
					err.Error() + "```")
				return
			}
			_, err = exec.Command("xcopy", fileName, pathKmis, "/Y").Output()
			if err != nil {
				conv.Reply("```Ошибка!\n" +
					err.Error() + "```")
				return
			}
			conv.Reply("Шаблоны в `" + pathKmis + "\\" + fileName + "`.")
		})
	slack.Register(cmd)

	cmd = hanu.NewCommand("!pull", "найдет файл в папке *Подписанные*, "+
		"сохранит резервную копию старых файлов из папки *dominodata* в папку *dominodata\\backup* "+
		"и распакует с заменой файлы из архива в папку *dominodata*",
		func(conv hanu.ConversationInterface) {
			if status != 1 {
				conv.Reply("```Ошибка выполнения команды `!pull`!\n" +
					"Нельзя изменять шаблоны пока служба не остановлена.```")
				return
			}

			signedArchives, err := ioutil.ReadDir(pathSigned)
			if err != nil {
				conv.Reply("```Ошибка при попытке прочитать файлы из папки" + pathSigned + "!\n" +
					err.Error() + "```")
				return
			}
			if len(signedArchives) == 0 {
				conv.Reply("```Ошибка при попытке найти подходящий файл для распаковки!\n" +
					"В папке " + pathSigned + " нет файлов.```")
				return
			}

			var signed os.FileInfo
			if len(signedArchives) == 1 {
				signed = signedArchives[0]
				conv.Reply("выбран файл `" + pathSigned + signed.Name() + "`")
			} else {
				text := "Список подписанных архивов:\n"
				for i := 1; i <= len(signedArchives); i++ {
					n := strconv.FormatInt(int64(i), 10)
					text += n + ". " + signedArchives[i-1].Name() + "\n"
				}
				conv.Reply("```" + text + "```")
				// TODO распаковать и поставить шаблоны
				conv.Reply("(`!pull <номер>` — распаковать и поставить шаблоны (в разработке))")
			}

			_, err = exec.Command("xcopy", pathSigned+signed.Name(), pathTemp, "/Y").Output()
			if err != nil {
				conv.Reply("```Ошибка при копировании архива во временную папку!\n" +
					err.Error() + "```")
				return
			}

			a, err := unarr.NewArchive(pathTemp + signed.Name())
			if err != nil {
				conv.Reply("```Ошибка инициализации архива " + pathTemp + signed.Name() + "!\n" +
					err.Error() + "```")
				return
			}
			defer a.Close()

			signedFilenames, err := a.List()
			if err != nil {
				conv.Reply("```Ошибка при чтении имен файлов архива!\n" + err.Error() + "```")
				return
			}
			for _, n := range signedFilenames {
				_, err := exec.Command("xcopy", pathData+n, pathBackup, "/Y").Output()
				if err != nil {
					conv.Reply("```Ошибка при попытке резервного копирования файла " + n + "!\n" +
						err.Error() + "```")
					return
				}
				os.Remove(pathData + n)
			}

			err = a.Extract(pathData)
			if err != nil {
				conv.Reply("```Ошибка при попытке распаковки архива!\n" +
					err.Error() + "```")
				return
			}

			os.Remove(pathSigned + signed.Name())

			conv.Reply("_Установка подписанных шаблонов завершена успешно._")
		})
	slack.Register(cmd)

	cmd = hanu.NewCommand("!pull <номер>",
		"ставит шаблоны из номера архива, указанного в аргументе",
		func(conv hanu.ConversationInterface) {
			if status != 1 {
				conv.Reply("```Ошибка!\n" +
					"Нельзя изменять шаблоны пока служба не остановлена.```")
				return
			}

			conv.Reply("Ничего не сделано, команда в разработке." +
				"Временное решение: оставьте один архив в папке " + pathSigned)
		})
	slack.Register(cmd)

	cmd = hanu.NewCommand("!user",
		"узнать свой идентификатор (никогда не знаешь что может пригодиться)",
		func(conv hanu.ConversationInterface) {
			user := conv.Message().User()
			conv.Reply("`" + user + "`")
		})
	slack.Register(cmd)

	cmd = hanu.NewCommand("!r",
		"камень",
		func(conv hanu.ConversationInterface) {
			user := conv.Message().User()
			playRps(user, rpsRock)
		})
	slack.Register(cmd)

	cmd = hanu.NewCommand("!s",
		"ножницы",
		func(conv hanu.ConversationInterface) {
			user := conv.Message().User()
			playRps(user, rpsScissors)
		})
	slack.Register(cmd)

	cmd = hanu.NewCommand("!p",
		"бумага",
		func(conv hanu.ConversationInterface) {
			user := conv.Message().User()
			playRps(user, rpsPaper)
		})
	slack.Register(cmd)

	cmd = hanu.NewCommand("\\+", // надо экранировать плюс иначе падает hanu
		"положительный ответ при голосовании (правильно `+`)",
		func(conv hanu.ConversationInterface) {
			if voteEnabled {
				user := conv.Message().User()
				vote(user, true)
			}
		})
	slack.Register(cmd)

	cmd = hanu.NewCommand("\\-", // надо экранировать плюс иначе падает hanu
		"отрицательный ответ при голосовании (правильно `-`)",
		func(conv hanu.ConversationInterface) {
			if voteEnabled {
				user := conv.Message().User()
				vote(user, false)
				conv.Reply("Голосование отменено пользователем " + lotusmen[user])
			}
		})
	slack.Register(cmd)

	slack.Listen()
}

func remove(f string) []string {
	var a []string
	for _, file := range files {
		if file != f {
			a = append(a, file)
		}
	}
	return a
}

// TODO доделать команду
func startPoll(conv hanu.ConversationInterface) {
	voters = []string{}
	voteEnabled = true
	conv.Reply("Запущено голосование за остановку службы. Необходимо получить еще " +
		strconv.Itoa(votes-1) + " голоса. Для голосования достаточно поставить знак `+`.")
	conv.Reply("Скрестите шпаги, лотусисты!")

	user := conv.Message().User()
	vote(user, true)

	<-voteChan
	voteMessage := "Голосование завершено успешно. Список голосовавших:\n"
	for i := 0; i < len(voters); i++ {
		lotusman, ok := lotusmen[voters[i]]
		if !ok {
			lotusman = voters[i]
		}
		voteMessage += lotusman + "\n"
	}
	conv.Reply("```" + voteMessage + "```")
}

func vote(user string, isPositive bool) {
	if isPositive {
		voters = append(voters, user)
		if votes > len(voters) {
			return
		}
		voteChan <- true
	}
	voteEnabled = false
}

func playRps(user string, rps string) {
	username, ok := lotusmen[user]
	if !ok {
		username = user
	}

	if !rpsEnabled {
		players = make(map[string]string)
		rpsEnabled = true
		convList[0].Reply(username + rpsMessage)
		go startRps()
	}
	players[username] = rps
}

func startRps() {
	for i := rpsCountdown; i > 0; i-- {
		time.Sleep(time.Second)
	}
	convList[0].Reply("Состязание завершено! Убрать шпаги в ножны!")

	resultMessage := "*Итоги состязания:*\n"
	var playersList []string
	for player, rps := range players {
		playersList = append(playersList, player)
		resultMessage += player + " _выбрал_ " + rps + "\n"
	}
	convList[0].Reply(resultMessage + "\n_Идет подсчет результатов..._")

	for i := 4; i > 0; i-- {
		time.Sleep(time.Second)
		randomPlayer := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(playersList))
		randomAction := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(rpsActions))
		convList[0].Reply(playersList[randomPlayer] + rpsActions[randomAction])
	}
	convList[0].Reply("_А впрочем, считайте сами, я все равно пока сам не умею_ :man-shrugging:")

	rpsEnabled = false
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

		zipfile, err := os.Open(pathData + file)
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
