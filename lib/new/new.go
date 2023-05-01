package new

import (
	"fmt"
	"strings"

	utils "github.com/antonio-leitao/nau/lib/utils"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Submission struct {
	project_name string
	folder_name  string
	repo_name    string
	description  string
	git          bool
}

type KeyMap struct {
	Next   key.Binding
	Prev   key.Binding
	Up     key.Binding
	Down   key.Binding
	Submit key.Binding
	Quit   key.Binding
	Enter  key.Binding
	Toggle key.Binding
	// Help toggle keybindings.
	ShowFullHelp  key.Binding
	CloseFullHelp key.Binding
}

var DefaultKeyMap = KeyMap{
	Next: key.NewBinding(
		key.WithKeys("tab"),         // actual keybindings
		key.WithHelp("tab", "next"), // corresponding help text
	),
	Prev: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "prev"),
	),
	// Browsing.
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Submit: key.NewBinding(
		key.WithKeys("ctrl+d"),
		key.WithHelp("ctrl+d", "submit"),
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
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "esc"),
		key.WithHelp("ctrl+c/esc", "quit"),
	),
	// Toggle help.
	ShowFullHelp: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "more"),
	),
	CloseFullHelp: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "close help"),
	),
}

type Model struct {
	//existing projects
	existing_codes []string
	existing_names []string
	//basic info
	showHelp        bool
	KeyMap          KeyMap
	Styles          Styles
	Help            help.Model
	index           int
	inputs          []textinput.Model
	summary         textarea.Model
	errors          []string
	status          string
	template        string
	base_color      string
	templates       []string
	template_colors []string
	cursor          int
	spinner         spinner.Model
	//git confirmatio
	confirmation bool
	//really just to pass it along
	config *utils.Config
	//adaptivsize
	width  int
	height int
}

func newModel(
	base_color string,
	initial_status string,
	template string,
	templates []string,
	template_colors []string,
	existing_names []string,
	existing_codes []string,
	config *utils.Config,
) Model {
	m := Model{
		showHelp:        true,
		KeyMap:          DefaultKeyMap,
		Styles:          DefaultStyles(base_color),
		Help:            help.New(),
		index:           0,
		inputs:          make([]textinput.Model, 2),
		summary:         textarea.New(),
		errors:          []string{"", ""},
		status:          initial_status,
		template:        template,
		base_color:      base_color,
		templates:       templates,
		template_colors: template_colors,
		existing_codes:  existing_codes,
		existing_names:  existing_names,
		cursor:          0,
		spinner:         spinner.New(),
		confirmation:    true,
		config:          config,
	}

	//style of the inputs
	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.CursorStyle = m.Styles.FocusedStyle
		t.CharLimit = 20

		switch i {
		case 0:
			t.Placeholder = "Project Name"
			t.Focus()
			t.PromptStyle = m.Styles.FocusedStyle
			t.TextStyle = m.Styles.FocusedStyle
		case 1:
			t.Placeholder = "XXX"
			t.CharLimit = 3
		}

		m.inputs[i] = t
	}
	//static props of the text area
	m.summary.Placeholder = "Describe you project"
	m.summary.ShowLineNumbers = false
	m.summary.SetWidth(50)
	m.summary.SetHeight(6)
	m.summary.CharLimit = 500
	//style of the text area
	m.summary.FocusedStyle.CursorLine = m.Styles.NoStyle

	return m
}

func (m Model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	switch m.status {
	case "choose":
		return m.UpdateChoose(msg)
	case "confirm":
		return m.UpdateConfirm(msg)
	case "waiting":
		return m.UpdateWaiting(msg)
	default:
		return m.UpdateInfo(msg)
	}
}

