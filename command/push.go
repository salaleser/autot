package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"

	unarr "github.com/gen2brain/go-unarr"
	"github.com/nlopes/slack"
	"salaleser.ru/autot/poster"
	"salaleser.ru/autot/util"
)

// PushHandler распаковывает архив с подписанными шаблонами из папки и dest-dir копирует с заменой
// файлы в папку data-dir
func PushHandler(c *slack.Client, rtm *slack.RTM, ev *slack.MessageEvent, data []string) {
	if len(data) < 2 {
		poster.PostError(ev.Channel, "Ошибка!", "Не указан номер файла в папке")
		return
	}
	key := data[1]

	n, err := strconv.Atoi(key)
	if err != nil {
		poster.PostError(ev.Channel, "Ошибка!", fmt.Sprintf("%q не является числом!", key))
		return
	}

	if util.Status != util.StatusStopped {
		poster.PostError(ev.Channel, "Ошибка!",
			"Нельзя изменять шаблоны пока служба не остановлена")
		return
	}

	DestDirFiles, err := ioutil.ReadDir(util.DestDir)
	if err != nil {
		poster.PostError(ev.Channel, fmt.Sprintf("Ошибка при попытке прочитать файлы из папки %q!",
			util.DestDir), err.Error())
		return
	}

	if len(DestDirFiles) == 0 {
		poster.PostError(ev.Channel, "Ошибка!", fmt.Sprintf("В папке %q нет файлов", util.DestDir))
		return
	}

	if len(DestDirFiles) > 1 {
		poster.PostError(ev.Channel, "Ошибка!",
			fmt.Sprintf("В папке %q несколько файлов, не могу выбрать (в разработке)",
				util.DestDir))
		return
	}

	archivePattern, err := regexp.Compile("^.+\\.7z$")
	if err != nil {
		poster.PostError(ev.Channel, "Ошибка!", err.Error())
		return
	}

	if n > len(DestDirFiles) {
		poster.PostError(ev.Channel, "Ошибка!", fmt.Sprintf("Файла с номером %d нет в папке", n))
		return
	}
	n--
	signed := DestDirFiles[n]

	if !archivePattern.MatchString(signed.Name()) {
		poster.PostError(ev.Channel, "Ошибка!",
			fmt.Sprintf("Файл %q не является архивом", signed.Name()))
		return
	}

	tempFile := os.Getenv("TEMP") + "\\" + signed.Name()
	util.CopyFile(util.DestDir+signed.Name(), tempFile)

	archive, err := unarr.NewArchive(tempFile)
	if err != nil {
		poster.PostError(ev.Channel, fmt.Sprintf("Ошибка инициализации архива %q", tempFile), err.Error())
		return
	}
	defer os.Remove(tempFile)
	defer archive.Close()

	if err = archive.Extract(util.DataDir); err != nil {
		poster.PostError(ev.Channel, "Ошибка при попытке распаковки архива", err.Error())
		return
	}

	poster.Post(ev.Channel, "Установка подписанных шаблонов завершена успешно", "", "")
}
