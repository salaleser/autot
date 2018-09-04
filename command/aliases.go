package command

import (
	"strings"

	"github.com/sbstjn/hanu"
	"salaleser.ru/autot/util"
)

// Aliases содержит функцию, которая отображает известные алиасы шаблонов.
// Эти алиасы содержатся в файле "aliases.list"
var Aliases = func(conv hanu.ConversationInterface) {
	text := "Список алиасов шаблонов:\n"
	columnWidth := 26
	for filename, alias := range util.Aliases {
		spaces := columnWidth - len(filename)
		if spaces < 1 {
			spaces = 1
		}
		text += filename + strings.Repeat(" ", spaces) + alias + "\n"
	}
	conv.Reply("```%s```", text)
}
