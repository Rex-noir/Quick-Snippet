package ui

import (
	"database/sql"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// NewBrowseModel creates a new browse model with the given database and snippets
func NewBrowseModel(db *sql.DB, snippets []Snippet) tea.Model {
	keys := newBrowseKeyMap()

	// Initialize filter input
	filterInput := textinput.New()
	filterInput.Placeholder = "Type to filter..."
	filterInput.CharLimit = 50

	// Initialize title input for add/edit
	titleInput := textinput.New()
	titleInput.Placeholder = "Enter title..."
	titleInput.CharLimit = 100

	// Initialize body textarea for add/edit
	bodyInput := textarea.New()
	bodyInput.Placeholder = "Enter body..."
	bodyInput.SetHeight(10)

	// Initialize table
	t := initializeTable()

	m := browseModel{
		items:         snippets,
		filteredItems: snippets,
		filterInput:   filterInput,
		titleInput:    titleInput,
		bodyInput:     bodyInput,
		keys:          keys,
		table:         t,
		mode:          browseMode,
		filtering:     false,
		currentSort:   sortByID,
		sortAscending: true,
		showHelp:      false, // Start with help visible
		db:            db,
	}

	m.updateTable()
	return &m
}

// Init initializes the model
func (m *browseModel) Init() tea.Cmd {
	return nil
}

// getSelectedSnippet returns the currently selected snippet
func (m *browseModel) getSelectedSnippet() *Snippet {
	cursor := m.table.Cursor()
	if cursor < 0 || cursor >= len(m.filteredItems) {
		return nil
	}
	return &m.filteredItems[cursor]
}
