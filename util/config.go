package util

import (
	"bufio"
	"errors"
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
	"github.com/nlopes/slack"
)

// Ver содержит номер версии
const Ver = "0.9.5"

// TemplateDateFormat содержит формат даты для архива
const TemplateDateFormat = "2006-01-02_15-04"

// alertChannel содержит имя канала для оповещений об изменении статуса службы
const alertChannel = "dominoserver"

const separator = "=" // разделитель для парсинга файлов
const name = "sc"     // имя утилиты для запросов

// Имена файлов
const (
	FilenameConfig    = "autot.cfg"    // имя конфигурационного файла
	FilenameBackup    = "files.bak"    // имя файла с перечнем отправляемых шаблонов
	FilenameAliasList = "aliases.list" // имя файла с перечнем алиасов шаблонов
)

// Списки
var (
	Files   = make(map[string]string) // список отправляемых файлов
	Aliases = make(map[string]string) // список алиасов шаблонов
	Config  = make(map[string]string) // список ключей и значений настроек
)

// TODO Сделать красиво, пока тут глобальные переменные
var (
	API         *slack.Client // API содержит ссылку на клиент новой библиотеки (nlopes/slack)
	ArcFullName string
)

// Разные переменные
var (
	// Status содержит код состояния службы
	Status int // 1 -- остановлена, 2 -- запускается, 3 -- останавливается, 4 -- запущена

	// OpStatus содержит ссылку на канал для отмены остановки службы командой "-"
	OpStatus chan bool

	// Sounds содержит массив со звуковыми файлами
	Sounds = []string{
		"Archspire — Involuntary Doppelgänger.mp3", // этот элемент никогда не используется, а зря
		"", // stopped-sound
		"", // start-pending-sound
		"", // stop-pending-sound
		"", // started-sound
		"", // stop-vote-sound
		"", // beep-sound
	}

/*
	scErrors = map[int]string{
		5:    "Отказано в доступе.",
		50:   "Такой запрос не поддерживается.",
		1060: "Указанная служба не установлена.",
		1061: "Служба в настоящее время не может принимать команды.",
		1062: "Служба не запущена.",
		1056: "Одна копия службы уже запущена.",
		1639: "?",
		1722: "Сервер RPC недоступен.",
	}
*/

)

// Переменные из конфигурационного файла
var (
	server    string // полное имя сервера
	Service   string // имя службы
	DestDir   string // путь к папке с подписанными шаблонами
	DataDir   string // путь к папке с шаблонами
	SrcDir    string // путь к папке, в которой содержатся шаблоны для отправки
	Cooldown  int    // время в миллисекундах между запросами о состоянии службы
	Countdown int    // время в секундах с момента запроса на остановку службы до ее остановки
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

// RandomNumber возвращает случайное число от нуля до max
func RandomNumber(max int) int {
	seed := time.Now().UnixNano()
	source := rand.NewSource(seed)
	return rand.New(source).Intn(max)
}

// ContainsFileAlready проверяет есть ли уже в списке отправляемых файлов указанный файл
func ContainsFileAlready(f string) bool {
	for _, x := range Files {
		if strings.ToLower(x) == f {
			return true
		}
	}
	return false
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

// GetAlertChannel возвращает канал для оповещений
func GetAlertChannel() (slack.Channel, error) {
	if API == nil {
		return slack.Channel{}, errors.New("nlopes/slack.Client is not initialized")
	}

	channels, err := API.GetChannels(true)
	if err != nil {
		return slack.Channel{}, err
	}

	for _, channel := range channels {
		if channel.Name == alertChannel {
			return channel, nil
		}
	}

	return slack.Channel{}, errors.New("Channel \"" + alertChannel + "\" not found")
}

// UpdateBackupFile обновляет файл-бэкап списком отправляемых файлов
func UpdateBackupFile() {
	if err := os.Remove(FilenameBackup); err != nil {
		log.Printf("Ошибка при попытке удаления файла бэкапа (%s)", err)
		return
	}

	var data string
	for number, filename := range Files {
		data += number + separator + filename + "\n"
	}
	writeToFile(FilenameBackup, os.O_CREATE|os.O_WRONLY, data)
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

// CopyFile копирует файл
func CopyFile(source string, destination string) error {
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
	const errorMessage = "Не найден ключ %q в конфигурационном файле. Использую %q"
	const errorMessageWithDigit = "Не найден ключ %q в конфигурационном файле. Использую \"%d\""
	const convertionErrorMessage = "Неверное значение ключа %q. Использую \"%d\" (%s)"

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

	key = "dest-dir"
	if DestDir, ok = Config[key]; !ok {
		log.Printf(errorMessage, key, DestDir)
	}

	key = "data-dir"
	if DataDir, ok = Config[key]; !ok {
		log.Printf(errorMessage, key, DataDir)
	}

	key = "src-dir"
	if SrcDir, ok = Config[key]; !ok {
		log.Printf(errorMessage, key, SrcDir)
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
