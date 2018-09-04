package command

import "github.com/sbstjn/hanu"

var User = func(conv hanu.ConversationInterface) {
	user := conv.Message().User()
	conv.Reply("`%s`", user)
}
