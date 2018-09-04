package command

import (
	"archive/zip"
	"io"
	"os"
	"time"

	"github.com/sbstjn/hanu"
	"salaleser.ru/autot/util"
)

//Pull содерижт функцию, которая упакует файлы и отправит в папку path-kmis
var Pull = func(conv hanu.ConversationInterface) {
	if util.Status != util.StatusStopped {
		errMsg := "Нельзя изменять шаблоны пока служба не остановлена"
		conv.Reply("```%s```", errMsg)
		return
	}

	if len(util.Files) == 0 {
		errMsg := "Список пустой!"
		conv.Reply("```%s```\n`!add <файл(-ы),через,запятую>` — добавить файл(-ы)", errMsg)
		return
	}

	time := time.Now()
	date := time.Format(util.TemplateDateFormat)
	arcName := "Templates_" + date + "_KMIS.zip"
	arcFullName := util.PathKmis + "\\" + arcName
	archiveFile, err := os.Create(arcFullName)
	if err != nil {
		conv.Reply("```Ошибка при попытке создать архив!\n%s```", err)
	}
	defer archiveFile.Close()

	zipWriter := zip.NewWriter(archiveFile)
	defer zipWriter.Close()

	for _, filename := range util.Files {
		file, err := os.Open(util.PathData + filename)
		if err != nil {
			errMsg := "Ошибка при попытке архивировать шаблоны!"
			conv.Reply("```%s\n%s```", errMsg, err)
			return
		}
		defer file.Close()

		info, err := file.Stat()
		if err != nil {
			errMsg := "Ошибка при попытке архивировать шаблоны!"
			conv.Reply("```%s\n%s```", errMsg, err)
			return
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			errMsg := "Ошибка при попытке архивировать шаблоны!"
			conv.Reply("```%s\n%s```", errMsg, err)
			return
		}

		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			errMsg := "Ошибка при попытке архивировать шаблоны!"
			conv.Reply("```%s\n%s```", errMsg, err)
			return
		}

		_, err = io.Copy(writer, file)
		if err != nil {
			errMsg := "Ошибка при попытке архивировать шаблоны!"
			conv.Reply("```%s\n%s```", errMsg, err)
			return
		}
	}

	conv.Reply("Шаблоны в `%s`.", arcFullName)
}
