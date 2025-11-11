package ui

import (
	"database/sql"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
)

type viewMode int

const (
	browseMode viewMode = iota
	addMode
	editMode
	deleteConfirmMode
)

type sortField int

const (
	sortByID sortField = iota
	sortByTitle
	sortByDate
)

type Snippet struct {
	ID    int
	Title string
	Body  string
}
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
