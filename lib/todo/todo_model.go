package todo

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type KeyMap struct {
	Submit key.Binding
	Quit   key.Binding
	// Help toggle keybindings.
	ShowFullHelp  key.Binding
	CloseFullHelp key.Binding
}

var DefaultKeyMap = KeyMap{
	Submit: key.NewBinding(
		key.WithKeys("ctrl+d"),
		key.WithHelp("ctrl+d", "submit"),
	),
	//confirmation keys
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

type Styles struct {
	HelpStyle lipgloss.Style
	AreaStyle lipgloss.Style
}

func DefaultStyle() Styles {
	var s Styles
	s.HelpStyle = lipgloss.NewStyle().Padding(1, 0, 0, 2)
	s.AreaStyle = lipgloss.NewStyle().Margin(1, 0, 0, 1)
	return s
}

type TodoModel struct {
	memo     textarea.Model
	aborted  bool
	quitting bool
	KeyMap   KeyMap
	Help     help.Model
	styles   Styles
}

func initialTodoModel() TodoModel {
	m := TodoModel{
		memo:     textarea.New(),
		aborted:  false,
		quitting: false,
		KeyMap:   DefaultKeyMap,
		Help:     help.New(),
		styles:   DefaultStyle(),
	}
	m.memo.Placeholder = "XXX: next todo"
	m.memo.ShowLineNumbers = false
	m.memo.SetWidth(50)
	m.memo.SetHeight(6)
	m.memo.CharLimit = 500
	//style of the text area
	m.memo.Focus()
	return m
}
func (m TodoModel) Init() tea.Cmd {
	return nil
}
func (m TodoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.Submit):
			//validate and submit
			m.quitting = true
			return m, m.Submit()
		case key.Matches(msg, m.KeyMap.Quit):
			m.aborted = true
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, m.KeyMap.ShowFullHelp):
			fallthrough
		case key.Matches(msg, m.KeyMap.CloseFullHelp):
			m.Help.ShowAll = !m.Help.ShowAll
		}
	}
	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m TodoModel) View() string {
	if m.quitting {
		return ""
	}

	var sections []string
	sections = append(sections, m.styles.AreaStyle.Render(m.memo.View()))
	sections = append(sections, m.helpView())
	return lipgloss.JoinVertical(lipgloss.Left, sections...)

}

func (m TodoModel) helpView() string {
	return m.styles.HelpStyle.Render(m.Help.View(m))
}

// FullHelp returns bindings to show the full help view. It's part of the
// help.KeyMap interface.
func (m TodoModel) FullHelp() [][]key.Binding {
	kb := [][]key.Binding{
		{m.KeyMap.Submit},
		{m.KeyMap.Quit},
		{m.KeyMap.CloseFullHelp},
	}
	return kb
}

func (m TodoModel) ShortHelp() []key.Binding {
	kb := []key.Binding{
		m.KeyMap.Submit,
		m.KeyMap.Quit,
		m.KeyMap.ShowFullHelp,
	}
	return kb
}
func (m *TodoModel) updateInputs(msg tea.Msg) tea.Cmd {
	//handles writing to the text area.
	var cmd tea.Cmd
	m.memo, cmd = m.memo.Update(msg)
	return cmd
}
func (m TodoModel) Submit() tea.Cmd {

	memo := m.memo.Value()
	// Remove line breaks from input string
	item := strings.ReplaceAll(memo, "\n", "")

	// Create path to the .nau/todos file
	homeDir, _ := os.UserHomeDir()
	todosDir := filepath.Join(homeDir, ".nau")
	_ = os.MkdirAll(todosDir, 0700)
	todosFile := filepath.Join(todosDir, "todos")

	// Open the file for writing
	f, _ := os.OpenFile(todosFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	defer f.Close()

	// Write the item to the file
	_, _ = fmt.Fprintln(f, item)
	return tea.Quit
}
