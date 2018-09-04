package command

import (
	"time"

	"github.com/sbstjn/hanu"
	"salaleser.ru/autot/util"
)

const (
	rpsRock      = ":fist: *КАМЕНЬ*"
	rpsPaper     = ":hand: *БУМАГА*"
	rpsScissors  = ":v: *НОЖНИЦЫ*"
	rpsCountdown = 15
)

var (
	rpsEnabled bool
	rpsActions = []string{
		"_смотрит сурово..._",
		"_хмурит брови..._",
		"_напрягся..._",
		"_уверенно стоит на ногах..._",
		"_постукивает ногой в такт музыки..._",
		"_затаился в укрытии..._",
		"_готовится к выпаду..._",
		"_внимательно рассматривает свою ладонь..._",
		"_сосредоточенно разглядывает пиксель..._",
		"_сдержанно улыбается..._",
		"_ожидает результатов..._",
		"_задумчиво смотрит в даль..._",
		"_заметно нервничает..._",
		"_разглядывает пятно на полу..._",
		"_выглядит гордым..._",
	}
)

var Rock = func(conv hanu.ConversationInterface) {
	user := conv.Message().User()
	playRps(user, rpsRock)
}

var Paper = func(conv hanu.ConversationInterface) {
	user := conv.Message().User()
	playRps(user, rpsPaper)
}

var Scissors = func(conv hanu.ConversationInterface) {
	user := conv.Message().User()
	playRps(user, rpsScissors)
}

func playRps(user string, rps string) {
	username, ok := util.Users[user]
	if !ok {
		username = user
	}

	if !rpsEnabled {
		util.Players = make(map[string]string)
		rpsEnabled = true
		util.Conv.Reply("%s начал состязание! %s секунд на ответ! (`!r`, `!p` или `!s`)", username,
			rpsCountdown)
		go startRps()
	}
	util.Players[username] = rps
}

func startRps() {
	for i := rpsCountdown; i > 0; i-- {
		time.Sleep(time.Second)
	}
	util.Conv.Reply("Состязание завершено! Убрать шпаги в ножны!")

	resultMessage := "*Итоги состязания:*\n"
	var playersList []string
	for player, rps := range util.Players {
		playersList = append(playersList, player)
		resultMessage += player + " _выбрал_ " + rps + "\n"
	}
	util.Conv.Reply(resultMessage, "\n_Идет подсчет результатов..._")

	for i := 3; i > 0; i-- {
		time.Sleep(time.Second)
		randomPlayer := util.RandomNumber(len(playersList))
		randomAction := util.RandomNumber(len(rpsActions))
		util.Conv.Reply(playersList[randomPlayer], rpsActions[randomAction])
	}
	util.Conv.Reply("_А впрочем, считайте сами, я все равно пока сам не умею_ :man-shrugging:")

	rpsEnabled = false
}
