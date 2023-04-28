package home

import (
	"fmt"
	utils "github.com/antonio-leitao/nau/lib/utils"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
)

type Styles struct {
	titleStyle  lipgloss.Style
	promptStyle lipgloss.Style
	sepStyle    lipgloss.Style
}

func defaultStyles(base_color string) Styles {
	verySubduedColor := lipgloss.AdaptiveColor{Light: "#DDDADA", Dark: "#3C3C3C"}
	subduedColor := lipgloss.AdaptiveColor{Light: "#9B9B9B", Dark: "#5C5C5C"}
	title_text := "#ffffd7" //230
	var s Styles
	s.titleStyle = lipgloss.NewStyle().
		Margin(1, 0, 1, 0).
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color(base_color)).
		Foreground(lipgloss.Color(title_text))
	s.sepStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderForeground(verySubduedColor).
		Foreground(subduedColor)
	s.promptStyle = lipgloss.NewStyle().Margin(0, 0, 0, 2).Foreground(lipgloss.Color(base_color))
	return s
}

type model struct {
	styles Styles
	config utils.Config
	width  int
	height int
}

func newModel(config utils.Config) model {
	m := model{
		styles: defaultStyles(config.Base_color),
		config: config,
	}
	return m
}

func (m model) Init() tea.Cmd {
	return nil
}
func (m model) View() string {
	output := lipgloss.JoinHorizontal(
		lipgloss.Center,
		renderBigArt(),
		m.renderConfig(),
	)
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		output,
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	//quit if the user click anywhere
	case tea.KeyMsg:
		return m, tea.Quit
	}
	return m, nil
}

