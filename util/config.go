package util

import (
	"bufio"
	"io"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"github.com/sbstjn/hanu"
)

const (
	//Ver содержит номер версии
	Ver = "0.9"

	//TemplateDateFormat содержит формат даты для архива
	TemplateDateFormat = "2006-01-02_15-04"

	//FilenameConfig содержит имя конфигурационного файла
	FilenameConfig = "autot.cfg"

	//FilenameFileList содержит имя файла с перечнем отправляемых шаблонов
	FilenameFileList = "files.bak"

	//FilenameUserList содержит имя файла с перечнем пользователей
	FilenameUserList = "users.list"

	//FilenameAliasList содержит имя файла с перечнем алиасов шаблонов
	FilenameAliasList = "aliases.list"
)

// Статусы службы
const (
	StatusStopped      = 1
	StatusStartPending = 2
	StatusStopPending  = 3
	StatusRunning      = 4
)

var (
	//Status содержит код состояния службы
	Status int // 1 -- остановлена, 2 -- запускается, 3 -- останавливается, 4 -- запущена

	//Sounds содержит массив со звуковыми файлами
	Sounds = []string{
		"Archspire — Involuntary Doppelgänger.mp3", // этот элемент никогда не используется, а зря
		"", // stopped-sound
		"", // start-pending-sound
		"", // stop-pending-sound
		"", // started-sound
		"", // stop-vote-sound
		"", // beep-sound
	}

	//Files содержит список отправляемых файлов
	Files = make(map[string]string)

	//Players содержит список игроков в КНБ
	Players = make(map[string]string)

	//Aliases содержит список алиасов шаблонов
	Aliases = make(map[string]string)

	//Users содержит список пользователей
	Users = make(map[string]string)

	//Config содержит список ключей и значений настроек
	Config = make(map[string]string)

	separator = "="
	name      = "sc"
	server    string

	// Service содержит имя службы
	Service string

	// PathSigned содержит путь к папке с подписанными шаблонами
	PathSigned string

	// PathTemp содержит путь к временной папке
	PathTemp string

	// PathData содержит путь к папке с шаблонами
	PathData string

	// PathKmis содержит путь к папке, в которой содержатся шаблоны для отправки
	PathKmis string

	// PathBackup содержит путь к папке, в которой содержатся резервные копии шаблонов
	PathBackup string

	// Cooldown содержит время в миллисекундах между запросами о состоянии службыы
	Cooldown int

	// Countdown содержит время в секундах с момента запроса на остановку службы до ее остановки
	Countdown int

	// Conv содержит ссылку на канал, в который необходимо сообщать об изменениях состояния службы
	Conv hanu.ConversationInterface

	// OpStatus содержит ссылку на канал для отмены остановки службы командой "-"
	OpStatus chan bool

	errors = map[int]string{
		5:    "Отказано в доступе.",
		50:   "Такой запрос не поддерживается.",
		1060: "Указанная служба не установлена.",
		1061: "Служба в настоящее время не может принимать команды.",
		1062: "Служба не запущена.",
		1056: "Одна копия службы уже запущена.",
		1639: "?",
		1722: "Сервер RPC недоступен.",
	}
)

// Execute выполняет команду name с аргументами
func Execute(s string) string {
	o, err := exec.Command(name, server, s, Service).Output()
	if err != nil {
		log.Printf("Ошибка выполнения команды \"%s %s %s %s\" (%s)", name, server, s, Service, err)
	}
	return string(o)
}

// ReadFileIntoMap считывает файл в мэп
func ReadFileIntoMap(filename string, _map map[string]string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Не удалось найти файл %q (%s)", filename, err)
		return
	}
	defer file.Close()

	var line []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line = strings.Split(scanner.Text(), separator)
		if len(line) == 2 {
			key := line[0]
			value := line[1]
			_map[key] = value
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
}

//RandomNumber возвращает случайное число от нуля до max
func RandomNumber(max int) int {
	seed := time.Now().UnixNano()
	source := rand.NewSource(seed)
	return rand.New(source).Intn(max)
}

// Beep воспроизводит звук
func Beep(filename string) {
	f, err := os.Open("sounds\\" + filename)
	if err != nil {
		log.Printf("Не удалось найти звуковой файл (%s)", err)
		return
	}

	s, format, _ := wav.Decode(f)
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	speaker.Play(s)
}

