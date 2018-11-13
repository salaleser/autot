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
	var refused []string
	var duplicates []string
	for _, arg := range args {
		for _, file := range files {
			if file.IsDir() {
				continue
			}

			// Сначала надо проверять БД ПФ!
			if patternDatabaseWithoutExtension.MatchString(arg) {
				arg += ".nsf"
			} else if patternTemplateWithoutExtension.MatchString(arg) {
				arg += ".ntf"
			}

			isTemplate := patternTemplate.MatchString(arg)
			isDatabase := patternDatabase.MatchString(arg)
			lowerCasedArg := strings.ToLower(arg)
			lowerCasedFilename := strings.ToLower(file.Name())
			if lowerCasedFilename == lowerCasedArg {
				if util.ContainsFileAlready(lowerCasedArg) {
					duplicates = append(duplicates, file.Name())
					for i := 0; i < len(args); i++ {
						if args[i] != file.Name() {
							refused = append(refused, file.Name())
						}
					}
					continue
				}
				if isTemplate || isDatabase {
					key := strconv.Itoa(lk + len(accepted) + 1)
					accepted = append(accepted, file.Name())
					util.Files[key] = file.Name()
					for i := 0; i < len(args); i++ {
						if args[i] != file.Name() {
							refused = append(refused, file.Name())
						}
					}
					continue
				}
			}
		}
	}
	util.UpdateBackupFile()

	if len(accepted) == 0 {
		text := fmt.Sprintf("Уже были добавлены ранее: %v\n", duplicates)
		poster.PostWarning(channel, "Ни один новый файл не добавлен!", text, "")
		return
	}

	if len(accepted) != len(args) {
		text := fmt.Sprintf("Добавлены сейчас: %v\nУже были добавлены ранее: %v",
			accepted, duplicates)
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

func expandCore(a []string) []string {
	var args []string
	for i := 1; i < len(a); i++ {
		args = append(args, a[i])
		if a[i] == "core" {
			args = append(args, "MKAmbul2.ntf")
			args = append(args, "MKAmbul2M.ntf")
			args = append(args, "MKCurrent2.ntf")
			args = append(args, "MKArhiv2.ntf")
		}
	}
	return args
}
