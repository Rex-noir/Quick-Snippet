package ui

import "github.com/charmbracelet/bubbles/key"

type browseKeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Add    key.Binding
	Edit   key.Binding
	Delete key.Binding
	Filter key.Binding
	Sort   key.Binding
	Help   key.Binding
	Quit   key.Binding
	Enter  key.Binding
	Escape key.Binding
	Copy   key.Binding
}

func newBrowseKeyMap() browseKeyMap {
	return browseKeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		Add: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add snippet"),
		),
		Edit: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit"),
		),
		Delete: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "delete"),
		),
		Filter: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "filter"),
		),
		Sort: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "sort"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "confirm"),
		),
		Escape: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel"),
		),
		Copy: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "copy"),
		),
	}
}

// ShortHelp returns a quick help string
func (k browseKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Add, k.Edit, k.Delete, k.Filter, k.Sort, k.Copy, k.Help, k.Quit}
}

// FullHelp returns extended help
func (k browseKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Add, k.Edit, k.Delete},
		{k.Filter, k.Sort, k.Copy, k.Help, k.Quit},
	}
}
