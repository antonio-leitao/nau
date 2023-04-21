package new

import (
	"fmt"
	"math"

	"github.com/charmbracelet/lipgloss"
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
	s.Header = lipgloss.NewStyle().Margin(0, 0, 1, 0)
	title_text := "#ffffd7" //230
	if !IsSufficientContrast(title_text, base_color) {
		title_text = "235"
	}

	s.Title = lipgloss.NewStyle().
		Margin(2, 1, 1, 0).
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

func IsSufficientContrast(bgColor string, fgColor string) bool {
	// Parse the background and foreground colors into RGB values
	bgR, bgG, bgB := hexToRGB(bgColor)
	fgR, fgG, fgB := hexToRGB(fgColor)

	// Calculate the relative luminance of the background and foreground colors
	bgL := calcRelativeLuminance(bgR, bgG, bgB)
	fgL := calcRelativeLuminance(fgR, fgG, fgB)

	// Calculate the contrast ratio between the background and foreground colors
	contrast := calcContrastRatio(bgL, fgL)

	// Check if the contrast ratio meets the minimum requirement
	return contrast >= 4.5
}

// Helper function to convert a hex color string to RGB values
func hexToRGB(hex string) (r, g, b float64) {
	r, g, b = float64(hexToInt(hex[1:3])), float64(hexToInt(hex[3:5])), float64(hexToInt(hex[5:7]))
	return
}

// Helper function to convert a 2-character hex string to an integer
func hexToInt(hex string) int {
	var result int
	fmt.Sscanf(hex, "%x", &result)
	return result
}

// Helper function to calculate the relative luminance of an RGB color
func calcRelativeLuminance(r, g, b float64) float64 {
	rsRGB := r / 255.0
	gsRGB := g / 255.0
	bsRGB := b / 255.0

	if rsRGB <= 0.03928 {
		rsRGB = rsRGB / 12.92
	} else {
		rsRGB = math.Pow(((rsRGB + 0.055) / 1.055), 2.4)
	}

	if gsRGB <= 0.03928 {
		gsRGB = gsRGB / 12.92
	} else {
		gsRGB = math.Pow(((gsRGB + 0.055) / 1.055), 2.4)
	}

	if bsRGB <= 0.03928 {
		bsRGB = bsRGB / 12.92
	} else {
		bsRGB = math.Pow(((bsRGB + 0.055) / 1.055), 2.4)
	}

	return 0.2126*rsRGB + 0.7152*gsRGB + 0.0722*bsRGB
}

// Helper function to calculate the contrast ratio between two colors
func calcContrastRatio(l1, l2 float64) float64 {
	if l1 < l2 {
		l1, l2 = l2, l1
	}
	return (l1 + 0.05) / (l2 + 0.05)
}
