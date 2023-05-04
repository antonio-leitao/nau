package todo

import (
	"github.com/charmbracelet/lipgloss"
)

type MemoStyle struct {
	// The Normal state.
	NormalTitle lipgloss.Style
	NormalDesc  lipgloss.Style

	// The selected item state.
	SelectedTitle lipgloss.Style
	SelectedDesc  lipgloss.Style

	// The dimmed state, for when the filter input is initially activated.
	DimmedTitle lipgloss.Style
	DimmedDesc  lipgloss.Style

	// Characters matching the current filter, if any.
	FilterMatch lipgloss.Style
}

// NewDefaultItemStyles returns style definitions for a default item. See
// DefaultItemView for when these come into play.
func NewDefaultItemStyles() MemoStyle {
	var s MemoStyle
	s.NormalTitle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"}).
		Padding(0, 0, 0, 2)

	s.NormalDesc = s.NormalTitle.Copy().
		Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"})

	s.SelectedTitle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"}).
		Foreground(lipgloss.AdaptiveColor{Light: "#EE6FF8", Dark: "#EE6FF8"}).
		Padding(0, 0, 0, 1)

	s.SelectedDesc = s.SelectedTitle.Copy().
		Foreground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"})

	s.DimmedTitle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"}).
		Padding(0, 0, 0, 2)

	s.DimmedDesc = s.DimmedTitle.Copy().
		Foreground(lipgloss.AdaptiveColor{Light: "#C2B8C2", Dark: "#4D4D4D"})

	s.FilterMatch = lipgloss.NewStyle().Underline(true)

	return s
}

type Memo struct {
	Title       string
	Description string
	Style       MemoStyle
}
func (m Memo) RenderSelected(width int) string {
    var sections []string
    if len(m.Title)>0{
        sections = append(sections,m.Style.SelectedTitle.Render(m.Title))
    }
    sections = append(sections, m.Style.SelectedDesc.Render(m.Description))
    return lipgloss.JoinVertical(lipgloss.Left,sections...)
}
func (m Memo) Render(width int)string{
    var sections []string
    if len(m.Title)>0{
        sections = append(sections,m.Style.NormalTitle.Render(m.Title))
    }
    sections = append(sections, m.Style.NormalDesc.Render(m.Description))
    return lipgloss.JoinVertical(lipgloss.Left,sections...)
}
