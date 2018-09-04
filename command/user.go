package command

import "github.com/sbstjn/hanu"

// User содержит функцию, которая отобразит слэк-идентификатор запросившего пользователя
var User = func(conv hanu.ConversationInterface) {
	user := conv.Message().User()
	conv.Reply("`%s`", user)
}
