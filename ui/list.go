package ui

import (
	"database/sql"
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
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
	list list.Model
	keys *listKeyMap
	db   *sql.DB
}

func newListModel(snippets []Snippet, db *sql.DB) *listModel {
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
	l.KeyMap.CursorUp.SetKeys("up", "k")
	l.KeyMap.CursorDown.SetKeys("down", "j")
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{keys.quit, keys.edit, keys.delete}
	}

	return &listModel{list: l, keys: &keys, db: db}

}

func (m *listModel) Init() tea.Cmd {
	return nil
}

func (m *listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.delete):
			i := m.list.Index()
			if i >= 0 && i < len(m.list.Items()) {
				it := m.list.Items()[i].(snippetItem)
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
	return appStyle.Render(m.list.View())
}
