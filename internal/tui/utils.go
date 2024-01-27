package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	color            = lipgloss.AdaptiveColor{Light: "#111222", Dark: "#FAFAFA"}
	defaultTextStyle = lipgloss.NewStyle().Foreground(color)
)

func boldString(s string) string {
	style := defaultTextStyle.Copy().Bold(true)
	return style.Render(s)
}

func errorString(s string) string {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#CC0000")).Bold(true)
	return style.Render(s)
}

func successString(s string) string {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#198754")).Bold(true)
	return style.Render(s)
}

func printError(err error) string {
	return errorString("Error: ") + err.Error() + "\n\n" + "Press CTRL+C or q to quit"
}

func printSuccess(msg string) string {
	return successString("Success: ") + msg + "\n\n" + "Press CTRL+C or q to quit"
}