func renderBigArt() string {
	s := `⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢠⣧⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣄⠀⠀⢀⣿⣤⡄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣿⠀⠈⠉⢻⡆⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⢀⡠⣤⣴⡶⠿⢿⠀⠀⠀⢸⢻⡀⠀⠀⠀⠀⠀⣀⣠⣤⡄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠰⠛⠋⠉⠀⠀⠀⢸⡀⠀⠀⢸⣀⣧⣤⠴⠒⠚⠛⢻⡍⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⢀⣀⣼⠷⠖⠚⠉⠉⢹⡀⠀⠀⠀⠀⠘⣇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠰⢞⡟⠉⣁⣠⠖⠉⠀⠀⠀⠀⢷⠀⠀⠀⠀⡀⠘⡆⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⢀⡞⣠⣾⠟⠁⠀⣀⣀⣀⣀⣀⣈⣧⣰⣺⠿⣿⣿⣿⣆⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⢸⣽⣫⣥⣴⣶⣾⡟⠛⣉⣿⣇⣀⣽⣦⣤⠴⠚⠛⠛⢻⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⣾⣿⣻⠟⢻⣁⣸⣧⣴⣿⣿⣿⠿⣿⣿⡄⠀⠀⠀⠀⠈⢷⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⢀⣀⣿⡮⠶⠶⠛⣛⡿⣿⠟⠋⠀⣿⣧⣿⠹⣵⡄⠀⢀⣤⣽⣾⣷⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⣀⡈⢩⡏⢀⡠⠾⠛⠁⠀⠀⠀⠀⢠⡏⠘⣿⠀⢳⡽⡿⠿⠿⢿⣿⣷⣿⣶⠾⠇⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⢹⡟⠰⠋⢀⣀⣤⣤⣶⡶⠿⠟⠛⡇⠀⡏⣇⠀⠀⠻⡶⠛⠛⠋⠁⠀⣌⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠘⣧⣶⣾⣿⣻⣿⣿⣯⣤⣤⣴⣾⡇⠀⡇⠉⠀⠀⠀⠹⣄⠀⠀⠀⠀⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀
⣀⠀⠀⢠⣿⣿⠿⣿⣿⣿⣿⣿⣿⣿⣿⣾⡇⠀⡇⠀⠀⠀⠀⠀⠙⣆⠀⠀⢸⠇⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠹⡟⠶⣟⣻⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠀⠀⢸⠀⠀⠀⠀⠀⠀⠘⣷⡀⢸⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⣷⠀⠈⠹⡟⠛⠛⣿⠛⢛⡟⠛⣿⢿⣿⡄⠀⠸⡄⠀⠀⠀⠀⠀⠀⠈⢟⣮⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠘⡆⠀⠀⢹⠀⠀⡏⠀⣼⠇⠀⡇⠀⢸⡇⠀⠀⣧⠀⢸⠀⠀⠀⠀⠀⠈⢿⡳⣄⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⢳⡀⠀⠈⣇⢸⠀⠀⣿⠀⠀⡇⠀⣼⢧⢸⡆⠘⣆⡞⣤⠶⠶⣶⣤⣤⣤⣹⡌⠳⣄⡀⠀⠀⠀⠀⠀
⠀⠀⠘⣇⠀⠀⢹⣸⠀⠀⡟⠀⠀⣧⣼⣿⢸⣹⣷⠈⠉⠉⢻⣷⣶⠶⠾⠽⣿⣿⣿⣦⠀⠙⢦⣀⣀⣤⠆
⠀⠀⠀⠸⡄⢠⣿⣏⡇⠀⡇⠀⣸⣿⣿⣿⣌⣿⣇⣀⠀⠀⢸⣿⣿⣷⣄⠀⠀⠀⠀⠉⠉⣹⣶⠟⠉⠀⠀
⠀⠀⠀⠀⢳⣸⣽⢻⣧⠀⣷⢀⣷⢿⣿⣿⡇⠀⣆⣿⣿⣿⡞⠛⠿⠿⠿⣷⣄⣀⣠⣴⠿⠋⠀⠀⠀⠀⠀
⣠⣀⣀⣀⡼⢯⣿⢺⣿⡀⢹⠘⠛⢸⣿⣿⡇⣼⡟⢿⣻⣿⣿⣆⠀⠀⢀⡠⣿⣿⠟⠁⠀⠀⠀⠀⠀⠀⠀
⠙⣯⡻⢷⣤⣼⣧⠈⣿⣧⠘⡆⠀⣾⣿⣿⣷⣿⣤⣤⡽⣿⣿⣿⣧⡴⣯⡾⠋⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠸⠖⠋⠉⠉⣩⡧⡽⣿⢷⣇⣀⣿⣿⡾⠿⣿⣿⠠⢶⣾⣿⣿⣷⣾⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⢀⣤⣿⣏⣍⣥⣽⣏⣽⣥⣬⣿⠷⠖⠉⢀⣠⡴⠿⠛⣹⣿⣿⡿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠸⣷⣝⣿⣷⣶⣾⣿⣿⣿⠉⢀⣠⠤⠚⠋⠀⠀⢶⠀⢹⣿⣿⠁⡠⣤⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⢹⣿⣙⣿⣿⣿⣿⣿⣿⣿⣟⡇⠀⠀⢀⣠⣄⣼⣤⢾⣻⠿⢛⣻⣿⣴⣤⣠⡄⠀⠀⠀⠀⠀⠀
⠀⣠⣤⣔⡛⢿⣟⣻⣉⣁⣬⡟⠻⢿⣿⣷⠶⢶⡏⠈⢻⢀⣄⣼⢁⣠⣽⡋⣙⣴⣧⣿⣧⡄⠀⠀⠀⠀⠀
⠈⠉⠉⠙⠛⠛⠛⠋⠀⠹⠛⠷⠴⢿⣾⣾⣾⣿⣿⢿⡿⣿⣿⣿⠿⠿⣿⡿⠿⠛⠋⠉⠁⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠁⠘⠿⠛⠙⠉⠠⠛⠻⠿⠛⠛⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀`
	return s

}
func (m model) renderConfig() string {
	//get list of templates
	keys := make([]string, 0, len(m.config.Templates))
	for key := range m.config.Templates {
		keys = append(keys, key)
	}
	//print fecth info
	var lines []string
	lines = append(lines, m.styles.promptStyle.Render("VERSION: ")+m.config.Version)
	lines = append(lines, m.styles.sepStyle.Render(m.styles.promptStyle.Render("URL: ")+m.config.Url))
	lines = append(lines, m.styles.promptStyle.Render("AUTHOR: ")+m.config.Author)
	lines = append(lines, m.styles.promptStyle.Render("EMAIL: ")+m.config.Email)
	lines = append(lines, m.styles.promptStyle.Render("WEBSITE: ")+m.config.Website)
	lines = append(lines, m.styles.promptStyle.Render("REMOTE: ")+m.config.Remote)
	lines = append(lines, m.styles.promptStyle.Render("EDITOR: ")+m.config.Editor)
	lines = append(lines, m.styles.promptStyle.Render("TEMPLATES: ")+fmt.Sprintf("%v", keys))
	lines = append(lines, m.styles.promptStyle.Render("PROJECTS_PATH: ")+m.config.Projects_path)
	lines = append(lines, m.styles.promptStyle.Render("TEMPLATES_PATH: ")+m.config.Templates_path)
	lines = append(lines, m.styles.promptStyle.Render("ARCHIVES_PATH: ")+m.config.Archives_path)
	//Add header to the lines
	header := m.styles.titleStyle.Render(`|\| /\ |_|`)
	return lipgloss.JoinVertical(lipgloss.Center, header, lipgloss.JoinVertical(lipgloss.Left, lines...))
}

func Home(config utils.Config) {

	m := newModel(config)
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
