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

func RunAddInteractive(db *sql.DB, initialTitle, initialBody *string) error {
	model := newAddInteractiveModel(initialTitle, initialBody, db)
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		log.Println("Error running TUI", err)
		return err
	}
	return nil
}

func RunListModel(db *sql.DB, snippets []Snippet) error {

	model := newListModel(snippets, db)
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		log.Println("Error running TUI", err)
		return err
	}
	return nil

}
