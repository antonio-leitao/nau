package new

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	Title lipgloss.Style
	HelpStyle       lipgloss.Style
	FocusedStyle lipgloss.Style
	WarningStyle lipgloss.Style
	ErrorStyle lipgloss.Style
	BlurredStyle lipgloss.Style
	NoStyle lipgloss.Style

}

func DefaultStyles()(s Styles){
	//verySubduedColor := lipgloss.AdaptiveColor{Light: "#DDDADA", Dark: "#3C3C3C"}
	//subduedColor := lipgloss.AdaptiveColor{Light: "#9B9B9B", Dark: "#5C5C5C"}
	s.Title = lipgloss.NewStyle().
	    Margin(2, 1,1,0).
		Background(lipgloss.Color("62")).
		Foreground(lipgloss.Color("230"))
		
	s.HelpStyle = lipgloss.NewStyle().Padding(1, 0, 0, 2)
	s.FocusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("230"))
	s.BlurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("230"))
	s.NoStyle             = lipgloss.NewStyle()
	s.WarningStyle	= lipgloss.NewStyle().Foreground(lipgloss.Color("#F9E2AF")) //yellow
	s.ErrorStyle	= lipgloss.NewStyle().Background(lipgloss.Color("#d20f39"))  //red
	return s
}