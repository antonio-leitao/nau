package new

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	HelpStyle       lipgloss.Style
	FocusedStyle lipgloss.Style
	BlurredStyle lipgloss.Style
	NoStyle lipgloss.Style

}

func DefaultStyles()(s Styles){
	//verySubduedColor := lipgloss.AdaptiveColor{Light: "#DDDADA", Dark: "#3C3C3C"}
	//subduedColor := lipgloss.AdaptiveColor{Light: "#9B9B9B", Dark: "#5C5C5C"}
	s.HelpStyle = lipgloss.NewStyle().Padding(1, 0, 0, 2)
	s.FocusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	s.BlurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	s.NoStyle             = lipgloss.NewStyle()
	return s
}