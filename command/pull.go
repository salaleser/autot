package command

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/nlopes/slack"
	"salaleser.ru/autot/poster"
	"salaleser.ru/autot/util"
)

// PullHandler пакует файлы и отправляет в папку path-kmis
func PullHandler(c *slack.Client, rtm *slack.RTM, ev *slack.MessageEvent, data []string) {
	if util.Status != util.StatusStopped {
		poster.PostError(ev.Channel, "Ошибка!",
			"Нельзя изменять шаблоны пока служба не остановлена")
		return
	}

	if len(util.Files) == 0 {
		poster.PostError(ev.Channel,
			"Список пустой!", "*!add <файл(-ы)_через_пробелы>* — добавить файл(-ы)")
		return
	}

	time := time.Now()
	date := time.Format(util.TemplateDateFormat)
	arcName := "Templates_" + date + "_KMIS.zip"
	arcFullName := util.SrcDir + "\\" + arcName
	archiveFile, err := os.Create(arcFullName)
	if err != nil {
		poster.PostError(ev.Channel, "Ошибка при попытке создать архив!", err.Error())
		return
	}
	defer archiveFile.Close()

	zipWriter := zip.NewWriter(archiveFile)
	defer zipWriter.Close()

	for _, filename := range util.Files {
		file, err := os.Open(util.DataDir + filename)
		if err != nil {
			poster.PostError(ev.Channel, "Ошибка при попытке архивировать шаблоны!", err.Error())
			return
		}
		defer file.Close()

		info, err := file.Stat()
		if err != nil {
			poster.PostError(ev.Channel, "Ошибка при попытке архивировать шаблоны!", err.Error())
			return
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			poster.PostError(ev.Channel, "Ошибка при попытке архивировать шаблоны!", err.Error())
			return
		}

		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			poster.PostError(ev.Channel, "Ошибка при попытке архивировать шаблоны!", err.Error())
			return
		}

		_, err = io.Copy(writer, file)
		if err != nil {
			poster.PostError(ev.Channel, "Ошибка при попытке архивировать шаблоны!", err.Error())
			return
		}
	}

	util.ArcFullName = arcFullName // временный ужос
	poster.Post(ev.Channel, "Файлы успешно отправлены",
		fmt.Sprintf("Шаблоны в `%s`.", arcFullName), "file://kserver/!Common/КМИС%20ОП/")

	util.Files = map[string]string{}
	util.UpdateBackupFile()
}
