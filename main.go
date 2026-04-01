package main

import (
	"fmt"
)

type Task struct {
	ID    int
	Title string
	Done  bool
}

func main() {
	list := []Task{}
	list = append(list, Task{1, "Первая", false}, Task{2, "Вторая", false}, Task{3, "Посмотреть фильм", false})
	fmt.Println(list)
}
