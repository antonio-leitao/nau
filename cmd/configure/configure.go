package configure

// A simple example demonstrating the use of multiple text input components
// from the Bubbles component library.

import (
	"fmt"
	"log"
	"os"
	"strings"

	lib "github.com/antonio-leitao/nau/lib"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Styles struct {
	titleStyle    lipgloss.Style
	focusedStyle  lipgloss.Style
	blurredStyle  lipgloss.Style
	noStyle       lipgloss.Style
	focusedButton lipgloss.Style
	blurredButton lipgloss.Style
	promptPlace   lipgloss.Style
	warningStyle  lipgloss.Style
	errorStyle    lipgloss.Style
}

func DefaultStyles(base_color string) Styles {
	verySubduedColor := lipgloss.AdaptiveColor{Light: "#DDDADA", Dark: "#3C3C3C"}
	subduedColor := lipgloss.AdaptiveColor{Light: "#9B9B9B", Dark: "#5C5C5C"}
	title_text := "#ffffd7" //230
	var s Styles
	s.titleStyle = lipgloss.NewStyle().
		Margin(3, 0, 2, 0).
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color(base_color)).
		Foreground(lipgloss.Color(title_text))

	s.promptPlace = lipgloss.NewStyle().Width(50)
	s.focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(base_color))
	s.blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	s.noStyle = lipgloss.NewStyle()
	s.focusedButton = lipgloss.NewStyle().Background(lipgloss.Color(base_color)).Foreground(lipgloss.Color(title_text)).Align(lipgloss.Center).Padding(0, 3).Margin(1, 1)
	s.blurredButton = lipgloss.NewStyle().Background(lipgloss.Color("235")).Foreground(lipgloss.Color("254")).Align(lipgloss.Center).Padding(0, 3).Margin(1, 1)
	s.warningStyle = lipgloss.NewStyle().Foreground(verySubduedColor)
	s.errorStyle = lipgloss.NewStyle().Foreground(subduedColor)
	return s
}

type keyMap struct {
	Next   key.Binding
	Prev   key.Binding
	Submit key.Binding
	Help   key.Binding
	Quit   key.Binding
}

var keys = keyMap{
	Next: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "next"),
	),
	Prev: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "prev"),
	),
	Submit: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "submit"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

// show keybindings
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Next, k.Prev, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Next, k.Prev},   // first column
		{k.Submit, k.Quit}, // second column
	}
}

type Field struct {
	name  string
	input textinput.Model
}
type model struct {
	focusIndex int
	inputs     []Field
	keys       keyMap
	help       help.Model
	styles     Styles
	errors     []string
	line_width int
	height     int
	width      int
}

