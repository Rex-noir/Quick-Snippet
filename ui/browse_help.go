package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m *browseModel) renderCommandBar() string {
	commands := []struct {
		key  string
		desc string
	}{
		{"↑/k", "up"},
		{"↓/j", "down"},
		{"a", "add"},
		{"e", "edit"},
		{"d", "delete"},
		{"/", "filter"},
		{"s", "sort"},
		{"c", "copy body"},
		{"?", "help"},
		{"q", "quit"},
	}

	var parts []string
	for _, cmd := range commands {
		parts = append(parts, commandStyle.Render(cmd.key)+descStyle.Render(" "+cmd.desc))
	}

	return helpStyle.Render(strings.Join(parts, " • "))

}

func (m *browseModel) renderModeSpecificHelp() string {
	switch m.mode {
	case addMode, editMode:
		return helpStyle.Render(
			commandStyle.Render("enter") + descStyle.Render(" save") + " • " +
				commandStyle.Render("esc") + descStyle.Render(" cancel"))
	case deleteConfirmMode:
		return helpStyle.Render(
			commandStyle.Render("y") + descStyle.Render(" confirm") + " • " +
				commandStyle.Render("n/esc") + descStyle.Render(" cancel"))
	case browseMode:
		if m.filtering {
			return helpStyle.Render(
				commandStyle.Render("enter") + descStyle.Render(" apply filter") + " • " +
					commandStyle.Render("esc") + descStyle.Render(" cancel"))
		}
		return m.renderCommandBar()
	default:
		return m.renderCommandBar()
	}
}

// renderFullHelp renders expanded help when showHelp is true
func (m *browseModel) renderFullHelp() string {
	if !m.showHelp {
		return ""
	}

	helpText := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1, 2).
		Render(`Navigation & Actions:
  ↑/k, ↓/j    Navigate up/down
  a           Add new snippet
  e           Edit selected snippet
  d           Delete selected snippet
  c           Copy snippet to clipboard

Filtering & Sorting:
  /           Start filtering
  s           Cycle sort (ID → Title → Date)
  esc         Clear filter/Cancel

Other:
  ?           Toggle this help
  q, Ctrl+C   Quit application`)

	return "\n" + helpText + "\n"
}
