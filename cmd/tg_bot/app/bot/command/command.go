package command

type Name string

const (
	CommandStart      Name = "start"
	CommandRestart    Name = "restart"
	CommandHelp       Name = "help"
	CommandStatus     Name = "status"
	CommandAdd        Name = "add"
	CommandTaskDone   Name = "taskdone"
	CommandTaskDelete Name = "taskdelete"
	CommandTaskList   Name = "tasklist"
)

type CommandDto struct {
	Name
	State
}
