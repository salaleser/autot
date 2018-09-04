package command

import "github.com/sbstjn/hanu"

var Autot = func(conv hanu.ConversationInterface) {
	conv.Reply("в разработке...")
}
