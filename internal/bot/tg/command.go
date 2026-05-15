package tg

type Command string

const (
	CommandStart      Command = "start"
	CommandRestart    Command = "restart"
	CommandHelp       Command = "help"
	CommandStatus     Command = "status"
	CommandAdd        Command = "add"
	CommandTaskDone   Command = "taskdone"
	CommandTaskDelete Command = "taskdelete"
	CommandTaskList   Command = "tasklist"
)
