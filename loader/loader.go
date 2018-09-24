package loader

import (
	"github.com/sbstjn/hanu"
	"salaleser.ru/autot/command"
)

// LoadCommands загружает модули (подключает к библиотеке hanu команды)
func LoadCommands() []hanu.Command {
	var name string
	var description string
	var function func(conv hanu.ConversationInterface)
	var commands []hanu.Command

	name = "!autot"
	description = "_сделать красиво_ (остановит, запакует архив, копирует в КМИС ОП и запустит)"
	function = command.Autot
	commands = append(commands, hanu.NewCommand(name, description, function))

	name = "!query"
	description = "проверяет состояние службы"
	function = command.Query
	commands = append(commands, hanu.NewCommand(name, description, function))

	name = "!stop"
	description = "останавливает службу"
	function = command.Stop
	commands = append(commands, hanu.NewCommand(name, description, function))

	name = "!start"
	description = "запускает службу"
	function = command.Start
	commands = append(commands, hanu.NewCommand(name, description, function))

	name = "!ver"
	description = "информация об этом боте"
	function = command.Ver
	commands = append(commands, hanu.NewCommand(name, description, function))

	name = "!aliases"
	description = "показывает список алиасов шаблонов"
	function = command.Aliases
	commands = append(commands, hanu.NewCommand(name, description, function))

	name = "!config"
	description = "показать настройки"
	function = command.Config
	commands = append(commands, hanu.NewCommand(name, description, function))

	name = "!config <key> <value>"
	description = "изменяет значение ключа `<key>` на `<value>`"
	function = command.ConfigReplaceValue
	commands = append(commands, hanu.NewCommand(name, description, function))

	name = "!config reload"
	description = "перезагрузить конфиг из файла"
	function = command.ConfigReload
	commands = append(commands, hanu.NewCommand(name, description, function))

	name = "!add <файлы,через,запятую,без,пробелов>"
	description = "обновляет список отправляемых файлов"
	function = command.Add
	commands = append(commands, hanu.NewCommand(name, description, function))

	name = "!rm <номер>"
	description = "удаляет файл из списка отправляемых по номеру"
	function = command.Rm
	commands = append(commands, hanu.NewCommand(name, description, function))

	name = "!clear"
	description = "обнуляет список файлов"
	function = command.Clear
	commands = append(commands, hanu.NewCommand(name, description, function))

	name = "!status"
	description = "показывает список отправляемых файлов"
	function = command.Status
	commands = append(commands, hanu.NewCommand(name, description, function))

	name = "!pull"
	description = "_отправляет_ подготовленные файлы"
	function = command.Pull
	commands = append(commands, hanu.NewCommand(name, description, function))

	name = "!push"
	description = "найдет файл в папке *Подписанные*, сохранит резервную копию старых файлов из " +
		"папки *dominodata* в папку *dominodata\\backup* и распакует с заменой файлы из архива в " +
		"папку *dominodata* (работает нестабильно, лучше вручную пока ставить)"
	function = command.Push
	commands = append(commands, hanu.NewCommand(name, description, function))

	name = "\\-"
	description = "отменить запланированную остановку службы (правильно не `\\-`, а `-`)"
	function = command.VoteNegative
	commands = append(commands, hanu.NewCommand(name, description, function))

	name = "!ping"
	description = "отправляет личное сообщение начальнику ОП о шаблонах"
	function = command.Ping
	commands = append(commands, hanu.NewCommand(name, description, function))

	return commands
}
