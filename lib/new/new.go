package new

import (
	"fmt"
	"os"

	structs "github.com/antonio-leitao/nau/lib/structs"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Styles struct {
	BorderColor lipgloss.Color
	InputField  lipgloss.Style
	InactiveField lipgloss.Style
}

func DefaultStyles() *Styles {
	s := new(Styles)
	s.BorderColor = lipgloss.Color("36")
	s.InputField = lipgloss.NewStyle().BorderForeground(s.BorderColor).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(80)
	s.InactiveField = lipgloss.NewStyle().BorderForeground(lipgloss.Color("12")).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(80)
	return s
}

type Main struct {
	styles    *Styles
	index     int
	questions []Question
	width     int
	height    int
	done      bool
}

type Question struct {
	question string
	answer   string
	input    Input
}

func newQuestion(q string) Question {
	return Question{question: q}
}

func newShortQuestion(q string) Question {
	question := newQuestion(q)
	model := NewShortAnswerField()
	question.input = model
	return question
}

func newLongQuestion(q string) Question {
	question := newQuestion(q)
	model := NewLongAnswerField()
	question.input = model
	return question
}

func New(questions []Question) *Main {
	styles := DefaultStyles()
	return &Main{styles: styles, questions: questions}
}

func (m Main) Init() tea.Cmd {
	return m.questions[m.index].input.Blink
}

func (m Main) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "ctrl+d":
			m.done = true
		case "tab":
			m.Next()
		}
	}
	for i := range m.questions{
		if i==m.index{
			m.questions[i].input.Focus()
			m.questions[i].answer = m.questions[i].input.Value()
			m.questions[i].input, cmd = m.questions[i].input.Update(msg)
		}else{
			m.questions[i].input.Blur()
		}
	}
	return m, cmd
}



func (m Main) View() string {
	if m.done {
		var output string
		for _, q := range m.questions {
			output += fmt.Sprintf("%s: %s\n", q.question, q.answer)
		}
		return output
	}
	if m.width == 0 {
		return "loading..."
	}

	blocks := []string{}
	for index,current := range m.questions{
		var style lipgloss.Style
		style = m.styles.InactiveField
		if index == m.index{
			style = m.styles.InputField
		}
		blocks=append(blocks,
			lipgloss.JoinVertical(
				lipgloss.Left,
				current.question,
				style.Render(current.input.View()),
			),
		)
	}
	// stack some left-aligned strings together in the center of the window
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Left,
			blocks...,
		),
	)
}

func (m *Main) Next() {
	if m.index < len(m.questions)-1 {
		m.index++
	} else {
		m.index = 0
	}
}

func NewPrompt(config structs.Config, query string) {
	// init styles; optional, just showing as a way to organize styles
	// start bubble tea and init first model
	questions := []Question{newShortQuestion("what is your name?"), newShortQuestion("what is your favourite editor?"), newLongQuestion("what's your favourite quote?")}
	main := New(questions)

	p := tea.NewProgram(*main, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		os.Exit(1)
	}
}