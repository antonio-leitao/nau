package root

import (
	"fmt"
	"os"
	"sort"
	"strings"

	archive "github.com/antonio-leitao/nau/cmd/archive"
	open "github.com/antonio-leitao/nau/cmd/open"
	lib "github.com/antonio-leitao/nau/lib"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
	"github.com/sahilm/fuzzy"
)

var (
	targetProject lib.Project
	targetAction  string
	confirmation  bool
)

// we gonna need this for filtering and stuff
type Projects []lib.Project

func (p Projects) String(i int) string {
	return p[i].Name + p[i].Lang
}
func (p Projects) Len() int {
	return len(p)
}

const (
	bullet   = "•"
	ellipsis = "…"
)

type Styles struct {
	Header       lipgloss.Style
	BlurredStyle lipgloss.Style
	MutedStyle   lipgloss.Style
	//confirmation
	PromptStyle     lipgloss.Style
	notFound        lipgloss.Style
	SelectedStyle   lipgloss.Style
	UnselectedStyle lipgloss.Style
	//statusbar
	Title     lipgloss.Style
	Count     lipgloss.Style
	StatusBar lipgloss.Style
	//deco
	DividerDot lipgloss.Style
	//colors
	verySubduedColor lipgloss.AdaptiveColor
	subduedColor     lipgloss.AdaptiveColor
}

func DefaultStyles(base_color string) (s Styles) {
	title_text := "#ffffd7" //230
	if !lib.IsSufficientContrast(title_text, base_color) {
		title_text = "235"
	}
	s.verySubduedColor = lipgloss.AdaptiveColor{Light: "#DDDADA", Dark: "#3C3C3C"}
	s.subduedColor = lipgloss.AdaptiveColor{Light: "#9B9B9B", Dark: "#5C5C5C"}
	s.Header = lipgloss.NewStyle().Margin(1, 0, 0, 0)
	s.BlurredStyle = lipgloss.NewStyle().Foreground(s.subduedColor)
	s.MutedStyle = lipgloss.NewStyle().Foreground(s.verySubduedColor)
	//other stuff
	s.notFound = lipgloss.NewStyle().
		Foreground(s.subduedColor).
		Margin(1, 0, 1, 4)
	//statusbar
	s.Title = lipgloss.NewStyle().
		Foreground(s.subduedColor)
	s.Count = lipgloss.NewStyle().
		Foreground(s.verySubduedColor)
	s.StatusBar = lipgloss.NewStyle().
		Padding(1, 0, 1, 1)
	//deco
	s.DividerDot = lipgloss.NewStyle().
		Foreground(s.verySubduedColor).
		SetString(" " + bullet + " ")
	s.PromptStyle = lipgloss.NewStyle().
		Margin(2, 0, 1, 0)
	s.SelectedStyle = lipgloss.NewStyle().
		Background(lipgloss.Color(base_color)).
		Foreground(lipgloss.Color(title_text)).
		Align(lipgloss.Center).
		Padding(0, 3).
		Margin(1, 1)
	s.UnselectedStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("235")).
		Foreground(lipgloss.Color("254")).
		Align(lipgloss.Center).
		Padding(0, 3).
		Margin(1, 1)
	return s
}

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type keyMap struct {
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
	//commands
	Open    key.Binding
	Archive key.Binding
	//filter
	Filter               key.Binding
	ClearFilter          key.Binding
	AcceptWhileFiltering key.Binding
	//confirmation
	Enter  key.Binding
	Toggle key.Binding
	//general
	Help   key.Binding
	Quit   key.Binding
	Cancel key.Binding
	// The quit-no-matter-what keybinding. This will be caught when filtering.
	ForceQuit key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Open, k.Archive, k.Filter, k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right}, // first column
		{k.Open, k.Archive},
		{k.Filter, k.Help, k.Quit},
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "right"),
	),
	Open: key.NewBinding(
		key.WithKeys("o", "enter"),
		key.WithHelp("o", "open"),
	),
	Archive: key.NewBinding(
		key.WithKeys("a", "delete"),
		key.WithHelp("a", "archive"),
	),
	Filter: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "filter"),
	),
	ClearFilter: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "clear filter"),
	),
	// Filtering.
	Cancel: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	),
	AcceptWhileFiltering: key.NewBinding(
		key.WithKeys("enter", "tab", "shift+tab", "ctrl+k", "up", "ctrl+j", "down"),
		key.WithHelp("enter", "apply filter"),
	),
	//confirmation keys
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "confirm"),
	),
	Toggle: key.NewBinding(
		key.WithKeys("left", "shift+tab", "tab", "right", "h", "l"),
		key.WithHelp("←/→/tab", "choose"),
	),
	//general
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "crtl+c"),
		key.WithHelp("q", "quit"),
	),
	ForceQuit: key.NewBinding(key.WithKeys("ctrl+c")),
}

type model struct {
	keys             keyMap
	help             help.Model
	inputStyle       lipgloss.Style
	lastKey          string
	quitting         bool
	state            string
	cursor           int
	initial_projects Projects
	projects         []lib.Project
	width            int
	columnWidth      int
	numCols          int
	selectedProject  lib.Project
	styles           Styles
	gap              int
	filter           textinput.Model
}

