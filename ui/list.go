package ui

import (
	"QuickSnip/db"
	"QuickSnip/db/models"
	"database/sql"
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type snippetItem struct {
	Snippet
}

func (i snippetItem) Description() string { return i.Snippet.Body }
func (i snippetItem) FilterValue() string { return i.Snippet.Title }
func (i snippetItem) Title() string       { return i.Snippet.Title }

type listKeyMap struct {
	delete key.Binding
	quit   key.Binding
	edit   key.Binding
}

func newListKeyMap() listKeyMap {
	return listKeyMap{
		delete: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "delete"),
		),
		quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "quit"),
		),
		edit: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit"),
		),
	}
}

type listModel struct {
	list        list.Model
	keys        *listKeyMap
	db          *sql.DB
	editing     bool
	editForm    *editFormModel
	editingItem snippetItem
	editIndex   int
}

type editFormModel struct {
	titleInput textinput.Model
	bodyArea   textarea.Model
	focusIndex int
}

func newEditForm(snippet Snippet) *editFormModel {
	ti := textinput.New()
	ti.Placeholder = "Snippet Title"
	ti.Focus()
	ti.CharLimit = 200
	ti.Width = 50
	ti.SetValue(snippet.Title)

	ta := textarea.New()
	ta.Placeholder = "Snippet Body"
	ta.CharLimit = 5000
	ta.SetWidth(80)
	ta.SetHeight(10)
	ta.SetValue(snippet.Body)

	return &editFormModel{
		titleInput: ti,
		bodyArea:   ta,
		focusIndex: 0,
	}
}

func (m *editFormModel) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab":
			if m.focusIndex == 0 {
				m.focusIndex = 1
				m.titleInput.Blur()
				cmd = m.bodyArea.Focus()
				cmds = append(cmds, cmd)
			} else {
				m.focusIndex = 0
				m.bodyArea.Blur()
				m.titleInput.Focus()
			}
			return tea.Batch(cmds...)
		}
	}

	if m.focusIndex == 0 {
		m.titleInput, cmd = m.titleInput.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		m.bodyArea, cmd = m.bodyArea.Update(msg)
		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
}

func (m *editFormModel) View() string {
	var s string
	s += titleStyle.Render("Edit Snippet") + "\n\n"
	s += "Title:\n"
	s += m.titleInput.View() + "\n\n"
	s += "Body:\n"
	s += m.bodyArea.View() + "\n\n"
	s += lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("Tab: switch fields • Ctrl+S: save • Esc: cancel")
	return s
}

func newListModel(snippets []Snippet, database *sql.DB) *listModel {
	items := make([]list.Item, len(snippets))
	for i, s := range snippets {
		items[i] = snippetItem{Snippet: s}
	}
	keys := newListKeyMap()
	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = true
	delegate.Styles.SelectedTitle = lipgloss.NewStyle().Foreground(lipgloss.Color("#25A065")).Bold(true)

	l := list.New(items, delegate, 80, 20)
	l.Title = "Snippets"
	l.Styles.Title = titleStyle
	l.SetShowStatusBar(false)
	l.KeyMap.CursorUp.SetKeys("up", "k")
	l.KeyMap.CursorDown.SetKeys("down", "j")
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{keys.quit, keys.edit, keys.delete}
	}

	return &listModel{
		list:    l,
		keys:    &keys,
		db:      database,
		editing: false,
	}
}

func (m *listModel) Init() tea.Cmd {
	return nil
}

func (m *listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle edit mode
	if m.editing {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc":
				m.editing = false
				m.editForm = nil
				return m, nil
			case "ctrl+s":
				// Save the edited snippet
				updatedSnippet := models.Snippet{
					ID:    m.editingItem.ID,
					Title: m.editForm.titleInput.Value(),
					Body:  m.editForm.bodyArea.Value(),
				}

				_, err := db.SaveSnippet(m.db, updatedSnippet)
				if err != nil {
					statusMsg := fmt.Sprintf("Error updating snippet: %v", err)
					cmd := m.list.NewStatusMessage(statusMessageStyle(statusMsg))
					m.editing = false
					return m, cmd
				}

				// Update the item in the list
				m.list.SetItem(m.editIndex, snippetItem{Snippet: Snippet{ID: updatedSnippet.ID, Title: updatedSnippet.Title, Body: updatedSnippet.Body}})

				statusMsg := fmt.Sprintf("Updated snippet #%d", updatedSnippet.ID)
				cmd := m.list.NewStatusMessage(statusMessageStyle(statusMsg))
				m.editing = false
				m.editForm = nil
				return m, cmd
			}
		}

		cmd := m.editForm.Update(msg)
		return m, cmd
	}

	// Normal list mode
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.edit):
			i := m.list.Index()
			if i >= 0 && i < len(m.list.Items()) {
				item := m.list.Items()[i].(snippetItem)
				m.editingItem = item
				m.editIndex = i
				m.editForm = newEditForm(item.Snippet)
				m.editing = true
				return m, nil
			}

		case key.Matches(msg, m.keys.delete):
			i := m.list.Index()
			if i >= 0 && i < len(m.list.Items()) {
				it := m.list.Items()[i].(snippetItem)
				err := db.DeleteSnippet(m.db, it.ID)
				if err != nil {
					delMsg := fmt.Sprintf("Error deleting snippet #%d: %v", it.ID, err)
					cmd := m.list.NewStatusMessage(statusMessageStyle(delMsg))
					return m, cmd
				}
				delMsg := fmt.Sprintf("Deleted snippet #%d", it.ID)
				delCmd := m.list.NewStatusMessage(statusMessageStyle(delMsg))
				m.list.RemoveItem(i)
				return m, delCmd
			}
		}
	}

	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	return m, cmd
}

func (m *listModel) View() string {
	if m.editing && m.editForm != nil {
		return appStyle.Render(m.editForm.View())
	}
	return appStyle.Render(m.list.View())
}
