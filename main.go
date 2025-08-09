package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)


type Task struct {
	ID int `json:"id"`
	Description string `json:"description"`
	Completed bool `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at"`
}

type TodoList struct {
	Tasks []Task `json:"tasks"`
	NextID int `json:"next_id"`
	filename string
}

func NewTodoList(filename string) *TodoList {
	todo := &TodoList{
		Tasks: []Task{},
		NextID: 1,
		filename: filename,
	}
	todo.Load()
	return todo
}

func (tl *TodoList) Add(description string) {
	task := Task{
		ID: tl.NextID,
		Description: description,
		Completed: false,
		CreatedAt: time.Now(),
	}

	tl.Tasks = append(tl.Tasks, task)
	tl.NextID++
	tl.save()

	fmt.Printf("Added task: %d - %s\n", task.ID, task.Description)
}

func (tl *TodoList) List() {
	if len(tl.Tasks) == 0 {
		fmt.Println("No tasks to list")
		return
	}

	fmt.Println("Tasks:")
	fmt.Println("================")

	for _, task := range tl.Tasks {
		status := "X"
		if task.Completed {
			status = "/"
		}

		fmt.Printf("%s [%d] %s\n", status, task.ID, task.Description)

		if task.Completed && task.CompletedAt != nil {
			fmt.Printf("Completed at: %s\n", task.CompletedAt.Format(time.RFC3339))

		} else {
			fmt.Printf("Created at: %s\n", task.CreatedAt.Format(time.RFC3339))
		}
		fmt.Println("================")
	}
}

func (tl *TodoList) Complete(id int) {
	for i := range tl.Tasks {
		if tl.Tasks[i].ID == id {
			if tl.Tasks[i].Completed {
				fmt.Printf("Task %d is already completed\n", id)
				return
			}

			now := time.Now()
			tl.Tasks[i].Completed = true
			tl.Tasks[i].CompletedAt = &now
			tl.save()

			fmt.Printf("Task %d completed\n", id)
			return 
		}
	}
}

func (tl *TodoList) Delete(id int) {
	for i, task := range tl.Tasks {
		if task.ID == id {
			tl.Tasks = append(tl.Tasks[:i], tl.Tasks[i+1:]...)
			tl.save()
			fmt.Printf("Task %d deleted\n", id)
			return
		}
	}
	fmt.Printf("Task %d not found\n", id)
}


func (tl *TodoList) ListPending() {
	pending := []Task{}
	for _, task := range tl.Tasks {
		if !task.Completed {
			pending = append(pending, task)
		}
	}

	if len(pending) == 0 {
		fmt.Println("No pending tasks")
		return
	}

	fmt.Printf("\n Pending Tasks (%d): \n", len(pending))
	fmt.Println("================")

	for _, task := range pending {
		fmt.Printf("[%d] %s\n", task.ID, task.Description)
		fmt.Printf("Created at: %s\n", task.CreatedAt.Format(time.RFC3339))
		fmt.Println("================")
	}
}

func (tl *TodoList) DeleteCompleted() {
	if len(tl.Tasks) == 0 {
		fmt.Println("No tasks to delete")
		return
	}

	deleted := []Task{}


	for _, task := range tl.Tasks {
		if task.Completed {
			fmt.Printf("Deleting completed task: %d - %s\n", task.ID, task.Description)
			deleted = append(deleted, task)
			tl.Delete(task.ID)
		}
	}

	if len(deleted) == 0 {
		fmt.Println("No completed tasks to delete")
		return
	}

	fmt.Printf("\n Deleting %d completed tasks: \n", len(deleted))
	fmt.Println("================")

	for _, task := range deleted {
		fmt.Printf("Deleted: %d - %s\n", task.ID, task.Description)
	}

	fmt.Printf("Deleted %d completed tasks\n", len(deleted))
}

func (tl *TodoList) Stats() {
	total := len(tl.Tasks)
	completed := 0

	for _, task := range tl.Tasks {
		if task.Completed {
			completed++
		}
	}

	pending := total - completed

	fmt.Println("\n Tasks Stats: \n")
	fmt.Printf("Total tasks: %d\n", total)
	fmt.Printf("Completed tasks: %d\n", completed)
	fmt.Printf("Pending tasks: %d\n", pending)

	if total > 0 {
		percentage := float64(completed) / float64(total) * 100
		fmt.Printf("Completion percentage: %.2f%%\n", percentage)
	}

}

func (tl *TodoList) save() {
	data, err := json.MarshalIndent(tl, "", " ")
	if err != nil {
		log.Printf("Error saving file: %v", err)
		return
	}

	err = ioutil.WriteFile(tl.filename, data, 0644)
	if err != nil {
		log.Printf("Error saving file: %v", err)
		return
	}
}

func (tl *TodoList) Load() {
	data, err := ioutil.ReadFile(tl.filename)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, tl)
	if err != nil {
		log.Printf("Error loading file: %v", err)
	}
}


func showUsage() {
	fmt.Println("📝 Go Todo CLI Application")
	fmt.Println("==========================")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  todo -add \"Task description\"     Add a new task")
	fmt.Println("  todo -list                       List all tasks")
	fmt.Println("  todo -pending                    List pending tasks only")
	fmt.Println("  todo -complete 1                 Mark task #1 as complete")
	fmt.Println("  todo -delete 1                   Delete task #1")
	fmt.Println("  todo -delete-completed           Delete all completed tasks")
	fmt.Println("  todo -stats                      Show task statistics")
	fmt.Println("  todo -help                       Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  todo -add \"Buy groceries\"")
	fmt.Println("  todo -add \"Finish Go project\"")
	fmt.Println("  todo -list")
	fmt.Println("  todo -complete 1")
	fmt.Println("  todo -delete 2")
	fmt.Println("  todo -delete-completed")
}

func main() {
	var (
		add = flag.String("add", "", "Add a new task")
		list = flag.Bool("list", false, "List all tasks")
		pending = flag.Bool("pending", false, "List pending tasks only")
		complete = flag.Int("complete", 0, "Mark task as complete")
		delete = flag.Int("delete", 0, "Delete a task")
		deleteCompleted = flag.Bool("delete-completed", false, "Delete all completed tasks")
		stats = flag.Bool("stats", false, "Show task statistics")
		help = flag.Bool("help", false, "Show this help message")
	)

	flag.Parse()


	if *help || flag.NFlag() == 0 {
		showUsage()
		return
	}

	todoFile := "todos.json"
	todo := NewTodoList(todoFile)

	switch {
	case *add != "":
		todo.Add(*add)

	case *list:
		todo.List()

	case *pending:
		todo.ListPending()

	case *complete > 0:
		todo.Complete(*complete)

	case *delete > 0:
		todo.Delete(*delete) 

	case *deleteCompleted:
		todo.DeleteCompleted()

	case *stats:
		todo.Stats()

	default:
		fmt.Println("Invalid command. Use -help for usage information.")
		os.Exit(1)
	}
}