// DeleteFileList удаляет файл-бэкап списка отправляемых файлов
func DeleteFileList() {
	if err := os.Remove(FilenameFileList); err != nil {
		log.Printf("Ошибка при попытке удаления файла %s (%s)", FilenameFileList, err)
	}
}

// SaveFileList обновляет файл-бэкап списком отправляемых файлов
func SaveFileList() {
	var data string
	for number, filename := range Files {
		data += number + separator + filename + "\n"
	}
	writeToFile(FilenameFileList, os.O_CREATE|os.O_WRONLY, data)
}

func writeToFile(filename string, flag int, data string) {
	f, err := os.OpenFile(filename, flag, 0600)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	if _, err = f.WriteString(data); err != nil {
		log.Println(err)
	}
}

// not used now, just example
func copy(source string, destination string) error {
	src, err := os.Open(source)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.OpenFile(destination, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	return nil
}

// ReloadConfig перечитывет конфигурационный файл
func ReloadConfig() {
	var ok bool
	var err error
	var key string
	errorMessage := "Не найден ключ %q в конфигурационном файле. Использую %q"
	errorMessageWithDigit := "Не найден ключ %q в конфигурационном файле. Использую \"%d\""
	convertionErrorMessage := "Неверное значение ключа %q. Использую \"%d\" (%s)"

	key = "server"
	if server, ok = Config[key]; !ok {
		server = "\\\\127.0.0.1"
		log.Printf(errorMessage, key, server)
	}

	key = "service"
	if Service, ok = Config[key]; !ok {
		Service = "Audiosrv"
		log.Printf(errorMessage, key, Service)
	}

	key = "path-signed"
	if PathSigned, ok = Config[key]; !ok {
		log.Printf(errorMessage, key, PathSigned)
	}

	key = "path-temp"
	if PathTemp, ok = Config[key]; !ok {
		log.Printf(errorMessage, key, PathTemp)
	}

	key = "path-data"
	if PathData, ok = Config[key]; !ok {
		log.Printf(errorMessage, key, PathData)
	}

	key = "path-kmis"
	if PathKmis, ok = Config[key]; !ok {
		log.Printf(errorMessage, key, PathKmis)
	}

	key = "path-backup"
	if PathBackup, ok = Config[key]; !ok {
		log.Printf(errorMessage, key, PathBackup)
	}

	key = "stopped-sound"
	if Sounds[1], ok = Config[key]; !ok {
		Sounds[1] = "stopped.wav"
		log.Printf(errorMessage, key, Sounds[1])
	}

	key = "start-pending-sound"
	if Sounds[2], ok = Config[key]; !ok {
		Sounds[2] = "start_pending.wav"
		log.Printf(errorMessage, key, Sounds[2])
	}

	key = "stop-pending-sound"
	if Sounds[3], ok = Config[key]; !ok {
		Sounds[3] = "stop_pending.wav"
		log.Printf(errorMessage, key, Sounds[3])
	}

	key = "started-sound"
	if Sounds[4], ok = Config[key]; !ok {
		Sounds[4] = "started.wav"
		log.Printf(errorMessage, key, Sounds[4])
	}

	key = "stop-poll-sound"
	if Sounds[5], ok = Config[key]; !ok {
		Sounds[5] = "stop-poll.wav"
		log.Printf(errorMessage, key, Sounds[5])
	}

	key = "beep-sound"
	if Sounds[6], ok = Config[key]; !ok {
		Sounds[6] = "beep.wav"
		log.Printf(errorMessage, key, Sounds[6])
	}

	var defaultValue int

	key = "cooldown"
	defaultValue = 666
	if _, ok = Config[key]; !ok {
		Cooldown = defaultValue
		log.Printf(errorMessageWithDigit, key, Cooldown)
	} else if Cooldown, err = strconv.Atoi(Config[key]); err != nil {
		Cooldown = defaultValue
		log.Printf(convertionErrorMessage, key, Cooldown, err)
	}

	key = "countdown"
	defaultValue = 13
	if _, ok = Config[key]; !ok {
		Countdown = defaultValue
		log.Printf(errorMessageWithDigit, key, Countdown)
	} else if Countdown, err = strconv.Atoi(Config[key]); err != nil {
		Countdown = defaultValue
		log.Printf(convertionErrorMessage, key, Countdown, err)
	}
}