func initialModel(base_color string) model {
	m := model{
		inputs:     make([]Field, len(lib.CustomizableFields)),
		keys:       keys,
		help:       help.New(),
		styles:     DefaultStyles(base_color),
		errors:     make([]string, len(lib.CustomizableFields)),
		line_width: 48,
	}
	//start all errors as empty
	for i := range m.errors {
		m.errors[i] = ""
	}
	//start all the fields
	var f Field
	for i, field := range lib.CustomizableFields {
		f.input = textinput.New()
		f.input.CharLimit = m.line_width
		f.name = field
		f.input.Prompt = " • "
		f.input.Placeholder = normalizeString(field)
		//focus on the first field right away
		if i == 0 {
			f.input.Focus()
			f.input.Prompt = " > "
			f.input.PromptStyle = m.styles.focusedStyle
			f.input.TextStyle = m.styles.focusedStyle
		}
		m.inputs[i] = f
	}
	return m
}
func normalizeString(s string) string {
	// Replace underscores with spaces
	s = strings.ReplaceAll(s, "_", " ")
	// Capitalize each word
	s = strings.Title(strings.ToLower(s))
	return s
}
func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Next):
			m.focusIndex++
			return m.handleChange()
		case key.Matches(msg, m.keys.Prev):
			m.focusIndex--
			return m.handleChange()
		case key.Matches(msg, m.keys.Submit):
			if m.focusIndex == len(m.inputs) {
				//Evaluate first and submit later
				m.Validate()
				//if there are not errors just go
				if allStringsEmpty(m.errors) {
					return m, m.Submit
				}
			}
		}
	}
	// Handle character input and blinking
	cmd := m.updateInputs(msg)
	return m, cmd
}
func (m model) handleChange() (tea.Model, tea.Cmd) {
	//inbounds check
	if m.focusIndex > len(m.inputs) {
		m.focusIndex = 0
	} else if m.focusIndex < 0 {
		m.focusIndex = len(m.inputs)
	}

	//Handle change of focus
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := 0; i <= len(m.inputs)-1; i++ {
		if i == m.focusIndex {
			// Set focused state
			cmds[i] = m.inputs[i].input.Focus()
			m.inputs[i].input.Prompt = " > "
			m.inputs[i].input.PromptStyle = m.styles.focusedStyle
			m.inputs[i].input.TextStyle = m.styles.focusedStyle
			continue
		}
		// Remove focused state
		m.inputs[i].input.Blur()
		m.inputs[i].input.Prompt = " • "
		m.inputs[i].input.PromptStyle = m.styles.noStyle
		m.inputs[i].input.TextStyle = m.styles.noStyle
	}
	return m, tea.Batch(cmds...)
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i].input, cmds[i] = m.inputs[i].input.Update(msg)
	}
	return tea.Batch(cmds...)
}
func allStringsEmpty(strList []string) bool {
	for _, str := range strList {
		if len(str) != 0 {
			return false
		}
	}
	return true
}
func (m model) Validate() (tea.Model, tea.Cmd) {
	//clear all errors
	for i := range m.errors {
		m.errors[i] = ""
	}
	//validate entries
	for i, field := range lib.CustomizableFields {
		value := m.inputs[i].input.Value()
		if len(value) == 0 {
			m.errors[i] = m.styles.warningStyle.Render("• Field cannot be empty")
			continue
		}
		error_string := lib.ValidateValue(field, value)
		if error_string != "" {
			m.errors[i] = m.styles.errorStyle.Render(error_string)
		}
	}
	return m, nil
}

func (m model) Submit() tea.Msg {
	for _, pair := range m.inputs {
		lib.UpdateConfigField(pair.name, pair.input.Value())
	}
	return tea.Quit() //this will make program run ad infinitum. Change to tea.Quit()
}
func (m model) getLongestEntry() int {
	max_len := 0 //minwidth
	for _, input := range m.inputs {
		place_len := len(input.input.Placeholder)
		value_len := len(input.input.Value())
		if place_len >= max_len {
			max_len = place_len
		}
		if value_len >= max_len {
			max_len = value_len
		}
	}
	return max_len
}
func (m model) View() string {
	var b strings.Builder
	//Add header title
	title := m.styles.titleStyle.Render("Configure NAU")
	b.WriteString(title)
	b.WriteString("\n")
	//get max_width
	max_len := m.getLongestEntry() + 5
	//Add all prompts
	for i := range m.inputs {
		input_string := m.inputs[i].input.View()
		line_string := lipgloss.JoinHorizontal(
			lipgloss.Left,
			lipgloss.NewStyle().Width(max_len).Render(input_string),
			m.errors[i],
		)
		b.WriteString(m.styles.promptPlace.Render(line_string))
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}
	//submit button
	button := m.styles.blurredButton.Render("Submit")
	if m.focusIndex == len(m.inputs) {
		button = m.styles.focusedButton.Render("Submit")
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", button)
	//help
	b.WriteString("\n")
	b.WriteString(m.help.View(m.keys))

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		b.String(),
	)
}
func init_config(base_color string) {
	model := initialModel(base_color)
	if _, err := tea.NewProgram(model, tea.WithAltScreen()).Run(); err != nil {
		log.Printf("NAU ERROR could not start program: %s\n", err)
		os.Exit(1)
	}
}
func Execute(args []string) {
	if len(args) == 0 {
        config, err := lib.ReadConfig()
        if err != nil{
            log.Println("Could not load default config")
        }
		// No arguments provided
		init_config(config.Base_color)
	} else if len(args) == 1 {
        config, err := lib.LoadConfig()
        if err != nil{
            log.Println("Could not print config: ",err)
            fmt.Println(`Run:
nau config                 #to set all config values
nau config [field] [value] #to set individual ones
                `)
        }
		// Only field provided
		lib.OutputField(config, args[0])
	} else if len(args) > 1 {
		// Run program for field and value case
		err := lib.UpdateConfigField(args[0], strings.Join(args[1:], " "))
		if err != nil {
			log.Println("Error updating config file file:", err)
			os.Exit(1)
		}
	}
}