func (m Model) UpdateChoose(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {

		//next field
		case key.Matches(msg, m.KeyMap.Up):
			m.cursor--
			if m.cursor > len(m.templates) {
				m.cursor = 0
			} else if m.cursor < 0 {
				m.cursor = len(m.templates)
			}
			return m, nil

		case key.Matches(msg, m.KeyMap.Down):
			m.cursor++
			if m.cursor > len(m.templates) {
				m.cursor = 0
			} else if m.cursor < 0 {
				m.cursor = len(m.templates)
			}
			return m, nil

		//submission
		case key.Matches(msg, m.KeyMap.Enter):
			m.template = m.templates[m.cursor]
			m.base_color = m.template_colors[m.cursor]
			//recompute styles based on selection
			m.Styles = DefaultStyles(m.base_color)
			if allStringsEmpty(m.errors) {
				m.status = "info"
			}

		case key.Matches(msg, m.KeyMap.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.KeyMap.ShowFullHelp):
			fallthrough
		case key.Matches(msg, m.KeyMap.CloseFullHelp):
			m.Help.ShowAll = !m.Help.ShowAll
		}
	}
	return m, nil
}

func (m Model) UpdateInfo(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {

		//next field
		case key.Matches(msg, m.KeyMap.Next):
			m.index++
			m, cmd := m.forceBounds()
			return m, cmd

		//previous field
		case key.Matches(msg, m.KeyMap.Prev):
			m.index--
			m, cmd := m.forceBounds()
			return m, cmd

		//submission
		case key.Matches(msg, m.KeyMap.Submit):
			if allStringsEmpty(m.errors) {
				m.status = "confirm"
			}
		case key.Matches(msg, m.KeyMap.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.KeyMap.ShowFullHelp):
			fallthrough
		case key.Matches(msg, m.KeyMap.CloseFullHelp):
			m.Help.ShowAll = !m.Help.ShowAll
		}
	}
	// Handle character input
	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m Model) UpdateConfirm(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {

		//next field
		case key.Matches(msg, m.KeyMap.Toggle):
			m.confirmation = !m.confirmation
			return m, nil
		//submission
		case key.Matches(msg, m.KeyMap.Enter):
			m.status = "waiting"
			return m, m.Submit()
		case key.Matches(msg, m.KeyMap.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.KeyMap.ShowFullHelp):
			fallthrough
		case key.Matches(msg, m.KeyMap.CloseFullHelp):
			m.Help.ShowAll = !m.Help.ShowAll
		}
	}
	return m, nil
}

func (m Model) UpdateWaiting(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.KeyMap.ShowFullHelp):
			fallthrough
		case key.Matches(msg, m.KeyMap.CloseFullHelp):
			m.Help.ShowAll = !m.Help.ShowAll
		}
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m Model) View() string {
	header := m.Styles.Header.Render(
		lipgloss.JoinVertical(
			lipgloss.Center,
			"New Project",
			m.Styles.ErrorStyle.Render("Start a new project from a template"),
		),
	)

	var output string
	switch m.status {
	case "choose":
		output = lipgloss.JoinVertical(lipgloss.Center, header, m.Styles.App.Render(m.ChooseView()))
	case "info":
		output = lipgloss.JoinVertical(lipgloss.Center, header, m.Styles.App.Render(m.InfoView()))
	case "confirm":
		output = lipgloss.JoinVertical(lipgloss.Center, header, m.Styles.App.Render(m.ConfirmView()))
	case "waiting":
		output = lipgloss.JoinVertical(lipgloss.Center, header, m.Styles.App.Render(m.WaitingView()))
	}
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		output,
	)
}

func (m Model) ChooseView() string {
	var sections []string
	var choice string
	for i, lang := range m.templates {
		title_text := "#ffffd7" //230
		if m.cursor == i {
			color := m.template_colors[i]
			if !IsSufficientContrast(title_text, color) {
				title_text = "235"
			}
			choice = m.Styles.SelectedTemplate.
				Background(lipgloss.Color(color)).
				Foreground(lipgloss.Color(title_text)).
				Render(lang)
		} else {
			choice = m.Styles.UnselectedTemplate.
				Render(lang)
		}
		sections = append(sections, choice)
	}

	if m.showHelp {
		sections = append(sections, m.helpView())
	}
	return lipgloss.JoinVertical(lipgloss.Center, sections...)
}

