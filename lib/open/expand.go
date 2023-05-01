package open

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	list "github.com/antonio-leitao/nau/lib/list"
	utils "github.com/antonio-leitao/nau/lib/utils"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Styles struct {
	App    lipgloss.Style
	Header lipgloss.Style
}

var (
	verySubduedColor = lipgloss.AdaptiveColor{Light: "#DDDADA", Dark: "#3C3C3C"}
	subduedColor     = lipgloss.AdaptiveColor{Light: "#9B9B9B", Dark: "#5C5C5C"}
	highlight        = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special          = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}

	divider = lipgloss.NewStyle().
		SetString("•").
		Padding(0, 1).
		Foreground(verySubduedColor).
		String()

	url       = lipgloss.NewStyle().Foreground(special).Render
	docStyle  = lipgloss.NewStyle().Width(50).MarginTop(2)
	descStyle = lipgloss.NewStyle().MarginTop(5)

	infoStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderTop(true).
			BorderForeground(verySubduedColor).
			Foreground(subduedColor)
)

type model struct {
	list       list.Model
	list_width int
	width      int
	height     int
	templates  []string
	base_color string
	editor     string
	done       bool
	paths      string
}

func (m model) HeaderView() string {
	temps := strings.Join(m.templates, " • ")
	desc := lipgloss.JoinVertical(lipgloss.Left,
		descStyle.Render("NAU: project manager"),
		infoStyle.Render(temps),
	)
	return desc
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
		if msg.String() == "enter" && m.list.FilterState() != list.Filtering {
			m.done = true
			return m, m.Submit
		}
	case tea.WindowSizeMsg:
		//	h, v := docStyle.GetFrameSize()
		m.width = msg.Width
		m.height = int(float64(msg.Height) * 0.75)
		m.list.SetSize(m.list_width, m.height)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	var sections []string
	sections = append(sections, m.HeaderView())
	sections = append(sections, docStyle.Render(m.list.View()))
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center, lipgloss.JoinVertical(lipgloss.Center, sections...))
}
func (m model) Submit() tea.Msg {
	path := m.list.Submit()
	// Change to the specified directory
	if err := os.Chdir(path); err != nil {
		fmt.Println(err)
		return tea.Quit()
	}
	// Open Neovim
	cmd := exec.Command(m.editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	return tea.Quit()
}

// 1 Project
func Expand(config utils.Config) {
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
	//get list of templates
	keys := make([]string, 0, len(config.Templates))
	for k := range config.Templates {
		keys = append(keys, k)
	}

	delegate := list.NewDefaultDelegate()
	m := model{
		list_width: 50,
		editor:     config.Editor,
		done:       false,
		templates:  keys,
		base_color: config.Base_color,
		list:       list.New(items, delegate, 0, 0),
	}
	m.list.Title = "Projects"
	m.list.Styles.Title.Background(lipgloss.Color(config.Base_color))

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

}
