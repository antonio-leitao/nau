package todo

import (
	"bufio"
	"fmt"
	utils "github.com/antonio-leitao/nau/lib/utils"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"regexp"
	"strings"
)

func Todo(config utils.Config) {
	//initialize at least once
	//add the new todo to the file
	model := initialTodoModel()
	if _, err := tea.NewProgram(model).Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}

}

// ----------READING AND DELETING TODOS---------------------

func parseMemoString(s string) Memo {
	re := regexp.MustCompile(`^\s*(([A-Z]{3})\s?:\s*)?(.*)$`)
	matches := re.FindStringSubmatch(s)

	title := ""
	description := strings.TrimSpace(matches[3])
	if matches[2] != "" {
		title = matches[2]
	}

	return Memo{Title: title, Description: description, Style: NewDefaultItemStyles()}
}
func readTodosFile() ([]Memo, error) {
	//open todo file and parse memos.
	filePath, _ := utils.ConvertPath(".nau/todos")
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	//each line is guarenteed to be a different memo
	//that is guarenteed by previous model
	memos := []Memo{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		memos = append(memos, parseMemoString(line))
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return memos, nil
}

func Todos(config utils.Config, query string) {
	//read all todos
	memos, err := readTodosFile()
	if err != nil {
		fmt.Println(err)
		fmt.Printf("could not read Todos: %s\n", err)
		os.Exit(1)
	}
	//separate into queried and not
	model := New(
		memos, query, config.Base_color,
	)
	if _, err := tea.NewProgram(model, tea.WithAltScreen()).Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}
