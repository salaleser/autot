package command

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/nlopes/slack"
	"salaleser.ru/autot/poster"
	"salaleser.ru/autot/util"
)

// AddHandler добавляет указанные шаблоны в список отправляемых файлов
func AddHandler(c *slack.Client, rtm *slack.RTM, ev *slack.MessageEvent, data []string) {
	channel := ev.Channel
	if len(data) < 2 {
		poster.PostError(channel, "Не указано файлов!", "Укажите файлы через пробелы")
		return
	}

	newFilenames := data[1:]

	allTemplates, err := ioutil.ReadDir(util.DataDir)
	if err != nil {
		poster.PostError(channel, fmt.Sprintf("Ошибка при попытке прочитать файлы из папки %s!",
			util.DataDir), err.Error())
		return
	}

	patternTemplateWithoutExtension, err := regexp.Compile("^\\w+$")
	if err != nil {
		poster.PostError(channel, "Ошибка!", err.Error())
		return
	}

	patternDatabaseWithoutExtension, err := regexp.Compile("^kmis_frmstr$")
	if err != nil {
		poster.PostError(channel, "Ошибка!", err.Error())
		return
	}

	patternTemplate, err := regexp.Compile("^\\w+\\.ntf$")
	if err != nil {
		poster.PostError(channel, "Ошибка!", err.Error())
		return
	}

	patternDatabase, err := regexp.Compile("^kmis_frmstr\\.nsf$")
	if err != nil {
		poster.PostError(channel, "Ошибка!", err.Error())
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

	var acceptedTemplates []string
	var refusedTemplates []string
	var duplicates []string
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
					duplicates = append(duplicates, templateFile.Name())
					continue
				}
				if isTemplate || isDatabase {
					key := strconv.Itoa(lk + len(acceptedTemplates) + 1)
					acceptedTemplates = append(acceptedTemplates, templateFile.Name())
					util.Files[key] = templateFile.Name()
					continue
				}

				refusedTemplates = append(refusedTemplates, templateFile.Name())

				text := fmt.Sprintf("Файл %q не является шаблоном или БД ПФ!", newFilename)
				poster.PostWarning(channel, "", text, "")
			}
		}
	}
	util.UpdateBackupFile()

	if len(acceptedTemplates) == 0 {
		text := "Уже были добавлены ранее: " + fmt.Sprintf("%v\n", duplicates) +
			"Не прошли проверку: " + fmt.Sprintf("%v", refusedTemplates)
		poster.PostWarning(channel, "Ни один новый файл не добавлен!", text, "")
		return
	}

	if len(acceptedTemplates) != len(newFilenames) {
		text := "Добавлены: " + fmt.Sprintf("%v\n", acceptedTemplates) +
			"Уже были добавлены ранее: " + fmt.Sprintf("%v\n", duplicates) +
			"Не прошли проверку: " + fmt.Sprintf("%v", refusedTemplates)
		poster.PostWarning(channel, "Не все файлы добавлены!", text, "")
		return
	}

	var flexion string
	if len(acceptedTemplates) > 1 {
		flexion = "ы"
	}

	title := fmt.Sprintf("Успешно добавлен%s", flexion)
	text := fmt.Sprintf("Добавлены: %v", acceptedTemplates)
	footer := fmt.Sprintf("`!pull` — _отправить_ в *%s*, `!autot` — автоотправка", util.SrcDir)
	poster.Post(channel, title, text, footer)
}
