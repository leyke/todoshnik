package tg

type Command string

const (
	CommandStart      Command = "start"
	CommandRestart    Command = "restart"
	CommandHelp       Command = "help"
	CommandStatus     Command = "status"
	CommandAdd        Command = "add"
	СommandTaskDone   Command = "taskdone"
	CommandTaskDelete Command = "taskdelete"
	СommandTaskList   Command = "tasklist"
)
