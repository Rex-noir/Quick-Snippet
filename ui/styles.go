package ui

import "github.com/charmbracelet/lipgloss"

var (
	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Padding(0, 1)

	commandStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true)

	descStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("246"))
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			Padding(0, 0, 1, 0)

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Padding(1, 0, 0, 0)
)
