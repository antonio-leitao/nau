package todo

import (
	"github.com/charmbracelet/lipgloss"
)

type MemoStyle struct {
	Content         lipgloss.Style
	SelectedContent lipgloss.Style
}

// NewDefaultItemStyles returns style definitions for a default item. See
// DefaultItemView for when these come into play.
func NewDefaultItemStyles() MemoStyle {
	var s MemoStyle
	verySubduedColor := lipgloss.AdaptiveColor{Light: "#DDDADA", Dark: "#3C3C3C"}
	// subduedColor := lipgloss.AdaptiveColor{Light: "#9B9B9B", Dark: "#5C5C5C"}
	s.Content = lipgloss.NewStyle().Foreground(verySubduedColor).Padding(0, 0, 0, 2).Margin(1,0,0,0)
		
	s.SelectedContent = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(lipgloss.Color("#e7e8e9")).
		Foreground(lipgloss.Color("#e7e8e9")).Padding(0, 0, 0, 1).Margin(1,0,0,0)
	return s
}

type Memo struct {
	Title       string
	Description string
	Style       MemoStyle
}

func (m Memo) RenderSelected(width int) string {
	return m.Style.Content.Width(width).Render(m.Title + m.Description)
}
func (m Memo) Render(width int) string {
	return m.Style.SelectedContent.Width(width).Render(m.Title + m.Description)
}
