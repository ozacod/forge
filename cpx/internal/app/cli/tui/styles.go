package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	cyan    = lipgloss.Color("#00D4FF")
	green   = lipgloss.Color("#00FF00")
	red     = lipgloss.Color("#FF0000")
	gray    = lipgloss.Color("#6B7280")
	white   = lipgloss.Color("#FFFFFF")
	dimGray = lipgloss.Color("#4B5563")

	// Styles for animated Q&A
	cyanBold = lipgloss.NewStyle().
			Foreground(cyan).
			Bold(true)

	greenCheck = lipgloss.NewStyle().
			Foreground(green).
			Bold(true)

	greenStyle = lipgloss.NewStyle().
			Foreground(green)

	questionMark = lipgloss.NewStyle().
			Foreground(cyan).
			Bold(true)

	questionStyle = lipgloss.NewStyle().
			Foreground(white).
			Bold(true)

	dimStyle = lipgloss.NewStyle().
			Foreground(dimGray)

	errorStyle = lipgloss.NewStyle().
			Foreground(red)

	selectedStyle = lipgloss.NewStyle().
			Foreground(cyan).
			Bold(true)

	inputPromptStyle = lipgloss.NewStyle().
				Foreground(cyan)

	inputTextStyle = lipgloss.NewStyle().
			Foreground(cyan)

	cursorStyle = lipgloss.NewStyle().
			Foreground(cyan)

	spinnerStyle = lipgloss.NewStyle().
			Foreground(cyan)
)

// renderCursor returns a cursor for the given position
func (m Model) renderCursor(pos int) string {
	if m.cursor == pos {
		return selectedStyle.Render("❯")
	}
	return " "
}
