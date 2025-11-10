package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Snippet struct {
	ID    int
	Title string
	Body  string
}

type browseModel struct {
	table table.Model
	items []Snippet
}

func newBrowseModel(snippets []Snippet) tea.Model {
	columns := []table.Column{
		{Title: "ID", Width: 4},
		{Title: "Title", Width: 20},
	}

	rows := make([]table.Row, len(snippets))
	for i, s := range snippets {
		rows[i] = table.Row{fmt.Sprintf("%d", s.ID), s.Title}
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return browseModel{table: t, items: snippets}
}

func (m browseModel) Init() tea.Cmd { return nil }

func (m browseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			row := m.table.SelectedRow()
			if len(row) > 1 {
				return m, tea.Printf("Opening note: %s", row[1])
			}
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m browseModel) View() string {
	baseStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))
	return baseStyle.Render(m.table.View()) + "\nPress ↑/↓ to navigate, Enter to open, q to quit.\n"
}
