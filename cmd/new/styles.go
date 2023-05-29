package new

import (
	"github.com/charmbracelet/lipgloss"
    lib "github.com/antonio-leitao/nau/lib"
)

type Styles struct {
	Header       lipgloss.Style
	App          lipgloss.Style
	Title        lipgloss.Style
	HelpStyle    lipgloss.Style
	FocusedStyle lipgloss.Style
	WarningStyle lipgloss.Style
	ErrorStyle   lipgloss.Style
	BlurredStyle lipgloss.Style
	NoStyle      lipgloss.Style

	//choose
	SelectedTemplate   lipgloss.Style
	UnselectedTemplate lipgloss.Style
	//confirmation
	PromptStyle     lipgloss.Style
	SelectedStyle   lipgloss.Style
	UnselectedStyle lipgloss.Style
}

func DefaultStyles(base_color string) (s Styles) {
	verySubduedColor := lipgloss.AdaptiveColor{Light: "#DDDADA", Dark: "#3C3C3C"}
	subduedColor := lipgloss.AdaptiveColor{Light: "#9B9B9B", Dark: "#5C5C5C"}

	s.App = lipgloss.NewStyle().Width(52).Align(lipgloss.Center)
	s.Header = lipgloss.NewStyle().Margin(1, 0, 0, 0)
	title_text := "#ffffd7" //230
	if !lib.IsSufficientContrast(title_text, base_color) {
		title_text = "235"
	}

	s.Title = lipgloss.NewStyle().
		Margin(1, 1, 1, 0).
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color(base_color)).
		Foreground(lipgloss.Color(title_text))
		//Background(lipgloss.Color("32")).
		//Foreground(lipgloss.Color("230"))

	s.HelpStyle = lipgloss.NewStyle().Padding(1, 0, 0, 2)
	s.FocusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("230"))
	s.BlurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("230"))
	s.NoStyle = lipgloss.NewStyle()
	s.WarningStyle = lipgloss.NewStyle().Foreground(verySubduedColor)
	s.ErrorStyle = lipgloss.NewStyle().Foreground(subduedColor)
	s.PromptStyle = lipgloss.NewStyle().Margin(5, 0, 0, 0)

	s.SelectedTemplate = lipgloss.NewStyle().Width(15).Align(lipgloss.Center).Padding(0, 3).Margin(1, 0, 0, 0)
	s.UnselectedTemplate = lipgloss.NewStyle().Width(15).Background(lipgloss.Color("235")).Align(lipgloss.Center).Padding(0, 3).Margin(1, 0, 0, 0)

	s.SelectedStyle = lipgloss.NewStyle().Background(lipgloss.Color(base_color)).Foreground(lipgloss.Color(title_text)).Align(lipgloss.Center).Padding(0, 3).Margin(1, 1)
	s.UnselectedStyle = lipgloss.NewStyle().Background(lipgloss.Color("235")).Foreground(lipgloss.Color("254")).Align(lipgloss.Center).Padding(0, 3).Margin(1, 1)
	return s
}

