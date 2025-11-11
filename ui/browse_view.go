package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			Padding(0, 0, 1, 0)

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Padding(1, 0, 0, 0)
)

// View renders the browse model
func (m *browseModel) View() string {
	var b strings.Builder

	// Title
	b.WriteString(titleStyle.Render("📝 Snippet Manager"))
	b.WriteString("\n\n")

	// Mode-specific view
	switch m.mode {
	case browseMode:
		b.WriteString(m.viewBrowseMode())
	case addMode:
		b.WriteString(m.viewAddMode())
	case editMode:
		b.WriteString(m.viewEditMode())
	case deleteConfirmMode:
		b.WriteString(m.viewDeleteConfirmMode())
	}

	// Status message
	if m.statusMsg != "" {
		b.WriteString("\n")
		b.WriteString(statusStyle.Render(m.statusMsg))
	}

	// Always show command bar at the bottom (context-sensitive)
	b.WriteString("\n\n")
	b.WriteString(m.renderModeSpecificHelp())

	// Show full help if toggled
	if m.showHelp {
		b.WriteString(m.renderFullHelp())
	}

	return b.String()
}

func (m *browseModel) viewBrowseMode() string {
	var b strings.Builder

	// Filter input
	if m.filtering {
		b.WriteString("Filter: ")
		b.WriteString(m.filterInput.View())
		b.WriteString("\n\n")
	} else if m.filterQuery != "" {
		b.WriteString(fmt.Sprintf("Active filter: %q (press / to change, esc to clear)\n\n",
			m.filterQuery))
	}

	// Sort indicator
	sortStr := "Sort: "
	switch m.currentSort {
	case sortByID:
		sortStr += "ID"
	case sortByTitle:
		sortStr += "Title"
	}
	if m.sortAscending {
		sortStr += " ↑"
	} else {
		sortStr += " ↓"
	}
	b.WriteString(sortStr)
	b.WriteString(fmt.Sprintf(" | Items: %d/%d\n\n",
		len(m.filteredItems), len(m.items)))

	// Table
	b.WriteString(m.table.View())

	return b.String()
}

func (m *browseModel) viewAddMode() string {
	var b strings.Builder

	b.WriteString("Add New Snippet\n\n")
	b.WriteString("Title:\n")
	b.WriteString(m.titleInput.View())
	b.WriteString("\n\nBody:\n")
	b.WriteString(m.bodyInput.View())

	return b.String()
}

func (m *browseModel) viewEditMode() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("Edit Snippet (ID: %d)\n\n", m.editingID))
	b.WriteString("Title:\n")
	b.WriteString(m.titleInput.View())
	b.WriteString("\n\nBody:\n")
	b.WriteString(m.bodyInput.View())

	return b.String()
}

func (m *browseModel) viewDeleteConfirmMode() string {
	snippet := m.getSelectedSnippet()
	if snippet == nil {
		return "No snippet selected"
	}

	return fmt.Sprintf(
		"Are you sure you want to delete:\n\n"+
			"ID: %d\n"+
			"Title: %s\n\n"+
			"Press 'y' to confirm or 'n'/esc to cancel",
		snippet.ID,
		snippet.Title,
	)
}
