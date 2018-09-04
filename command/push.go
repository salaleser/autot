package command

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"

	unarr "github.com/gen2brain/go-unarr"
	"github.com/sbstjn/hanu"
	"salaleser.ru/autot/util"
)

var Push = func(conv hanu.ConversationInterface) {
	for {
		conv.Reply("Команда нестабильная, поэтому пока отключена." +
			"Если вы хотели отправить шаблоны в КМИС ОП, то используйте команду `!pull`" +
			" (Андрей считает, что поменять местами пуш и пул будет логичнее)")
		return
	}

	if util.Status != 1 {
		conv.Reply("```Нельзя изменять шаблоны пока служба не остановлена```")
		return
	}

	signedArchives, err := ioutil.ReadDir(util.PathSigned)
	if err != nil {
		conv.Reply("```Ошибка при попытке прочитать файлы из папки %s!\n```", util.PathSigned, err)
		return
	}
	if len(signedArchives) == 0 {
		conv.Reply("```В папке %s нет файлов```", util.PathSigned)
		return
	}

	var signed os.FileInfo
	if len(signedArchives) == 1 {
		signed = signedArchives[0]
		conv.Reply("выбран файл `%s%s`", util.PathSigned, signed.Name())
	} else {
		text := "Список подписанных архивов:\n"
		for i := 1; i <= len(signedArchives); i++ {
			n := strconv.FormatInt(int64(i), 10)
			text += n + ". " + signedArchives[i-1].Name() + "\n"
		}
		conv.Reply("```%s```", text)
		// TODO распаковать и поставить шаблоны
		conv.Reply("(`!pull <номер>` — распаковать и поставить шаблоны (в разработке))")
	}

	signedDir := util.PathSigned + signed.Name()
	_, err = exec.Command("xcopy", signedDir, util.PathTemp, "/Y").Output()
	if err != nil {
		conv.Reply("```Ошибка при копировании архива во временную папку!\n%s```", err)
		return
	}

	arcFilename := util.PathTemp + signed.Name()
	a, err := unarr.NewArchive(arcFilename)
	if err != nil {
		conv.Reply("```Ошибка инициализации архива `%s`!\n%s```", arcFilename, err)
		return
	}
	defer a.Close()

	signedFilenames, err := a.List()
	if err != nil {
		conv.Reply("```Ошибка при чтении имен файлов архива!\n%s```", err)
		return
	}
	for _, n := range signedFilenames {
		_, err := exec.Command("xcopy", util.PathData+n, util.PathBackup, "/Y").Output()
		if err != nil {
			conv.Reply("```Ошибка при попытке резервного копирования файла %s!\n%s```", n, err)
			return
		}
		os.Remove(util.PathData + n)
	}

	err = a.Extract(util.PathData)
	if err != nil {
		conv.Reply("```Ошибка при попытке распаковки архива!\n%s```", err)
		return
	}

	os.Remove(util.PathSigned + signed.Name())

	conv.Reply("_Установка подписанных шаблонов завершена успешно_")
}

var PushNthFile = func(conv hanu.ConversationInterface) {
	if util.Status != 1 {
		conv.Reply("```Нельзя изменять шаблоны пока служба не остановлена```")
		return
	}

	conv.Reply("Ничего не сделано, команда в разработке. "+
		"Временное решение: оставьте один архив в папке %s", util.PathSigned)
}
