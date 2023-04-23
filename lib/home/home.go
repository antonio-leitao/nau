package home

import (
	"fmt"
	"os"
	"sort"

	list "github.com/antonio-leitao/nau/lib/list"
	utils "github.com/antonio-leitao/nau/lib/utils"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type model struct {
	list list.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

type CustomDelegate struct {
	list.DefaultDelegate
}

// 1 Project
func Home(config utils.Config) {
	//Run open and open project
	projects, err := utils.GetProjects(config)
	if err != nil {
		fmt.Println("Home error:", err)
		os.Exit(1)
	}

	// define the custom Less function
	lessFunc := func(i, j int) bool {
		return projects[i].Timestamp.After(projects[j].Timestamp)
	}

	// sort the list using the Less function
	sort.Slice(projects, lessFunc)
	items := make([]list.Item, len(projects))
	for i, p := range projects {
		items[i] = p
	}

	m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "Projects"
	m.list.Styles.Title.Background(lipgloss.Color(config.Base_color))

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

}
