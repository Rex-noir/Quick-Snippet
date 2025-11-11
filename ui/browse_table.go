package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

// updateTable refreshes the table with current filtered/sorted data
func (m *browseModel) updateTable() {
	rows := make([]table.Row, len(m.filteredItems))
	for i, item := range m.filteredItems {
		preview := strings.ReplaceAll(item.Body, "\n", " ")
		if len(preview) > 60 {
			preview = preview[:60] + "..."
		}
		rows[i] = table.Row{
			fmt.Sprintf("%d", item.ID),
			item.Title,
			preview,
		}
	}
	m.table.SetRows(rows)
}

// applyFilter filters items based on the current filter query
func (m *browseModel) applyFilter() {
	if m.filterQuery == "" {
		m.filteredItems = m.items
	} else {
		m.filteredItems = []Snippet{}
		query := strings.ToLower(m.filterQuery)
		for _, item := range m.items {
			if strings.Contains(strings.ToLower(item.Title), query) ||
				strings.Contains(strings.ToLower(item.Body), query) {
				m.filteredItems = append(m.filteredItems, item)
			}
		}
	}
	m.sortItems()
	m.updateTable()
}

// sortItems sorts the filtered items based on current sort settings
func (m *browseModel) sortItems() {
	sort.Slice(m.filteredItems, func(i, j int) bool {
		var less bool
		switch m.currentSort {
		case sortByID:
			less = m.filteredItems[i].ID < m.filteredItems[j].ID
		case sortByTitle:
			less = strings.ToLower(m.filteredItems[i].Title) <
				strings.ToLower(m.filteredItems[j].Title)
		default:
			less = m.filteredItems[i].ID < m.filteredItems[j].ID
		}

		if !m.sortAscending {
			less = !less
		}
		return less
	})
}

// cycleSortField cycles through sort fields
func (m *browseModel) cycleSortField() {
	switch m.currentSort {
	case sortByID:
		m.currentSort = sortByTitle
		m.statusMsg = "Sorting by Title"
	case sortByTitle:
		m.currentSort = sortByID
		m.sortAscending = !m.sortAscending
		if m.sortAscending {
			m.statusMsg = "Sorting by ID (ascending)"
		} else {
			m.statusMsg = "Sorting by ID (descending)"
		}
	default:
		panic("unhandled default case")
	}
	m.sortItems()
	m.updateTable()
}

// initializeTable creates and styles the table
func initializeTable() table.Model {
	columns := []table.Column{
		{Title: "ID", Width: 8},
		{Title: "Title", Width: 30},
		{Title: "Preview", Width: 60},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(20),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return t
}
