package ui

import (
	"database/sql"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type browseModel struct {
	table         table.Model
	items         []Snippet
	filteredItems []Snippet
	filterInput   textinput.Model
	titleInput    textinput.Model
	bodyInput     textarea.Model
	keys          browseKeyMap
	mode          viewMode
	filtering     bool
	filterQuery   string
	currentSort   sortField
	sortAscending bool
	statusMsg     string
	showHelp      bool
	width         int
	height        int
	editingID     int
	db            *sql.DB
}

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

	// Find next ID
	nextID := 1
	for _, snippet := range snippets {
		if snippet.ID >= nextID {
			nextID = snippet.ID + 1
		}
	}

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
		db:            db,
	}

	m.updateTable()
	return &m
}
