package command

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/nlopes/slack"
	"salaleser.ru/autot/gui"
	"salaleser.ru/autot/util"
)

// PullHandler содерижт функцию, которая упакует файлы и отправит в папку path-kmis
func PullHandler(c *slack.Client, rtm *slack.RTM, ev *slack.MessageEvent, data []string) {
	if util.Status != util.StatusStopped {
		errorParams := slack.PostMessageParameters{}
		attachment := slack.Attachment{
			Color: gui.Red,
			Title: "Нельзя изменять шаблоны пока служба не остановлена",
		}
		errorParams.Attachments = []slack.Attachment{attachment}
		util.API.PostMessage(ev.Channel, "", errorParams)
		return
	}

	if len(util.Files) == 0 {
		errorParams := slack.PostMessageParameters{}
		attachment := slack.Attachment{
			Color: gui.Red,
			Title: "Список пустой!",
			Text:  "*!add <файл(-ы)_через_пробелы>* — добавить файл(-ы)",
		}
		errorParams.Attachments = []slack.Attachment{attachment}
		util.API.PostMessage(ev.Channel, "", errorParams)
		return
	}

	time := time.Now()
	date := time.Format(util.TemplateDateFormat)
	arcName := "Templates_" + date + "_KMIS.zip"
	arcFullName := util.SrcDir + "\\" + arcName
	archiveFile, err := os.Create(arcFullName)
	if err != nil {
		errorParams := slack.PostMessageParameters{}
		attachment := slack.Attachment{
			Color: gui.Red,
			Title: "Ошибка при попытке создать архив!",
			Text:  err.Error(),
		}
		errorParams.Attachments = []slack.Attachment{attachment}
		util.API.PostMessage(ev.Channel, "", errorParams)
	}
	defer archiveFile.Close()

	zipWriter := zip.NewWriter(archiveFile)
	defer zipWriter.Close()

	for _, filename := range util.Files {
		file, err := os.Open(util.DataDir + filename)
		if err != nil {
			errorParams := slack.PostMessageParameters{}
			attachment := slack.Attachment{
				Color: gui.Red,
				Title: "Ошибка при попытке архивировать шаблоны!",
				Text:  err.Error(),
			}
			errorParams.Attachments = []slack.Attachment{attachment}
			util.API.PostMessage(ev.Channel, "", errorParams)
			return
		}
		defer file.Close()

		info, err := file.Stat()
		if err != nil {
			errorParams := slack.PostMessageParameters{}
			attachment := slack.Attachment{
				Color: gui.Red,
				Title: "Ошибка при попытке архивировать шаблоны!",
				Text:  err.Error(),
			}
			errorParams.Attachments = []slack.Attachment{attachment}
			util.API.PostMessage(ev.Channel, "", errorParams)
			return
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			errorParams := slack.PostMessageParameters{}
			attachment := slack.Attachment{
				Color: gui.Red,
				Title: "Ошибка при попытке архивировать шаблоны!",
				Text:  err.Error(),
			}
			errorParams.Attachments = []slack.Attachment{attachment}
			util.API.PostMessage(ev.Channel, "", errorParams)
			return
		}

		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			errorParams := slack.PostMessageParameters{}
			attachment := slack.Attachment{
				Color: gui.Red,
				Title: "Ошибка при попытке архивировать шаблоны!",
				Text:  err.Error(),
			}
			errorParams.Attachments = []slack.Attachment{attachment}
			util.API.PostMessage(ev.Channel, "", errorParams)
			return
		}

		_, err = io.Copy(writer, file)
		if err != nil {
			errorParams := slack.PostMessageParameters{}
			attachment := slack.Attachment{
				Color: gui.Red,
				Title: "Ошибка при попытке архивировать шаблоны!",
				Text:  err.Error(),
			}
			errorParams.Attachments = []slack.Attachment{attachment}
			util.API.PostMessage(ev.Channel, "", errorParams)
			return
		}
	}

	util.ArcFullName = arcFullName // временный ужос
	params := slack.PostMessageParameters{}
	attachment := slack.Attachment{
		Color:     gui.Green,
		TitleLink: "file://kserver/!Common/КМИС%20ОП/",
		Title:     fmt.Sprintf("Шаблоны в `%s`.", arcFullName),
	}
	params.Attachments = []slack.Attachment{attachment}
	params.AsUser = true
	util.API.PostMessage(ev.Channel, "", params)

	util.Files = map[string]string{}
	util.UpdateBackupFile()
}
