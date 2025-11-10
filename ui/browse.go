package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type Snippet struct {
	ID    int
	Title string
	Body  string
}

type browseModel struct {
	cursor     int
	items      []Snippet
	openNote   bool
	deleteMode bool
	searchMode bool
	selected   Snippet
}

func newBrowseModel(snippets []Snippet) browseModel {
	return browseModel{items: snippets}
}

func (m browseModel) SelectedInput() Snippet {
	if len(m.items) == 0 {
		return Snippet{}
	}

	return m.items[m.cursor]

}

func (m browseModel) Init() tea.Cmd {
	return nil
}

func (m browseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "j", "down":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}

		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "enter":
			m.openNote = true
			m.selected = m.SelectedInput()

		case "d":
			m.deleteMode = true
			m.selected = m.SelectedInput()

		case "/":
			m.searchMode = true
		}

	}
	return m, nil
}

func (m browseModel) View() string {
	if len(m.items) == 0 {
		return "No snippets yet. \n Press 'q' to quit.\n"
	}

	s := "📒 Your snippets:\n\n"
	for i, snip := range m.items {
		cursor := " "
		if i == m.cursor {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, snip.Title)
	}

	s += "\n↑/↓ move  •  Enter open  •  d delete  •  / search  •  q quit\n"
	return s

}
