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

	args := expandCore(data)

	files, err := ioutil.ReadDir(util.DataDir)
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

	var accepted []string
	var duplicates []string
	var refused []string
	for arg := range args {
		var isAcceptedOrDuplicate bool
		for _, file := range files {
			if file.IsDir() {
				continue
			}

			var argWithExtension string
			if patternDatabaseWithoutExtension.MatchString(arg) { // Сначала проверять на БД ПФ!
				argWithExtension = arg + ".nsf"
			} else if patternTemplateWithoutExtension.MatchString(arg) {
				argWithExtension = arg + ".ntf"
			}

			isTemplate := patternTemplate.MatchString(argWithExtension)
			isDatabase := patternDatabase.MatchString(argWithExtension)
			lowerCasedArg := strings.ToLower(argWithExtension)
			lowerCasedFilename := strings.ToLower(file.Name())
			if lowerCasedFilename == lowerCasedArg {
				if util.ContainsFileAlready(lowerCasedArg) {
					duplicates = append(duplicates, file.Name())
					isAcceptedOrDuplicate = true
					continue
				}
				if isTemplate || isDatabase {
					key := strconv.Itoa(lk + len(accepted) + 1)
					accepted = append(accepted, file.Name())
					util.Files[key] = file.Name()
					isAcceptedOrDuplicate = true
					continue
				}
			}
		}
		if !isAcceptedOrDuplicate {
			refused = append(refused, arg)
		}
	}
	util.UpdateBackupFile()

	if len(accepted) == 0 {
		text := fmt.Sprintf("Уже были добавлены ранее: %v\nНе прошли проверку: %v",
			duplicates, refused)
		poster.PostWarning(channel, "Ни один новый файл не добавлен!", text, "")
		return
	}

	if len(accepted) != len(args) {
		text := fmt.Sprintf("Добавлены сейчас: %v\nУже были добавлены ранее: %v\n"+
			"Не прошли проверку: %v", accepted, duplicates, refused)
		poster.PostWarning(channel, "Не все файлы добавлены!", text, "")
		return
	}

	var flexion string
	if len(accepted) > 1 {
		flexion = "ы"
	}

	title := fmt.Sprintf("Успешно добавлен%s", flexion)
	text := fmt.Sprintf("Добавлены сейчас: %v", accepted)
	footer := fmt.Sprintf("`!pull` — _отправить_ в *%s*, `!autot` — автоотправка", util.SrcDir)
	poster.Post(channel, title, text, footer)
}

func hasElement(array []string, element string) bool {
	for _, e := range array {
		if element == e {
			return true
		}
	}
	return false
}

func expandCore(a []string) map[string]int {
	var args = make(map[string]int)
	for i := 1; i < len(a); i++ {
		args[a[i]] = 0
		if a[i] == "core" {
			args["MKAmbul2.ntf"] = 0
			args["MKAmbul2M.ntf"] = 0
			args["MKCurrent2.ntf"] = 0
			args["MKArhiv2.ntf"] = 0
		}
	}
	return args
}
