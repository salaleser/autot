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
		const errMsg = "Нельзя изменять шаблоны пока служба не остановлена"
		conv.Reply("```%s```", errMsg)
		return
	}

	if len(util.Files) == 0 {
		const errMsg = "Список пустой!"
		const addCommandName = "`!add <файл(-ы),через,запятую>`"
		conv.Reply("```%s```\n%s — добавить файл(-ы)", errMsg, addCommandName)
		return
	}

	time := time.Now()
	date := time.Format(util.TemplateDateFormat)
	arcName := "Templates_" + date + "_KMIS.zip"
	arcFullName := util.SrcDir + "\\" + arcName
	archiveFile, err := os.Create(arcFullName)
	if err != nil {
		conv.Reply("```Ошибка при попытке создать архив!\n%s```", err)
	}
	defer archiveFile.Close()

	zipWriter := zip.NewWriter(archiveFile)
	defer zipWriter.Close()

	for _, filename := range util.Files {
		file, err := os.Open(util.DataDir + filename)
		if err != nil {
			const errMsg = "Ошибка при попытке архивировать шаблоны!"
			conv.Reply("```%s\n%s```", errMsg, err)
			return
		}
		defer file.Close()

		info, err := file.Stat()
		if err != nil {
			const errMsg = "Ошибка при попытке архивировать шаблоны!"
			conv.Reply("```%s\n%s```", errMsg, err)
			return
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			const errMsg = "Ошибка при попытке архивировать шаблоны!"
			conv.Reply("```%s\n%s```", errMsg, err)
			return
		}

		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			const errMsg = "Ошибка при попытке архивировать шаблоны!"
			conv.Reply("```%s\n%s```", errMsg, err)
			return
		}

		_, err = io.Copy(writer, file)
		if err != nil {
			const errMsg = "Ошибка при попытке архивировать шаблоны!"
			conv.Reply("```%s\n%s```", errMsg, err)
			return
		}
	}

	util.ArcFullName = arcFullName // временный ужос
	conv.Reply("Шаблоны в `%s`.", arcFullName)

	util.Files = map[string]string{}
	util.UpdateBackupFile()
}
