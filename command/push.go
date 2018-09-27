package command

import (
	"io/ioutil"
	"os"
	"regexp"

	unarr "github.com/gen2brain/go-unarr"
	"github.com/sbstjn/hanu"
	"salaleser.ru/autot/util"
)

// Push содержит функцию, которая распакует архив с подписанными шаблонами из папки и dest-dir
// копирует с заменой файлы в папку data-dir
var Push = func(conv hanu.ConversationInterface) {
	if util.Status != util.StatusStopped {
		conv.Reply("```Нельзя изменять шаблоны пока служба не остановлена```")
		return
	}

	DestDirFiles, err := ioutil.ReadDir(util.DestDir)
	if err != nil {
		conv.Reply("```Ошибка при попытке прочитать файлы из папки %q\n```", util.DestDir, err)
		return
	}

	if len(DestDirFiles) == 0 {
		conv.Reply("```В папке %q нет файлов```", util.DestDir)
		return
	}

	if len(DestDirFiles) > 1 {
		conv.Reply("```В папке %q несколько файлов, не могу выбрать```", util.DestDir)
		return
	}

	archivePattern, err := regexp.Compile("^.+\\.7z$")
	if err != nil {
		conv.Reply("```Ошибка\n%s```", err)
		return
	}

	signed := DestDirFiles[0]

	if !archivePattern.MatchString(signed.Name()) {
		conv.Reply("```Файл %q не является архивом```", signed.Name())
		return
	}

	tempFile := os.Getenv("TEMP") + "\\" + signed.Name()
	util.CopyFile(util.DestDir+signed.Name(), tempFile)

	archive, err := unarr.NewArchive(tempFile)
	if err != nil {
		conv.Reply("```Ошибка инициализации архива %q\n%s```", tempFile, err)
		return
	}
	defer os.Remove(tempFile)
	defer archive.Close()

	if err = archive.Extract(util.DataDir); err != nil {
		conv.Reply("```Ошибка при попытке распаковки архива\n%s```", err)
		return
	}

	conv.Reply("_Установка подписанных шаблонов завершена успешно_")
}
