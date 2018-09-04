package command

import (
	"github.com/sbstjn/hanu"
	"salaleser.ru/autot/util"
)

// Hook содержит функцию, которая сохранит ссылку текущего канала в глобальную переменную
var Hook = func(conv hanu.ConversationInterface) {
	util.Conv = conv
	conv.Reply("*deprecated*\n`!hook` можно больше не использовать," +
		" так как по команде `!stop` этот канал автоматически добавляется в список оповещения")
}
