package command

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/nlopes/slack"
	"salaleser.ru/autot/gui"
	"salaleser.ru/autot/util"
)

var (
	params = slack.PostMessageParameters{}
)

// AddHandler добавляет указанные шаблоны в список отправляемых файлов
func AddHandler(c *slack.Client, rtm *slack.RTM, ev *slack.MessageEvent, data []string) {
	a := strings.Split(ev.Msg.Text, " ")

	if len(a) < 2 {
		errorParams := slack.PostMessageParameters{}
		attachment := slack.Attachment{
			Color: gui.Red,
			Title: "Не указано файлов!",
			Text:  "Укажите файлы через пробелы",
		}
		errorParams.Attachments = []slack.Attachment{attachment}
		util.API.PostMessage(ev.Channel, "", errorParams)
		return
	}

	newFilenames := a[1:]

	allTemplates, err := ioutil.ReadDir(util.DataDir)
	if err != nil {
		errorParams := slack.PostMessageParameters{}
		attachment := slack.Attachment{
			Color: gui.Red,
			Title: fmt.Sprintf("Ошибка при попытке прочитать файлы из папки %s!", util.DataDir),
			Text:  err.Error(),
		}
		errorParams.Attachments = []slack.Attachment{attachment}
		util.API.PostMessage(ev.Channel, "", errorParams)
		return
	}

	patternTemplateWithoutExtension, err := regexp.Compile("^\\w+$")
	if err != nil {
		errorParams := slack.PostMessageParameters{}
		attachment := slack.Attachment{
			Color: gui.Red,
			Title: "Ошибка!",
			Text:  err.Error(),
		}
		errorParams.Attachments = []slack.Attachment{attachment}
		util.API.PostMessage(ev.Channel, "", errorParams)
		return
	}

	patternDatabaseWithoutExtension, err := regexp.Compile("^kmis_frmstr$")
	if err != nil {
		errorParams := slack.PostMessageParameters{}
		attachment := slack.Attachment{
			Color: gui.Red,
			Title: "Ошибка!",
			Text:  err.Error(),
		}
		errorParams.Attachments = []slack.Attachment{attachment}
		util.API.PostMessage(ev.Channel, "", errorParams)
		return
	}

	patternTemplate, err := regexp.Compile("^\\w+\\.ntf$")
	if err != nil {
		errorParams := slack.PostMessageParameters{}
		attachment := slack.Attachment{
			Color: gui.Red,
			Title: "Ошибка!",
			Text:  err.Error(),
		}
		errorParams.Attachments = []slack.Attachment{attachment}
		util.API.PostMessage(ev.Channel, "", errorParams)
		return
	}

	patternDatabase, err := regexp.Compile("^kmis_frmstr\\.nsf$")
	if err != nil {
		errorParams := slack.PostMessageParameters{}
		attachment := slack.Attachment{
			Color: gui.Red,
			Title: "Ошибка!",
			Text:  err.Error(),
		}
		errorParams.Attachments = []slack.Attachment{attachment}
		util.API.PostMessage(ev.Channel, "", errorParams)
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
				warningParams := slack.PostMessageParameters{}
				attachment := slack.Attachment{
					Color: gui.Orange,
					Title: fmt.Sprintf("Файл %q не является шаблоном или БД ПФ!", newFilename),
				}
				warningParams.Attachments = []slack.Attachment{attachment}
				util.API.PostMessage(ev.Channel, "", warningParams)
			}
		}
	}
	util.UpdateBackupFile()

	if len(acceptedTemplates) == 0 {
		warningParams := slack.PostMessageParameters{}
		attachment := slack.Attachment{
			Color: gui.Orange,
			Title: "Ни один новый файл не добавлен!",
			Text: "Уже были добавлены ранее: " + fmt.Sprint(duplicates) +
				"\nНе прошли проверку: " + fmt.Sprint(refusedTemplates),
		}
		warningParams.Attachments = []slack.Attachment{attachment}
		util.API.PostMessage(ev.Channel, "", warningParams)
		return
	}

	if len(acceptedTemplates) != len(newFilenames) {
		warningParams := slack.PostMessageParameters{}
		attachment := slack.Attachment{
			Color: gui.Orange,
			Title: "Не все файлы добавлены!",
			Text: "Добавлены: " + fmt.Sprint(acceptedTemplates) +
				"\nУже были добавлены ранее: " + fmt.Sprint(duplicates) +
				"\nНе прошли проверку: " + fmt.Sprint(refusedTemplates),
		}
		warningParams.Attachments = []slack.Attachment{attachment}
		util.API.PostMessage(ev.Channel, "", warningParams)
		return
	}

	var flexion string
	if len(acceptedTemplates) > 1 {
		flexion = "ы"
	}

	attachment := slack.Attachment{
		Color:  gui.Green,
		Title:  "Все файлы успешно добавлен" + flexion,
		Text:   "Добавлены: " + fmt.Sprint(acceptedTemplates),
		Footer: fmt.Sprintf("`!pull` — _отправить_ в `%s`, `!autot` — автоотправка", util.SrcDir),
	}
	params.Attachments = []slack.Attachment{attachment}
	util.API.PostMessage(ev.Channel, "", params)
}
