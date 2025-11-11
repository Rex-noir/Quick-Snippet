package ui

import (
	"database/sql"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func RunBrowse(db *sql.DB, snippets []Snippet) error {
	model := NewBrowseModel(db, snippets)
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		log.Println("Error running TUI", err)
		return err
	}
	return nil

}