func (m Model) InfoView() string {
	var sections []string
	//make styles here
	sections = append(sections, m.Styles.Title.Render("Name and ID"))
	for i := range m.inputs {
		sections = append(
			sections,
			lipgloss.JoinHorizontal(
				lipgloss.Left,
				lipgloss.NewStyle().Width(23).Render(m.inputs[i].View()),
				m.errors[i],
			),
		)
	}
	sections = append(sections, m.Styles.Title.Render("Description"))
	sections = append(sections, m.summary.View())
	if m.showHelp {
		sections = append(sections, m.helpView())
	}
	return lipgloss.NewStyle().Padding(0, 1).Render(lipgloss.JoinVertical(lipgloss.Left, sections...))
}

func (m Model) ConfirmView() string {

	var aff, neg string
	var sections []string
	sections = append(sections, m.Styles.PromptStyle.Render("Start Git?"))

	if m.confirmation {
		aff = m.Styles.SelectedStyle.Render("Yes")
		neg = m.Styles.UnselectedStyle.Render("No")
	} else {
		aff = m.Styles.UnselectedStyle.Render("Yes")
		neg = m.Styles.SelectedStyle.Render("No")
	}
	sections = append(sections,
		lipgloss.NewStyle().Margin(0, 0, 2, 0).Render(lipgloss.JoinHorizontal(lipgloss.Left, aff, neg)))
	if m.showHelp {
		sections = append(sections, m.helpView())
	}

	return lipgloss.JoinVertical(lipgloss.Center, sections...)
}

func (m Model) WaitingView() string {
	var sections []string
	spinner := fmt.Sprintf("\n\n   %s Creating new project\n\n", m.spinner.View())
	sections = append(sections, spinner)
	if m.showHelp {
		sections = append(sections, m.helpView())
	}
	return lipgloss.JoinVertical(lipgloss.Center, sections...)
}

