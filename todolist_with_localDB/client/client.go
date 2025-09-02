package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Todo struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

const baseURL = "http://localhost:8080/todos" // adjust to your server

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n==== TODO LIST CLIENT ====")
		fmt.Println("1. List Todos")
		fmt.Println("2. Get one Todo")
		fmt.Println("3. Add Todo")
		fmt.Println("4. Toggle Todo")
		fmt.Println("5. Update title")
		fmt.Println("6. Delete Todo")
		fmt.Println("7. Exit")
		fmt.Print("Choose an option: ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			listTodos()
		case "2":
			getOneTodo(reader)
		case "3":
			addTodo(reader)
		case "4":
			toggleTodo(reader)
		case "5":
			updateTitle(reader)
		case "6":
			deleteTodo(reader)
		case "7":
			fmt.Println("Bye 👋")
			return
		default:
			fmt.Println("Invalid choice")
		}
	}
}

func listTodos() {
	resp, err := http.Get(baseURL)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	var todos []Todo
	if err := json.NewDecoder(resp.Body).Decode(&todos); err != nil {
		fmt.Println("Error decoding:", err)
		return
	}

	if len(todos) == 0 {
		fmt.Println("No todos yet!")
		return
	}

	for _, t := range todos {
		printDataRespond(&t)
	}
}

func getOneTodo(reader *bufio.Reader) {
	fmt.Print("Enter Todo ID: ")
	idStr, _ := reader.ReadString('\n')
	idStr = strings.TrimSpace(idStr)

	url := fmt.Sprintf("%s/%s", baseURL, idStr)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: status code", resp.StatusCode)
		return
	}

	decodeDataRespond(resp)
}

func addTodo(reader *bufio.Reader) {
	fmt.Print("Enter todo title: ")
	title, _ := reader.ReadString('\n')
	title = strings.TrimSpace(title)

	body, _ := json.Marshal(map[string]string{"title": title})
	fmt.Printf("%s\n", body)

	resp, err := http.Post(baseURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusOK {
		fmt.Println("Todo added ✅")
	} else {
		fmt.Println("Failed to add todo:", resp.Status)
	}

	decodeDataRespond(resp)
}

func toggleTodo(reader *bufio.Reader) {
	fmt.Print("Enter todo ID to toggle: ")
	idStr, _ := reader.ReadString('\n')
	idStr = strings.TrimSpace(idStr)
	id, _ := strconv.Atoi(idStr)

	// First, fetch the todo
	resp, err := http.Get(fmt.Sprintf("%s/%d", baseURL, id))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Todo not found")
		return
	}

	var todo Todo
	if err := json.NewDecoder(resp.Body).Decode(&todo); err != nil {
		fmt.Println("Error decoding:", err)
		return
	}

	// Toggle Completed
	body, _ := json.Marshal(map[string]any{"completed": !todo.Completed})

	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/%d", baseURL, id), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		fmt.Println("Todo toggled 🔄")
	} else {
		fmt.Println("Failed to toggle:", res.Status)
	}

	decodeDataRespond(resp)
}

func updateTitle(reader *bufio.Reader) {
	fmt.Print("Enter todo ID to update title: ")
	idStr, _ := reader.ReadString('\n')
	idStr = strings.TrimSpace(idStr)
	id, _ := strconv.Atoi(idStr)

	// First, fetch the todo
	resp, err := http.Get(fmt.Sprintf("%s/%d", baseURL, id))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Todo not found")
		return
	}

	var todo Todo
	if err := json.NewDecoder(resp.Body).Decode(&todo); err != nil {
		fmt.Println("Error decoding:", err)
		return
	}

	// Update title
	fmt.Print("Enter title to update: ")
	enterTitle, _ := reader.ReadString('\n')
	enterTitle = strings.TrimSpace(enterTitle)

	body, _ := json.Marshal(map[string]string{"title": enterTitle})

	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/%d", baseURL, id), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		fmt.Println("Title updated ✅")
	} else {
		fmt.Println("Failed to update:", res.Status)
	}

	decodeDataRespond(resp)
}

func deleteTodo(reader *bufio.Reader) {
	fmt.Print("Enter todo ID to delete: ")
	idStr, _ := reader.ReadString('\n')
	idStr = strings.TrimSpace(idStr)
	id, _ := strconv.Atoi(idStr)

	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%d", baseURL, id), nil)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK || res.StatusCode == http.StatusNoContent {
		fmt.Println("Todo deleted ❌")
	} else {
		fmt.Println("Failed to delete:", res.Status)
	}
}

func decodeDataRespond(resp *http.Response) {
	var todo Todo
	if err := json.NewDecoder(resp.Body).Decode(&todo); err != nil {
		fmt.Println("Error decoding:", err)
		return
	}

	printDataRespond(&todo)
}

func printDataRespond(todo *Todo) {
	status := " "
	if todo.Completed {
		status = "✅"
	}

	fmt.Printf("[%s] %d: Title: %s\n", status, todo.ID, todo.Title)
	fmt.Printf("CreatedAt: %s\n", todo.CreatedAt)
	fmt.Printf("UpdatedAt: %s\n", todo.UpdatedAt)
}
