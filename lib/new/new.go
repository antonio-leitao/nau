package new

import (
	"fmt"

	structs "github.com/antonio-leitao/nau/lib/structs"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type KeyMap struct {
    Next key.Binding
    Prev key.Binding
	Submit key.Binding
	Quit key.Binding
	// Help toggle keybindings.
	ShowFullHelp  key.Binding
	CloseFullHelp key.Binding
	
}

var DefaultKeyMap = KeyMap{
    Next: key.NewBinding(
        key.WithKeys("tab"),        // actual keybindings
        key.WithHelp("tab", "next field"), // corresponding help text
    ),
    Prev: key.NewBinding(
        key.WithKeys("shift+tab"),
        key.WithHelp("shift+tab", "prev field"),
    ),
	Submit: key.NewBinding(
        key.WithKeys("ctrl+d"),
        key.WithHelp("ctrl+d", "submit"),
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
	showHelp   bool
	KeyMap 	   KeyMap
	Styles     Styles
	Help help.Model
	index int
	inputs     []textinput.Model
	summary textarea.Model
}

func newModel()Model{
	m:=Model{
		showHelp: true,
		KeyMap: DefaultKeyMap,
		Styles : DefaultStyles(),
		Help:help.New(),
		index: 0,
		inputs: make([]textinput.Model, 3),
		summary: textarea.New(),
	}

	//style of the inputs
	var t textinput.Model
	for i := range m.inputs{
		t = textinput.New()
		t.CursorStyle = m.Styles.FocusedStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Project Name"
			t.Focus()
			t.PromptStyle = m.Styles.FocusedStyle
			t.TextStyle = m.Styles.FocusedStyle
		case 1:
			t.Placeholder = "XXX"
			t.CharLimit = 3
		case 2:
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}

		m.inputs[i] = t
	}
	//style of the text area
	m.summary.Placeholder="Describe you project"
	m.summary.ShowLineNumbers = false
	m.summary.SetWidth(32)
	m.summary.SetHeight(5)
	m.summary.CharLimit = 120
	return m
}

func (m Model) Init()tea.Cmd{
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
    case tea.KeyMsg:
        switch {
		
		//next field 
		case key.Matches(msg,  m.KeyMap.Next):
			m.index++
			m,cmd := m.forceBounds()
			return m, cmd

		//previous field
		case key.Matches(msg,  m.KeyMap.Prev):
			m.index--
			m,cmd := m.forceBounds()
			return m, cmd

        case key.Matches(msg,  m.KeyMap.Quit):
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

func (m Model) View()string{
	var sections []string

	sections=append(sections,"Hello world\n\n")

	for i := range m.inputs {
		sections=append(sections,m.inputs[i].View())
	}
	sections=append(sections,m.summary.View())
	if m.showHelp {
		sections = append(sections, m.helpView())
	}
	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m Model) forceBounds()(tea.Model, tea.Cmd){
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
	if m.index == 3{
		cmd = m.summary.Focus()
	} else {
		m.summary.Blur()
	}
	cmds = append(cmds,cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	
	//update summary
	var cmd tea.Cmd
	m.summary, cmd = m.summary.Update(msg)
	cmds = append(cmds,cmd)
	return tea.Batch(cmds...)
}




func (m Model) helpView() string {
	return m.Styles.HelpStyle.Render(m.Help.View(m))
}

// FullHelp returns bindings to show the full help view. It's part of the
// help.KeyMap interface.
func (m Model) FullHelp() [][]key.Binding {
	kb := [][]key.Binding{{
		m.KeyMap.Next,
		m.KeyMap.Prev,
		m.KeyMap.Submit,
		m.KeyMap.Quit,
		m.KeyMap.CloseFullHelp,
	}}

	return kb
}

func (m Model) ShortHelp() []key.Binding {
	kb := []key.Binding{
		m.KeyMap.Submit,
		m.KeyMap.Quit,
		m.KeyMap.ShowFullHelp,
	}
	return kb
}


func New(config structs.Config, query string){
	p:=tea.NewProgram(newModel(),tea.WithAltScreen())
	if _, err :=p.Run();err !=nil{
		fmt.Println(err)
	}
}