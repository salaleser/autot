package command

import (
	"archive/zip"
	"io"
	"os"
	"time"

	"github.com/sbstjn/hanu"
	"salaleser.ru/autot/util"
)

//Pull упакует файлы и отправит в папку path-kmis
var Pull = func(conv hanu.ConversationInterface) {
	if util.Status != 1 {
		conv.Reply("```Нельзя изменять шаблоны пока служба не остановлена```")
		return
	}

	if len(util.Files) == 0 {
		conv.Reply("```Список пустой!```\n`!add <файл(-ы),через,запятую>` — добавить файл(-ы)")
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
			conv.Reply("```Ошибка при попытке архивировать шаблоны!\n%s```", err)
			return
		}
		defer file.Close()

		info, err := file.Stat()
		if err != nil {
			conv.Reply("```Ошибка при попытке архивировать шаблоны!\n%s```", err)
			return
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			conv.Reply("```Ошибка при попытке архивировать шаблоны!\n%s```", err)
			return
		}

		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			conv.Reply("```Ошибка при попытке архивировать шаблоны!\n%s```", err)
			return
		}

		_, err = io.Copy(writer, file)
		if err != nil {
			conv.Reply("```Ошибка при попытке архивировать шаблоны!\n%s```", err)
			return
		}
	}

	conv.Reply("Шаблоны в `%s`.", arcFullName)
}
