package todo

import (
	"fmt"
	utils "github.com/antonio-leitao/nau/lib/utils"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

func Todo(config utils.Config) {
	//add the new todo to the file
	model := initialTodoModel()
	if _, err := tea.NewProgram(model).Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}
func Todos(config utils.Config, query string) {
	// search for find
    //read all todos
    //separate into queried and not
}
