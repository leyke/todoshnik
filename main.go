package main

import (
	"errors"

	"fmt"
	"sort"
)

type Task struct {
	ID    int
	Title string
	Done  bool
}

type TaskManager struct {
	tasks  map[int]*Task
	nextID int
}

func (tm *TaskManager) AddTask(title string) *Task {
	newTask := &Task{
		ID:    tm.nextID,
		Title: title,
		Done:  false,
	}

	tm.tasks[tm.nextID] = newTask
	tm.nextID += 1

	return newTask
}

func (tm *TaskManager) ListTasks() []*Task {
	keys := make([]int, 0, len(tm.tasks))
	for k := range tm.tasks {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	result := make([]*Task, 0, len(keys))
	for _, k := range keys {
		result = append(result, tm.tasks[k])
	}
	return result
}

func (tm *TaskManager) DeleteTask(taskId int) error {
	_, ok := tm.tasks[taskId]
	if ok {
		delete(tm.tasks, taskId)
		return nil
	}
	return errors.New("Задача не найдена")
}

func (tm *TaskManager) MarkDone(taskId int) (*Task, error) {
	_, ok := tm.tasks[taskId]
	if ok {
		task := tm.tasks[taskId]
		task.Done = true
		return task, nil
	}

	return nil, errors.New("Задача не найдена")
}

func main() {
	tm := TaskManager{tasks: make(map[int]*Task), nextID: 1}
	tm.AddTask("first")
	printList(tm.ListTasks())

	tm.AddTask("second")
	tm.MarkDone(1)
	printList(tm.ListTasks())

	fmt.Println(tm.DeleteTask(14))
	tm.DeleteTask(1)
	printList(tm.ListTasks())
}

func printList(list []*Task) {
	for _, task := range list {
		fmt.Printf("ID: %d, Title: %s, Done: %v\n", task.ID, task.Title, task.Done)
	}
}
