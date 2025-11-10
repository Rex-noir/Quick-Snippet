package ui

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
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
	nextID        int
	editingID     int
}

func NewBrowseModel(snippets []Snippet) tea.Model {
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

	// Find next ID
	nextID := 1
	for _, s := range snippets {
		if s.ID >= nextID {
			nextID = s.ID + 1
		}
	}

	m := browseModel{
		items:         snippets,
		filteredItems: snippets,
		filterInput:   filterInput,
		titleInput:    titleInput,
		bodyInput:     bodyInput,
		keys:          keys,
		mode:          browseMode,
		filtering:     false,
		currentSort:   sortByID,
		sortAscending: true,
		nextID:        nextID,
	}

	m.updateTable()
	return &m
}
