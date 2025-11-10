package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Snippet struct {
	ID    int
	Title string
	Body  string
}

type viewMode int

const (
	browseMode viewMode = iota
	previewMode
	editMode
	addMode
)

type sortField int

const (
	sortByID sortField = iota
	sortByTitle
)

func (m *browseModel) updateTable() {
	columns := []table.Column{
		{Title: "ID", Width: 6},
		{Title: "Title", Width: 40},
		{Title: "Preview", Width: 50},
	}

	rows := make([]table.Row, len(m.filteredItems))
	for i, s := range m.filteredItems {
		preview := s.Body
		if len(preview) > 47 {
			preview = preview[:47] + "..."
		}
		preview = strings.ReplaceAll(preview, "\n", " ")
		rows[i] = table.Row{fmt.Sprintf("%d", s.ID), s.Title, preview}
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
		Bold(true).
		Foreground(lipgloss.Color("229"))

	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)

	t.SetStyles(s)
	m.table = t
}

func (m *browseModel) applyFilter() {
	if m.filterQuery == "" {
		m.filteredItems = m.items
	} else {
		query := strings.ToLower(m.filterQuery)
		var filtered []Snippet
		for _, item := range m.items {
			if strings.Contains(strings.ToLower(item.Title), query) ||
				strings.Contains(strings.ToLower(item.Body), query) {
				filtered = append(filtered, item)
			}
		}
		m.filteredItems = filtered
	}
	m.updateTable()
}

func (m *browseModel) sortItems() {
	switch m.currentSort {
	case sortByID:
		sort.Slice(m.filteredItems, func(i, j int) bool {
			if m.sortAscending {
				return m.filteredItems[i].ID < m.filteredItems[j].ID
			}
			return m.filteredItems[i].ID > m.filteredItems[j].ID
		})
	case sortByTitle:
		sort.Slice(m.filteredItems, func(i, j int) bool {
			if m.sortAscending {
				return strings.ToLower(m.filteredItems[i].Title) < strings.ToLower(m.filteredItems[j].Title)
			}
			return strings.ToLower(m.filteredItems[i].Title) > strings.ToLower(m.filteredItems[j].Title)
		})
	}
	m.updateTable()
}

func (m *browseModel) getSelectedSnippet() *Snippet {
	row := m.table.SelectedRow()
	if len(row) == 0 {
		return nil
	}

	var id int
	_, err := fmt.Sscanf(row[0], "%d", &id)
	if err != nil {
		return nil
	}

	for i := range m.filteredItems {
		if m.filteredItems[i].ID == id {
			return &m.filteredItems[i]
		}
	}
	return nil
}

func (m *browseModel) deleteSelected() {
	snippet := m.getSelectedSnippet()
	if snippet == nil {
		return
	}

	// Remove from main items
	for i, item := range m.items {
		if item.ID == snippet.ID {
			m.items = append(m.items[:i], m.items[i+1:]...)
			break
		}
	}

	m.applyFilter()
	m.statusMsg = fmt.Sprintf("Deleted: %s", snippet.Title)
}

func (m *browseModel) Init() tea.Cmd {
	return nil
}

