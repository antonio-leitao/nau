package new

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	Title        lipgloss.Style
	HelpStyle    lipgloss.Style
	FocusedStyle lipgloss.Style
	WarningStyle lipgloss.Style
	ErrorStyle   lipgloss.Style
	BlurredStyle lipgloss.Style
	NoStyle      lipgloss.Style
	//confirmation
	PromptStyle     lipgloss.Style
	SelectedStyle   lipgloss.Style
	UnselectedStyle lipgloss.Style
}

func DefaultStyles(base_color string) (s Styles) {
	verySubduedColor := lipgloss.AdaptiveColor{Light: "#DDDADA", Dark: "#3C3C3C"}
	subduedColor := lipgloss.AdaptiveColor{Light: "#9B9B9B", Dark: "#5C5C5C"}
	s.Title = lipgloss.NewStyle().
		Margin(2, 1, 1, 0).
		Background(lipgloss.Color(base_color)).
		Foreground(lipgloss.Color("230"))
		//Background(lipgloss.Color("32")).
		//Foreground(lipgloss.Color("230"))


	s.HelpStyle = lipgloss.NewStyle().Padding(1, 0, 0, 2)
	s.FocusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("230"))
	s.BlurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("230"))
	s.NoStyle = lipgloss.NewStyle()
	s.WarningStyle = lipgloss.NewStyle().Foreground(verySubduedColor) //yellow
	s.ErrorStyle = lipgloss.NewStyle().Foreground(subduedColor)   //red
	s.PromptStyle = lipgloss.NewStyle().Margin(5, 0, 0, 0)
	s.SelectedStyle = lipgloss.NewStyle().Background(lipgloss.Color(base_color)).Foreground(lipgloss.Color("230")).Padding(0, 3).Margin(1, 1)
	s.UnselectedStyle = lipgloss.NewStyle().Background(lipgloss.Color("235")).Foreground(lipgloss.Color("254")).Padding(0, 3).Margin(1, 1)
	return s
}
