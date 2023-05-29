package root
import (
	"fmt"
	"os"
	"sort"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
    lib "github.com/antonio-leitao/nau/lib"
    open "github.com/antonio-leitao/nau/cmd/open"
    // archive "github.com/antonio-leitao/nau/cmd/archive"
	"github.com/charmbracelet/lipgloss"
)

var (
	targetProject lib.Project
	targetAction  string
)

// we gonna need this for filtering and stuff
type Projects []lib.Project

func (p Projects) String(i int) string {
	return p[i].Name
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
	PromptStyle lipgloss.Style
	//statusbar
	Title     lipgloss.Style
	Count     lipgloss.Style
	StatusBar lipgloss.Style
	//decorations
	DividerDot lipgloss.Style
}

func DefaultStyles() (s Styles) {
	verySubduedColor := lipgloss.AdaptiveColor{Light: "#DDDADA", Dark: "#3C3C3C"}
	subduedColor := lipgloss.AdaptiveColor{Light: "#9B9B9B", Dark: "#5C5C5C"}
	s.Header = lipgloss.NewStyle().Margin(1, 0, 0, 0)
	s.BlurredStyle = lipgloss.NewStyle().Foreground(subduedColor)
	s.MutedStyle = lipgloss.NewStyle().Foreground(verySubduedColor)
	//statusbar
	s.Title = lipgloss.NewStyle().
		Foreground(subduedColor)
	s.Count = lipgloss.NewStyle().
		Foreground(verySubduedColor)
	s.StatusBar = lipgloss.NewStyle().
		Padding(1, 0, 1, 1)
	//deco
	s.DividerDot = lipgloss.NewStyle().
		Foreground(verySubduedColor).
		SetString(" " + bullet + " ")
	return s
}

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type keyMap struct {
	Up          key.Binding
	Down        key.Binding
	Left        key.Binding
	Right       key.Binding
	Open        key.Binding
	Archive     key.Binding
	Filter      key.Binding
	ClearFilter key.Binding
	// Keybindings used when setting a filter.
	CancelWhileFiltering key.Binding
	AcceptWhileFiltering key.Binding
	//general
	Help key.Binding
	Quit key.Binding
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
	CancelWhileFiltering: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	),
	AcceptWhileFiltering: key.NewBinding(
		key.WithKeys("enter", "tab", "shift+tab", "ctrl+k", "up", "ctrl+j", "down"),
		key.WithHelp("enter", "apply filter"),
	),
	//general
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

type model struct {
	keys            keyMap
	help            help.Model
	inputStyle      lipgloss.Style
	lastKey         string
	quitting        bool
	state           string
	cursor          int
	projects        Projects
	width           int
	columnWidth     int
	numCols         int
	selectedProject lib.Project
	styles          Styles
}

func newModel(projects Projects) model {
	m := model{
		cursor:   0,
		projects: projects,
		keys:     keys,
		state:    "browsing",
		help:     help.New(),
		width:    70,
		styles:   DefaultStyles(),
	}
	m.updateGrid()
	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// If we set a width on the help menu it can gracefully truncate
		// its view as needed.
		m.help.Width = msg.Width
		//maybe no need for width, just recalc number of columns

	case tea.KeyMsg:
		switch {
		//browsing
		case key.Matches(msg, m.keys.Up) && m.state == "browsing":
			m.lastKey = "↑"
			m.HandleKeyPress("up")
		case key.Matches(msg, m.keys.Down) && m.state == "browsing":
			m.lastKey = "↓"
			m.HandleKeyPress("down")
		case key.Matches(msg, m.keys.Left) && m.state == "browsing":
			m.lastKey = "←"
			m.HandleKeyPress("left")
		case key.Matches(msg, m.keys.Right) && m.state == "browsing":
			m.lastKey = "→"
			m.HandleKeyPress("right")
		case key.Matches(msg, m.keys.Open) && m.state == "browsing":
			m.lastKey = "opening"
			m.quitting = true
			targetAction = "open"
			return m, tea.Quit
		case key.Matches(msg, m.keys.Archive) && m.state == "browsing":
			m.lastKey = "archiving"
			m.quitting = true
			targetAction = "archive"
			return m, tea.Quit
		case key.Matches(msg, m.keys.Filter) && m.state == "browsing":
			m.state = "filtering"
		//filtering
		case key.Matches(msg, m.keys.CancelWhileFiltering) && m.state == "filtering":
			m.state = "browsing"
		case key.Matches(msg, m.keys.AcceptWhileFiltering) && m.state == "filtering":
			m.state = "browsing"
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Quit):
			m.quitting = true
			return m, tea.Quit
		}
	}

	return m, nil
}
func (m *model) HandleKeyPress(key string) {
	// Check the key and perform the corresponding action
	switch key {
	case "up":
		// Check if the cursor is in the top row
		if m.cursor >= m.numCols {
			m.cursor -= m.numCols
		}
	case "down":
		// Check if the cursor is in the last row
		if m.cursor+m.numCols < m.projects.Len() {
			m.cursor += m.numCols
		}
	case "left":
		// Check if the cursor is in the least column
		if m.cursor%m.numCols != 0 {
			m.cursor--
		}
	case "right":
		// Check if the cursor is in the last column
		if (m.cursor+1)%m.numCols != 0 && m.cursor < m.projects.Len()-1 {
			m.cursor++
		}
	}
	m.selectedProject = m.projects[m.cursor]
	targetProject = m.selectedProject
}
func (m model) View() string {
	if m.quitting {
		return ""
	}
	//header
	var sections []string
	//status
	sections = append(sections, m.statusView())
	//Grid
	sections = append(sections, m.GridView())
	//help
	helpView := m.help.View(m.keys)
	sections = append(sections, "\n"+helpView)
	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m model) statusView() string {
	var status string
	totalItems := m.projects.Len()

	var itemName string
	if totalItems != 1 {
		itemName = "projects"
	} else {
		itemName = "projects"
	}
	status += m.styles.Title.Render(`|\| /\ |_|`)
	status += m.styles.DividerDot.String()
	itemsDisplay := fmt.Sprintf("%d %s", totalItems, itemName)
	status += m.styles.Count.Render(itemsDisplay)
	return m.styles.StatusBar.Render(status)
}

func (m model) GridView() string {
	// Maybe this doesnt need to be called everytime
	m.updateGrid()
	// Create a new slice to store the modified strings
	modifiedStrings := make([]string, len(m.projects))
	// Print the grid
	var rows []string
	for i := 0; i < m.projects.Len(); i += m.numCols {
		end := i + m.numCols
		if end > m.projects.Len() {
			end = m.projects.Len()
		}
		// Apply the function on each element before printing
		for j := i; j < end; j++ {
			modifiedStrings[j] = m.renderProject(j)
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
			Width(m.columnWidth).
			Background(lipgloss.Color(project.Color)).
			Foreground(lipgloss.Color(title_text)).
			Render(project.Display_Name)

	} else {
		return lipgloss.NewStyle().
			Width(m.columnWidth).
			Render(project.Display_Name)
	}
}

func (m model) getColumnWidth() int {
	columnWidth := 0
	for _, project := range m.projects {
		if len(project.Display_Name) > columnWidth {
			columnWidth = len(project.Display_Name)
		}
	}
	return columnWidth + 1
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
		Projects(projects),
	)
	//run the cli
	if _, err := tea.NewProgram(model).Run(); err != nil {
		fmt.Printf("Could not start program :(\n%v\n", err)
		os.Exit(1)
	}
	switch targetAction {
	case "open":
        open.Open(targetProject.Path,config.Editor)
	}
}
