package new

import (
	"fmt"

	structs "github.com/antonio-leitao/nau/lib/structs"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
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

}

func newModel()Model{
	m:=Model{
		showHelp: true,
		KeyMap: DefaultKeyMap,
		Styles : DefaultStyles(),
		Help:help.New(),
	}
	return m
}

func (m Model) Init()tea.Cmd{
	
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
    case tea.KeyMsg:
        switch {
        case key.Matches(msg, DefaultKeyMap.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.KeyMap.ShowFullHelp):
			fallthrough
		case key.Matches(msg, m.KeyMap.CloseFullHelp):
			m.Help.ShowAll = !m.Help.ShowAll
        }
    }
	return m, nil
}

func (m Model) View()string{
	var sections []string

	sections=append(sections,"Hello world\n\n")

	if m.showHelp {
		sections = append(sections, m.helpView())
	}
	return lipgloss.JoinVertical(lipgloss.Left, sections...)
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