func (m *browseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.bodyInput.SetWidth(msg.Width - 4)

	case tea.KeyMsg:
		switch m.mode {
		case browseMode:
			if m.filtering {
				switch msg.String() {
				case "enter", "esc":
					m.filtering = false
					m.filterQuery = m.filterInput.Value()
					m.applyFilter()
					m.statusMsg = fmt.Sprintf("Filtered: %d results", len(m.filteredItems))
				default:
					m.filterInput, cmd = m.filterInput.Update(msg)
					return m, cmd
				}
			} else {
				switch {
				case key.Matches(msg, m.keys.quit):
					return m, tea.Quit

				case key.Matches(msg, m.keys.filter):
					m.filtering = true
					m.filterInput.SetValue("")
					m.filterInput.Focus()
					return m, nil

				case key.Matches(msg, m.keys.clearFilter):
					m.filterQuery = ""
					m.applyFilter()
					m.statusMsg = "Filter cleared"
					return m, nil

				case key.Matches(msg, m.keys.sort):
					if m.currentSort == sortByID {
						m.currentSort = sortByTitle
						m.sortAscending = true
					} else if m.currentSort == sortByTitle && m.sortAscending {
						m.sortAscending = false
					} else {
						m.currentSort = sortByID
						m.sortAscending = true
					}
					m.sortItems()
					sortDir := "↑"
					if !m.sortAscending {
						sortDir = "↓"
					}
					sortName := "ID"
					if m.currentSort == sortByTitle {
						sortName = "Title"
					}
					m.statusMsg = fmt.Sprintf("Sorted by %s %s", sortName, sortDir)
					return m, nil

				case key.Matches(msg, m.keys.help):
					m.showHelp = !m.showHelp
					return m, nil

				case key.Matches(msg, m.keys.preview), key.Matches(msg, m.keys.enter):
					if m.getSelectedSnippet() != nil {
						m.mode = previewMode
					}
					return m, nil

				case key.Matches(msg, m.keys.edit):
					snippet := m.getSelectedSnippet()
					if snippet != nil {
						m.mode = editMode
						m.editingID = snippet.ID
						m.titleInput.SetValue(snippet.Title)
						m.titleInput.Focus()
						m.bodyInput.SetValue(snippet.Body)
					}
					return m, nil

				case key.Matches(msg, m.keys.add):
					m.mode = addMode
					m.titleInput.SetValue("")
					m.titleInput.Focus()
					m.bodyInput.SetValue("")
					return m, nil

				case key.Matches(msg, m.keys.delete):
					m.deleteSelected()
					return m, nil
				}
			}

		case previewMode:
			switch {
			case key.Matches(msg, m.keys.back):
				m.mode = browseMode
				return m, nil
			case key.Matches(msg, m.keys.edit):
				snippet := m.getSelectedSnippet()
				if snippet != nil {
					m.mode = editMode
					m.editingID = snippet.ID
					m.titleInput.SetValue(snippet.Title)
					m.titleInput.Focus()
					m.bodyInput.SetValue(snippet.Body)
				}
				return m, nil
			}

		case editMode, addMode:
			switch msg.String() {
			case "esc":
				if m.titleInput.Focused() {
					m.mode = browseMode
					m.statusMsg = "Cancelled"
				} else {
					m.titleInput.Focus()
					m.bodyInput.Blur()
				}
				return m, nil

			case "tab":
				if m.titleInput.Focused() {
					m.titleInput.Blur()
					m.bodyInput.Focus()
				} else {
					m.bodyInput.Blur()
					m.titleInput.Focus()
				}
				return m, nil

			case "ctrl+s":
				title := strings.TrimSpace(m.titleInput.Value())
				body := strings.TrimSpace(m.bodyInput.Value())

				if title == "" {
					m.statusMsg = "Title cannot be empty"
					return m, nil
				}

				if m.mode == addMode {
					newSnippet := Snippet{
						ID:    m.nextID,
						Title: title,
						Body:  body,
					}
					m.items = append(m.items, newSnippet)
					m.nextID++
					m.statusMsg = fmt.Sprintf("Added: %s", title)
				} else {
					// Edit mode
					for i := range m.items {
						if m.items[i].ID == m.editingID {
							m.items[i].Title = title
							m.items[i].Body = body
							m.statusMsg = fmt.Sprintf("Updated: %s", title)
							break
						}
					}
				}

				m.applyFilter()
				m.mode = browseMode
				return m, nil
			}

			// Update inputs
			if m.titleInput.Focused() {
				m.titleInput, cmd = m.titleInput.Update(msg)
				cmds = append(cmds, cmd)
			} else {
				m.bodyInput, cmd = m.bodyInput.Update(msg)
				cmds = append(cmds, cmd)
			}
			return m, tea.Batch(cmds...)
		}
	}

	// Update table in browse mode
	if m.mode == browseMode && !m.filtering {
		m.table, cmd = m.table.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *browseModel) View() string {
	baseStyle := lipgloss.NewStyle().Padding(1, 2)
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("63")).
		Padding(0, 1)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Padding(1, 0, 0, 0)

	var content string

	switch m.mode {
	case browseMode:
		title := titleStyle.Render("📝 Snippet Manager")

		filterBar := ""
		if m.filtering {
			filterBar = "\n🔍 " + m.filterInput.View() + "\n"
		} else if m.filterQuery != "" {
			filterBar = fmt.Sprintf("\n🔍 Filter: %s (press esc to clear)\n", m.filterQuery)
		}

		statusBar := ""
		if m.statusMsg != "" {
			statusBar = lipgloss.NewStyle().
				Foreground(lipgloss.Color("42")).
				Render("\n✓ " + m.statusMsg)
		}

		help := ""
		if m.showHelp {
			help = helpStyle.Render("\nHelp:\n" +
				"  ↑/k, ↓/j: navigate  |  enter/p: preview  |  e: edit  |  a: add\n" +
				"  d: delete  |  /: filter  |  s: sort  |  ?: toggle help  |  q: quit")
		} else {
			help = helpStyle.Render("\nPress ? for help | " +
				fmt.Sprintf("%d snippets", len(m.filteredItems)))
		}

		content = title + "\n\n" + filterBar + m.table.View() + statusBar + help

	case previewMode:
		snippet := m.getSelectedSnippet()
		if snippet == nil {
			m.mode = browseMode
			return m.View()
		}

		title := titleStyle.Render(fmt.Sprintf("📖 Preview: %s", snippet.Title))

		bodyStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(1).
			Width(m.width - 8)

		body := bodyStyle.Render(snippet.Body)
		help := helpStyle.Render("\nPress esc/q to go back | e to edit")

		content = title + "\n\n" + body + help

	case editMode, addMode:
		modeTitle := "✏️  Edit Snippet"
		if m.mode == addMode {
			modeTitle = "➕ Add New Snippet"
		}
		title := titleStyle.Render(modeTitle)

		titleLabel := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render("Title:")
		bodyLabel := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render("Body:")

		help := helpStyle.Render("\ntab: switch fields | ctrl+s: save | esc: cancel")

		content = title + "\n\n" +
			titleLabel + "\n" + m.titleInput.View() + "\n\n" +
			bodyLabel + "\n" + m.bodyInput.View() + help
	}

	return baseStyle.Render(content)
}
