package ui

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func RunBrowse(snippets []Snippet) error {
	model := newBrowseModel(snippets)
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		log.Println("Error running TUI", err)
		return err
	}
	return nil

}
