package new

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	HelpStyle       lipgloss.Style
}

func DefaultStyles()(s Styles){
	//verySubduedColor := lipgloss.AdaptiveColor{Light: "#DDDADA", Dark: "#3C3C3C"}
	//subduedColor := lipgloss.AdaptiveColor{Light: "#9B9B9B", Dark: "#5C5C5C"}
	s.HelpStyle = lipgloss.NewStyle().Padding(1, 0, 0, 2)
	return s
}