func newModel(base_color string, projects Projects) model {
	m := model{
		cursor:           0,
		initial_projects: projects,
		keys:             keys,
		state:            "browsing",
		help:             help.New(),
		width:            70,
		styles:           DefaultStyles(base_color),
		gap:              2,
		filter:           textinput.New(),
	}
	confirmation = false
	m.filter.Prompt = ""
	m.applyFilter()
	//update really has to be after
	m.updateGrid()
	return m
}
func (m *model) applyFilter() {
	m.cursor = 0
	//get query
	query := m.filter.Value()
	if query == "" {
		m.projects = m.initial_projects
	} else {
		//get matches
		results := fuzzy.FindFrom(query, m.initial_projects)
		var filteredProjects []lib.Project
		for _, r := range results {
			filteredProjects = append(filteredProjects, m.initial_projects[r.Index])
		}
		//update projects
		m.projects = filteredProjects
	}
	//toggle state if necessary
	if len(m.projects) < 1 {
		m.state = "stasis"
	}
}
func (m *model) resetFilter() {
	m.filter.SetValue("")
	m.applyFilter()
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width
	case tea.KeyMsg:
		if key.Matches(msg, m.keys.ForceQuit) {
			m.quitting = true
			return m, tea.Quit
		}
	}
	switch m.state {
	case "filtering":
		cmds = append(cmds, m.handleFiltering(msg))
	case "stasis":
		cmds = append(cmds, m.handleStasis(msg))
	case "confirm":
		cmds = append(cmds, m.handleConfirmation(msg))
	default:
		cmds = append(cmds, m.handleBrowsing(msg))
	}
	return m, tea.Batch(cmds...)
}
func (m *model) handleFiltering(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	// Handle keys
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch {
		case key.Matches(msg, m.keys.Cancel):
			m.resetFilter()
			m.applyFilter()
			m.state = "browsing"

		case key.Matches(msg, m.keys.AcceptWhileFiltering):
			m.filter.Blur()
			m.state = "browsing"
			m.applyFilter()
		}
	}
	// Update the filter text input component
	newFilterInputModel, inputCmd := m.filter.Update(msg)
	filterChanged := m.filter.Value() != newFilterInputModel.Value()
	if filterChanged {
		m.applyFilter()
	}
	m.filter = newFilterInputModel
	cmds = append(cmds, inputCmd)

	return tea.Batch(cmds...)
}
func (m *model) handleStasis(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			m.quitting = true
			return tea.Quit
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Filter):
			m.state = "filtering"
			m.filter.Focus()
		}
	}
	return tea.Batch(cmds...)
}

func (m *model) handleBrowsing(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {

		case key.Matches(msg, m.keys.Quit):
			m.quitting = true
			return tea.Quit
		case key.Matches(msg, m.keys.Up):
			// Check if the cursor is in the top row
			if m.cursor >= m.numCols {
				m.cursor -= m.numCols
			}
			m.selectedProject = m.projects[m.cursor]
			targetProject = m.selectedProject
		case key.Matches(msg, m.keys.Down):
			// Check if the cursor is in the last row
			if m.cursor+m.numCols < len(m.projects) {
				m.cursor += m.numCols
			}
			m.selectedProject = m.projects[m.cursor]
			targetProject = m.selectedProject
		case key.Matches(msg, m.keys.Left):
			// Check if the cursor is in the least column
			if m.cursor%m.numCols != 0 {
				m.cursor--
			}
			m.selectedProject = m.projects[m.cursor]
			targetProject = m.selectedProject
		case key.Matches(msg, m.keys.Right):
			// Check if the cursor is in the last column
			if (m.cursor+1)%m.numCols != 0 && m.cursor < len(m.projects)-1 {
				m.cursor++
			}
			m.selectedProject = m.projects[m.cursor]
			targetProject = m.selectedProject
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Open):
			targetAction = "open"
			m.state = "confirm"
		case key.Matches(msg, m.keys.Archive):
			targetAction = "archive"
			m.state = "confirm"
		case key.Matches(msg, m.keys.Filter):
			m.state = "filtering"
			m.filter.Focus()
		}
	}
	return tea.Batch(cmds...)
}
func (m *model) handleConfirmation(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		//next field
		case key.Matches(msg, m.keys.Toggle):
			confirmation = !confirmation
		//submission
		case key.Matches(msg, m.keys.Enter):
			if confirmation {
				m.quitting = true
				return tea.Quit
			} else {
				m.state = "browsing"
			}
		case key.Matches(msg, m.keys.Cancel):
			m.state = "browsing"
		case key.Matches(msg, m.keys.Quit):
			m.state = "browsing"
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		}
	}
	return tea.Batch(cmds...)
}

