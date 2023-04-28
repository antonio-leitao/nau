package configure

// A simple example demonstrating the use of multiple text input components
// from the Bubbles component library.

import (
	"fmt"
	"strings"

	"github.com/antonio-leitao/nau/lib/utils"
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
}

func DefaultStyles(base_color string) Styles {
	title_text := "#ffffd7" //230
	var s Styles
	s.titleStyle = lipgloss.NewStyle().
		Margin(3, 1, 1, 0).
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color(base_color)).
		Foreground(lipgloss.Color(title_text))

	s.promptPlace = lipgloss.NewStyle().Margin(0, 0, 0, 1)
	s.focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(base_color))
	s.blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	s.noStyle = lipgloss.NewStyle()
	s.focusedButton = lipgloss.NewStyle().Background(lipgloss.Color(base_color)).Foreground(lipgloss.Color(title_text)).Align(lipgloss.Center).Padding(0, 3).Margin(1, 1)
	s.blurredButton = lipgloss.NewStyle().Background(lipgloss.Color("235")).Foreground(lipgloss.Color("254")).Align(lipgloss.Center).Padding(0, 3).Margin(1, 1)
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
}

func initialModel(base_color string) model {
	m := model{
		inputs: make([]Field, len(customizableFields)),
		keys:   keys,
		help:   help.New(),
		styles: DefaultStyles(base_color),
	}

	var f Field
	for i, field := range customizableFields {
		f.input = textinput.New()
		f.input.CharLimit = 32
		f.input.Placeholder = field
		f.name = utils.ToDisplayName(field)
		//focus on the first field right away
		if i == 0 {
			f.input.Focus()
			f.input.PromptStyle = m.styles.focusedStyle
			f.input.TextStyle = m.styles.focusedStyle
		}
		m.inputs[i] = f
	}
	return m
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
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
				return m, tea.Quit
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
			m.inputs[i].input.PromptStyle = m.styles.focusedStyle
			m.inputs[i].input.TextStyle = m.styles.focusedStyle
			continue
		}
		// Remove focused state
		m.inputs[i].input.Blur()
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

func (m model) View() string {
	var b strings.Builder
	//Add header title
	title := m.styles.titleStyle.Render("Configure NAU")
	b.WriteString(title)
	b.WriteString("\n")
	//Add all prompts
	for i := range m.inputs {
		b.WriteString(m.styles.promptPlace.Render(m.inputs[i].input.View()))
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

	return b.String()
}
