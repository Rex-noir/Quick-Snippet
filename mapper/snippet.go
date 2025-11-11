package mapper

import (
	"QuickSnip/db/models"
	"QuickSnip/ui"
)

func ToUISnippet(m models.Snippet) ui.Snippet {
	return ui.Snippet{
		ID:    m.ID,
		Title: m.Title,
		Body:  m.Body,
	}
}

func ToModelSnippet(m ui.Snippet) models.Snippet {
	return models.Snippet{
		ID:    m.ID,
		Title: m.Title,
	}
}

func ToModelSnippets(m []ui.Snippet) []models.Snippet {
	snippets := make([]models.Snippet, len(m))
	for i, s := range m {
		snippets[i] = ToModelSnippet(s)
	}
	return snippets
}

func ToUISnippets(m []models.Snippet) []ui.Snippet {
	snippets := make([]ui.Snippet, len(m))
	for i, s := range m {
		snippets[i] = ToUISnippet(s)
	}
	return snippets
}
