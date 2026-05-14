package cli

import (
	"errors"
	"strconv"
)

func getIntFromArgs(args []string, index int) (int, error) {
	if len(args) <= index {
		return 0, errors.New("Недостаточно аргументов")
	}
	return strconv.Atoi(args[index])
}
