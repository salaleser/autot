package command

import (
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"

	"github.com/sbstjn/hanu"
	"salaleser.ru/autot/util"
)

var Add = func(conv hanu.ConversationInterface) {
	s, err := conv.String("файлы,через,запятую,без,пробелов")
	if err != nil {
		conv.Reply("```Ошибка при разборе перечня файлов!\n%s```", err)
		return
	}

	newFilenames := strings.Split(s, ",")

	allTemplates, err := ioutil.ReadDir(util.PathData)
	if err != nil {
		conv.Reply("```Ошибка при попытке прочитать файлы из папки %s!\n%s```", util.PathData, err)
		return
	}

	patternTemplateWithoutExtension, err := regexp.Compile("^\\w+$")
	if err != nil {
		conv.Reply("```Ошибка!\n%s```", err)
		return
	}

	patternDatabaseWithoutExtension, err := regexp.Compile("^kmis_frmstr$")
	if err != nil {
		conv.Reply("```Ошибка!\n%s```", err)
		return
	}

	patternTemplate, err := regexp.Compile("^\\w+\\.ntf$")
	if err != nil {
		conv.Reply("```Ошибка!\n%s```", err)
		return
	}

	patternDatabase, err := regexp.Compile("^kmis_frmstr\\.nsf$")
	if err != nil {
		conv.Reply("```Ошибка!\n%s```", err)
		return
	}

	const templateExtension = ".ntf"
	const databaseExtension = ".nsf"
	lentgh := len(util.Files)
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
					count++
					key := strconv.Itoa(lentgh + count)
					util.Files[key] = newFilename
					continue
				}
				conv.Reply("```Файл %s не является шаблоном или БД ПФ!\n```", newFilename)
			}
		}
	}

	if count == 0 {
		conv.Reply("Ни один файл не прошел проверку (регистр учитывается)")
		return
	}

	if count != len(newFilenames) {
		conv.Reply("Не все файлы прошли проверку")
		return
	}

	util.SaveFileList()

	var flexion string
	if count > 1 {
		flexion = "ы"
	}
	conv.Reply("Успешно добавлен%s (`!pull` — _отправить_ в `%s`)", flexion, util.PathKmis)
}
