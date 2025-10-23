package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

const path = "tasks.json"

type Task struct {
	id          int
	description string
	status      string
	createdAt   time.Time
	updatedAt   time.Time
}

func NewTask(id int, description string) *Task {
	return &Task{
		id:          id,
		description: description,
		status:      "todo",
		createdAt:   time.Now(),
		updatedAt:   time.Now(),
	}
}

func logging(id int, action string) {
	fmt.Println(fmt.Sprintf("Task %s successfully: (ID: %d)", action, id))
}

func getNewUniqId(data []Task) int {
	if len(data) == 0 {
		return 1
	}
	newId := data[len(data)-1].id + 1
	return newId
}

func Add(description string) error {
	tasks, err := readTasks()
	if err != nil {
		return err
	}
	id := getNewUniqId(tasks)
	tasks = append(tasks, *NewTask(id, description))
	err = writeTasks(tasks)
	if err != nil {
		return err
	}
	logging(id, "added")
	return nil
}

func Update(id int, description string) error {
	tasks, err := readTasks()
	if err != nil {
		return err
	}

	for i, v := range tasks {
		if id == v.id {
			tasks[i].description = description
			tasks[i].updatedAt = time.Now()
		}
	}

	err = writeTasks(tasks)
	if err != nil {
		return err
	}
	logging(id, "updated")
	return nil
}

func Delete(id int) error {
	tasks, err := readTasks()
	if err != nil {
		return err
	}

	for i, v := range tasks {
		if id == v.id {
			tasks = append(tasks[:i], tasks[i+1:]...)
		}
	}

	err = writeTasks(tasks)
	if err != nil {
		return err
	}
	logging(id, "deleted")
	return nil
}

func markStatus(id int, status string) error {
	tasks, err := readTasks()
	if err != nil {
		return err
	}
	for i, v := range tasks {
		if id == v.id {
			tasks[i].status = status
		}
	}
	err = writeTasks(tasks)
	if err != nil {
		return err
	}
	return nil
}

func markInProgress(id int) error {
	err := markStatus(id, "in-progress")
	if err != nil {
		return err
	}
	logging(id, "marked `in-progress`")
	return nil
}

func markDone(id int) error {
	err := markStatus(id, "done")
	if err != nil {
		return err
	}
	logging(id, "marked `done`")
	return nil
}

func getListTasksByStatus(status string) ([]Task, error) {
	tasks, err := readTasks()
	if err != nil {
		return nil, err
	}
	if status == "" {
		return tasks, nil
	}
	var res []Task
	for _, v := range tasks {
		if status == v.status {
			res = append(res, v)
		}
	}
	return res, nil
}

func readTasks() ([]Task, error) {
	data, err := readJSON()
	if err != nil {
		return nil, err
	}
	tasks, err := convertInTasks(data)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func writeTasks(tasks []Task) error {
	err := writeJSON(convertTasksInJson(tasks))
	if err != nil {
		return err
	}
	return nil
}

func convertInTasks(tasks map[string][]map[string]interface{}) ([]Task, error) {
	var res []Task

	for _, v := range tasks["tasks"] {
		createdAt, err := time.Parse(time.RFC3339, v["createdAt"].(string))
		if err != nil {
			return nil, err
		}
		updatedAt, err := time.Parse(time.RFC3339, v["updatedAt"].(string))
		if err != nil {
			return nil, err
		}
		res = append(res, Task{
			id:          int(v["id"].(float64)),
			description: v["description"].(string),
			status:      v["status"].(string),
			createdAt:   createdAt,
			updatedAt:   updatedAt,
		})
	}

	return res, nil
}

func convertTasksInJson(tasks []Task) map[string][]map[string]interface{} {
	var res = make(map[string][]map[string]interface{})
	var temp []map[string]interface{}

	for _, v := range tasks {
		temp = append(temp, map[string]interface{}{
			"id":          v.id,
			"description": v.description,
			"status":      v.status,
			"createdAt":   v.createdAt,
			"updatedAt":   v.updatedAt,
		})
	}

	res["tasks"] = temp

	return res
}

func writeJSON(tasks map[string][]map[string]interface{}) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func readJSON() (map[string][]map[string]interface{}, error) {
	file, err := os.Open(path)
	var tasks map[string][]map[string]interface{}

	if os.IsNotExist(err) {
		err = os.WriteFile(path, []byte("{}"), 0644)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	} else {
		defer file.Close()

		decoder := json.NewDecoder(file)
		err = decoder.Decode(&tasks)
		if err != nil {
			return nil, err
		}
	}
	return tasks, nil
}

func main() {
	command := os.Args[1]

	switch command {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("To less arguments, add <description>")
			return
		}
		err := Add(os.Args[2])
		if err != nil {
			fmt.Println(err)
			return
		}
	case "update":
		if len(os.Args) < 4 {
			fmt.Println("To less arguments, update <id> <description>")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println(err)
			return
		}
		err = Update(id, os.Args[3])
		if err != nil {
			fmt.Println(err)
			return
		}
	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("To less arguments, delete <id>")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println(err)
			return
		}
		err = Delete(id)
		if err != nil {
			fmt.Println(err)
			return
		}
	case "mark-in-progress":
		if len(os.Args) < 3 {
			fmt.Println("To less arguments, mark-in-progress <id>")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println(err)
			return
		}
		err = markInProgress(id)
		if err != nil {
			fmt.Println(err)
			return
		}
	case "mark-done":
		if len(os.Args) < 3 {
			fmt.Println("To less arguments, mark-done <id>")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println(err)
			return
		}
		err = markDone(id)
		if err != nil {
			fmt.Println(err)
			return
		}
	case "list":
		var (
			list []Task
			err  error
		)

		if len(os.Args) < 3 {
			list, err = getListTasksByStatus("")
			if err != nil {
				fmt.Println(err)
				return
			}
		} else {
			list, err = getListTasksByStatus(os.Args[2])
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		for _, v := range list {
			fmt.Println(v)
		}
	default:
		fmt.Println("Unknown command:", command)
	}
}