func contained(str string, slice []string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

func allStringsEmpty(strList []string) bool {
	for _, str := range strList {
		if len(str) != 0 {
			return false
		}
	}
	return true
}

func (m Model) Validate() {
	//if empty conditions
	m.errors[0] = ""
	m.errors[1] = ""
	//validate name and code
	name := utils.ToHyphenName(m.inputs[0].Value())
	if contained(name, m.existing_names) {
		m.errors[0] = m.Styles.ErrorStyle.Render("• Name already in use")
	}
	code := strings.ToUpper(m.inputs[1].Value())
	if contained(code, m.existing_codes) {
		m.errors[1] = m.Styles.ErrorStyle.Render("• Code already in use")
	}

	if len(m.inputs[0].Value()) == 0 {
		m.errors[0] = m.Styles.WarningStyle.Render("• Name cannot be empty")
	}
	if len(m.inputs[1].Value()) < 3 {
		m.errors[1] = m.Styles.WarningStyle.Render("• Code needs three letters")
	}
}

func (m Model) Submit() tea.Cmd {
	folder_name := utils.ToFolderName(m.inputs[0].Value())
	code := strings.ToUpper(m.inputs[1].Value())
	sub := Submission{
		project_name: utils.ToDunderName(m.inputs[0].Value()),
		folder_name:  code + "_" + folder_name,
		repo_name:    utils.ToHyphenName(m.inputs[0].Value()),
		description:  m.summary.Value(),
		git:          m.confirmation,
	}
	createNewProject(sub, m.config, m.template)
	return tea.Quit
}

func (m Model) forceBounds() (tea.Model, tea.Cmd) {
	if m.index > len(m.inputs) {
		m.index = 0
	} else if m.index < 0 {
		m.index = len(m.inputs)
	}

	cmds := make([]tea.Cmd, len(m.inputs))
	for i := 0; i <= len(m.inputs)-1; i++ {
		if i == m.index {
			// Set focused state
			cmds[i] = m.inputs[i].Focus()
			m.inputs[i].PromptStyle = m.Styles.FocusedStyle
			m.inputs[i].TextStyle = m.Styles.FocusedStyle
			continue
		}
		// Remove focused state
		m.inputs[i].Blur()
		m.inputs[i].PromptStyle = m.Styles.NoStyle
		m.inputs[i].TextStyle = m.Styles.NoStyle
	}

	//handle summary
	var cmd tea.Cmd
	if m.index == len(m.inputs) {
		cmd = m.summary.Focus()
	} else {
		m.summary.Blur()
	}
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) updateInputs(msg tea.Msg) tea.Cmd {
	m.Validate()
	cmds := make([]tea.Cmd, len(m.inputs))
	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	//update summary
	var cmd tea.Cmd
	m.summary, cmd = m.summary.Update(msg)
	cmds = append(cmds, cmd)
	return tea.Batch(cmds...)
}

func (m Model) helpView() string {
	return m.Styles.HelpStyle.Render(m.Help.View(m))
}

// FullHelp returns bindings to show the full help view. It's part of the
// help.KeyMap interface.
func (m Model) FullHelp() [][]key.Binding {
	switch m.status {
	case "choose":
		kb := [][]key.Binding{{
			m.KeyMap.Up,
			m.KeyMap.Down,
			m.KeyMap.Enter,
			m.KeyMap.Quit,
			m.KeyMap.CloseFullHelp,
		}}
		return kb

	case "confirm":
		kb := [][]key.Binding{{
			m.KeyMap.Toggle,
			m.KeyMap.Enter,
			m.KeyMap.Quit,
			m.KeyMap.CloseFullHelp,
		}}
		return kb
	case "waiting":
		kb := [][]key.Binding{{
			m.KeyMap.Quit,
			m.KeyMap.CloseFullHelp,
		}}
		return kb

	default:
		kb := [][]key.Binding{{
			m.KeyMap.Next,
			m.KeyMap.Prev,
			m.KeyMap.Submit,
			m.KeyMap.Quit,
			m.KeyMap.CloseFullHelp,
		}}
		return kb
	}
}

func (m Model) ShortHelp() []key.Binding {
	switch m.status {
	case "choose":
		kb := []key.Binding{
			m.KeyMap.Enter,
			m.KeyMap.Quit,
			m.KeyMap.ShowFullHelp,
		}
		return kb
	case "confirm":
		kb := []key.Binding{
			m.KeyMap.Enter,
			m.KeyMap.Quit,
			m.KeyMap.ShowFullHelp,
		}
		return kb
	case "waiting":
		kb := []key.Binding{
			m.KeyMap.Quit,
			m.KeyMap.ShowFullHelp,
		}
		return kb
	default:
		kb := []key.Binding{
			m.KeyMap.Submit,
			m.KeyMap.Quit,
			m.KeyMap.ShowFullHelp,
		}
		return kb
	}
}

func New(config utils.Config, query string) {
	//where do we start?
	initial_state, template, base_color := HandleArgs(config, query)
	//get all projects names
	projects, err := utils.GetProjects(config)
	if err != nil {
		fmt.Println(err)
	}
	//separate them
	var codes, repoNames []string
	for _, project := range projects {
		codes = append(codes, project.Code)
		repoNames = append(repoNames, project.Repo_name)
	}

	templates := []string{"Empty"}
	template_colors := []string{base_color}
	for lang, color := range config.Templates {
		templates = append(templates, lang)
		template_colors = append(template_colors, color)
	}

	//start application
	p := tea.NewProgram(
		newModel(
			base_color,
			initial_state,
			template,
			templates,
			template_colors,
			repoNames,
			codes,
			&config,
		),
		tea.WithAltScreen(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
	}
}
