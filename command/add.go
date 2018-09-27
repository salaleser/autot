package command

import (
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/sbstjn/hanu"
	"salaleser.ru/autot/util"
)

// Add содержит функцию, которая добавляет указанные шаблоны в список отправляемых файлов
var Add = func(conv hanu.ConversationInterface) {
	s, err := conv.String("файлы,через,запятую,без,пробелов")
	if err != nil {
		conv.Reply("```Ошибка при разборе перечня файлов!\n%s```", err)
		return
	}

	newFilenames := strings.Split(s, ",")

	allTemplates, err := ioutil.ReadDir(util.DataDir)
	if err != nil {
		conv.Reply("```Ошибка при попытке прочитать файлы из папки %s!\n%s```", util.DataDir, err)
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

	// Здесь нужно не количество элементов, а номер последнего ключа (пропуская несуществующие)
	var lk int
	for k := range util.Files {
		ik, err := strconv.Atoi(k)
		if err != nil {
			log.Println(err)
		}
		if ik > lk {
			lk = ik
		}
	}

	var count int
	for _, newFilename := range newFilenames {
		for _, templateFile := range allTemplates {
			if templateFile.IsDir() {
				continue
			}

			// Сначала надо проверять БД ПФ!
			if patternDatabaseWithoutExtension.MatchString(newFilename) {
				newFilename += ".nsf"
			} else if patternTemplateWithoutExtension.MatchString(newFilename) {
				newFilename += ".ntf"
			}

			isTemplate := patternTemplate.MatchString(newFilename)
			isDatabase := patternDatabase.MatchString(newFilename)
			lowerCasedTemplate := strings.ToLower(templateFile.Name())
			lowerCasedNewFilename := strings.ToLower(newFilename)
			if lowerCasedTemplate == lowerCasedNewFilename {
				if util.ContainsFileAlready(lowerCasedNewFilename) {
					continue
				}
				if isTemplate || isDatabase {
					count++
					key := strconv.Itoa(lk + count)
					util.Files[key] = templateFile.Name()
					continue
				}
				conv.Reply("```Файл %q не является шаблоном или БД ПФ!\n```", newFilename)
			}
		}
	}

	if count == 0 {
		conv.Reply("Ни один файл не прошел проверку")
		return
	}

	if count != len(newFilenames) {
		conv.Reply("Не все файлы прошли проверку")
		return
	}

	util.UpdateBackupFile()

	var flexion string
	if count > 1 {
		flexion = "ы"
	}
	conv.Reply("Успешно добавлен%s (`!pull` — _отправить_ в `%s`)", flexion, util.SrcDir)
}
