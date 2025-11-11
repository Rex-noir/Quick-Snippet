package ui

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type addInteractiveModel struct {
	title     textinput.Model
	body      textarea.Model
	focused   bool
	statusMsg string
	done      bool
	cancelled bool
	db        *sql.DB
}

func newAddInteractiveModel(initialTitle, initialBody *string, db *sql.DB) *addInteractiveModel {
	titleInput := textinput.New()
	titleInput.Placeholder = "Enter title..."
	titleInput.Focus()

	bodyInput := textarea.New()
	bodyInput.Placeholder = "Enter body..."
	bodyInput.SetHeight(10)

	return &addInteractiveModel{
		title:   titleInput,
		body:    bodyInput,
		focused: true,
		db:      db,
	}
}

func (m *addInteractiveModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *addInteractiveModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			if m.title.Focused() {
				m.title.Blur()
				m.body.Focus()
			} else {
				m.body.Blur()
				m.title.Focus()
			}
			return m, nil

		case "esc":
			m.cancelled = true
			return m, tea.Quit

		case "enter":
			// If in body and shift+enter not pressed, save
			if m.body.Focused() {
				m.done = true
				return m, tea.Quit
			}
			// Otherwise, move to body
			m.title.Blur()
			m.body.Focus()
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.body.SetWidth(msg.Width - 4)
	}

	// Update inputs
	if m.title.Focused() {
		m.title, cmd = m.title.Update(msg)
	} else {
		m.body, cmd = m.body.Update(msg)
	}

	return m, cmd
}

func (m *addInteractiveModel) View() string {
	var b strings.Builder

	b.WriteString("┌─ Add New Snippet ───────────────────────┐\n\n")
	b.WriteString(fmt.Sprintf("Title:\n%s\n\n", m.title.View()))
	b.WriteString(fmt.Sprintf("Body:\n%s\n\n", m.body.View()))
	b.WriteString("──────────────────────────────────────────\n")
	b.WriteString("Press [Tab] to switch, [Enter] to save, [Esc] to cancel\n")

	if m.statusMsg != "" {
		b.WriteString(fmt.Sprintf("\n%s\n", m.statusMsg))
	}

	return b.String()
}

func (m *addInteractiveModel) Value() (string, string, bool) {
	return m.title.Value(), m.body.Value(), m.done
}