func (m model) View() string {
	if m.quitting {
		return ""
	}
	if m.state == "confirm" {
		return m.confirmationView()
	}
	//header
	var sections []string
	//status
	sections = append(sections, m.statusView())
	//Grid
	if m.state == "stasis" {
		sections = append(
			sections,
			m.styles.notFound.Render("No projects found :("),
		)
	} else {
		sections = append(sections, m.GridView())
	}
	//help
	helpView := m.help.View(m.keys)
	sections = append(sections, "\n"+helpView)
	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m model) confirmationView() string {

	var aff, neg string
	var sections []string
	query := fmt.Sprintf("%s %s?", strings.Title(targetAction), targetProject.Display_Name)
	sections = append(sections, m.styles.PromptStyle.Render(query))

	if confirmation {
		aff = m.styles.SelectedStyle.Render("Yes")
		neg = m.styles.UnselectedStyle.Render("No")
	} else {
		aff = m.styles.UnselectedStyle.Render("Yes")
		neg = m.styles.SelectedStyle.Render("No")
	}
	sections = append(sections,
		lipgloss.NewStyle().Margin(0, 0, 2, 0).Render(lipgloss.JoinHorizontal(lipgloss.Left, aff, neg)))

	return lipgloss.JoinVertical(lipgloss.Center, sections...)
}
func (m model) statusView() string {
	var status string
	if m.state != "filtering" {
		//title
		status += m.styles.Title.Render(`|\| /\ |_|`)
		//filter status
		if m.filter.Value() != "" {
			status += m.styles.DividerDot.String()
			f := strings.TrimSpace(m.filter.Value())
			f = truncate.StringWithTail(f, 10, "…")
			status += m.styles.Count.Render(fmt.Sprintf("“%s” ", f))
		}
		//number of items
		totalItems := len(m.projects)
		var itemName string
		if totalItems != 1 {
			itemName = "projects"
		} else {
			itemName = "projects"
		}
		status += m.styles.DividerDot.String()
		itemsDisplay := fmt.Sprintf("%d %s", totalItems, itemName)
		status += m.styles.Count.Render(itemsDisplay)
	} else {
		//if we are filtering show the input thing
		status += m.styles.Title.Render("Filter: ")
		status += m.filter.View()
	}
	return m.styles.StatusBar.Render(status)
}

func (m model) GridView() string {
	// Create a new slice to store the modified strings
	modifiedStrings := make([]string, len(m.projects))
	// Print the grid
	var rows []string
	for i := 0; i < len(m.projects); i += m.numCols {
		end := i + m.numCols
		if end > len(m.projects) {
			end = len(m.projects)
		}
		// Apply the function on each element before printing
		for j := i; j < end; j++ {
			if m.state == "filtering" {
				modifiedStrings[j] = m.renderDimmedProject(j)
			} else {
				modifiedStrings[j] = m.renderProject(j)
			}
		}
		row := lipgloss.JoinHorizontal(
			lipgloss.Left,
			modifiedStrings[i:end]...,
		)
		rows = append(rows, row)
	}
	return lipgloss.JoinVertical(
		lipgloss.Left,
		rows...,
	)
}
func (m model) renderProject(index int) string {
	project := m.projects[index]
	title_text := "#ffffd7" //230
	if !lib.IsSufficientContrast(title_text, project.Color) {
		title_text = "235"
	}
	if m.cursor == index {
		return lipgloss.NewStyle().
			Width(m.columnWidth-m.gap).
			Margin(0, m.gap, 0, 0).
			Background(lipgloss.Color(project.Color)).
			Foreground(lipgloss.Color(title_text)).
			Render(project.Display_Name)

	} else {
		return lipgloss.NewStyle().
			Width(m.columnWidth-m.gap).
			Margin(0, m.gap, 0, 0).
			Render(project.Display_Name)
	}
}
func (m model) renderDimmedProject(index int) string {
	project := m.projects[index]

	return lipgloss.NewStyle().
		Width(m.columnWidth-m.gap).
		Foreground(m.styles.subduedColor).
		Margin(0, m.gap, 0, 0).
		Render(project.Display_Name)
}

func (m model) getColumnWidth() int {
	columnWidth := 0
	for _, project := range m.projects {
		if len(project.Display_Name) > columnWidth {
			columnWidth = len(project.Display_Name)
		}
	}
	return columnWidth + m.gap
}
func (m *model) updateGrid() {
	// Calculate the column widths
	m.columnWidth = m.getColumnWidth()
	//number of columns
	m.numCols = m.width / m.columnWidth
}

func Execute(config lib.Config) {
	//read projects
	projects, _ := lib.GetProjects(config)
	// Extract the names from the projects
	sort.SliceStable(projects, func(i, j int) bool {
		return projects[i].Timestamp.After(projects[j].Timestamp)
	})
	//instantiate model
	model := newModel(
		config.Base_color,
		Projects(projects),
	)
	//run the cli
	if _, err := tea.NewProgram(model).Run(); err != nil {
		fmt.Printf("Could not start program :(\n%v\n", err)
		os.Exit(1)
	}
	if !confirmation {
		return
	}
	switch targetAction {
	case "open":
		open.Open(targetProject.Path, config.Editor)
	case "archive":
		archive.Archive(targetProject.Path, config.Archives_path)
	}
}
