package ui

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
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
