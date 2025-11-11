package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// Update handles messages and updates the model
func (m *browseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Update active input based on mode
	switch m.mode {
	case browseMode:
		if m.filtering {
			m.filterInput, cmd = m.filterInput.Update(msg)
		}
		m.table, cmd = m.table.Update(msg)

	case addMode, editMode:
		// Check if title input is focused
		if m.titleInput.Focused() {
			m.titleInput, cmd = m.titleInput.Update(msg)
		} else {
			m.bodyInput, cmd = m.bodyInput.Update(msg)
		}
	default:
		panic("unhandled default case")
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		var cmd2 tea.Cmd
		_, cmd2 = m.handleKeyPress(msg)
		cmd = cmd2
	}

	return m, cmd
}

func (m *browseModel) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Global quit
	if msg.String() == "ctrl+c" {
		return m, tea.Quit
	}

	switch m.mode {
	case browseMode:
		return m.handleBrowseModeKeys(msg)
	case addMode, editMode:
		return m.handleEditModeKeys(msg)
	case deleteConfirmMode:
		return m.handleDeleteConfirmKeys(msg)
	}

	return m, nil
}

func (m *browseModel) handleBrowseModeKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.filtering {
		switch msg.String() {
		case "enter":
			m.filterQuery = m.filterInput.Value()
			m.filtering = false
			m.applyFilter()
			m.statusMsg = fmt.Sprintf("Filter applied: %q", m.filterQuery)
			return m, nil
		case "esc":
			m.filtering = false
			m.filterInput.SetValue("")
			m.statusMsg = "Filter cancelled"
			return m, nil
		}
		return m, nil
	}

	switch msg.String() {
	case "q":
		return m, tea.Quit

	case "a":
		m.mode = addMode
		m.titleInput.SetValue("")
		m.titleInput.Focus()
		m.bodyInput.SetValue("")
		m.statusMsg = ""
		return m, nil

	case "e", "enter":
		snippet := m.getSelectedSnippet()
		if snippet != nil {
			m.mode = editMode
			m.editingID = snippet.ID
			m.titleInput.SetValue(snippet.Title)
			m.titleInput.Focus()
			m.bodyInput.SetValue(snippet.Body)
			m.statusMsg = ""
		}
		return m, nil

	case "d":
		if m.getSelectedSnippet() != nil {
			m.mode = deleteConfirmMode
			m.statusMsg = ""
		}
		return m, nil

	case "/":
		m.filtering = true
		m.filterInput.Focus()
		m.statusMsg = ""
		return m, nil

	case "esc":
		if m.filterQuery != "" {
			m.filterQuery = ""
			m.applyFilter()
			m.statusMsg = "Filter cleared"
		}
		return m, nil

	case "s":
		m.cycleSortField()
		return m, nil

	case "?":
		m.showHelp = !m.showHelp
		return m, nil

	case "c":
		snippet := m.getSelectedSnippet()
		if snippet != nil {
			m.statusMsg = fmt.Sprintf("Copied snippet %d to clipboard", snippet.ID)
			// TODO: Implement actual clipboard copy
		}
		return m, nil
	}

	return m, nil
}

func (m *browseModel) handleEditModeKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.mode = browseMode
		m.titleInput.Blur()
		m.bodyInput.Blur()
		m.statusMsg = "Cancelled"
		return m, nil

	case "enter":
		if m.titleInput.Focused() {
			// Move to body input
			m.titleInput.Blur()
			m.bodyInput.Focus()
			return m, nil
		}
		// Otherwise save (handled by textarea's ctrl+enter or separate save key)
		return m, nil

	case "ctrl+s":
		return m.saveSnippet()

	case "tab":
		if m.titleInput.Focused() {
			m.titleInput.Blur()
			m.bodyInput.Focus()
		} else {
			m.bodyInput.Blur()
			m.titleInput.Focus()
		}
		return m, nil

	}

	return m, nil
}

func (m *browseModel) handleDeleteConfirmKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y":
		snippet := m.getSelectedSnippet()
		if snippet != nil {
			// TODO: Delete from database
			// Remove from items
			for i, s := range m.items {
				if s.ID == snippet.ID {
					m.items = append(m.items[:i], m.items[i+1:]...)
					break
				}
			}
			m.applyFilter()
			m.statusMsg = fmt.Sprintf("Deleted snippet %d", snippet.ID)
		}
		m.mode = browseMode
		return m, nil

	case "n", "esc":
		m.mode = browseMode
		m.statusMsg = "Delete cancelled"
		return m, nil
	}

	return m, nil
}

func (m *browseModel) saveSnippet() (tea.Model, tea.Cmd) {
	title := m.titleInput.Value()
	body := m.bodyInput.Value()

	if title == "" {
		m.statusMsg = "Title cannot be empty"
		return m, nil
	}

	if m.mode == addMode {
		// Find next ID
		nextID := 1
		for _, s := range m.items {
			if s.ID >= nextID {
				nextID = s.ID + 1
			}
		}

		newSnippet := Snippet{
			ID:    nextID,
			Title: title,
			Body:  body,
		}

		// TODO: Save to database
		m.items = append(m.items, newSnippet)
		m.statusMsg = fmt.Sprintf("Added snippet %d", nextID)

	} else if m.mode == editMode {
		// Update existing snippet
		for i := range m.items {
			if m.items[i].ID == m.editingID {
				m.items[i].Title = title
				m.items[i].Body = body
				break
			}
		}
		// TODO: Update database
		m.statusMsg = fmt.Sprintf("Updated snippet %d", m.editingID)
	}

	m.applyFilter()
	m.mode = browseMode
	m.titleInput.Blur()
	m.bodyInput.Blur()

	return m, nil
}
