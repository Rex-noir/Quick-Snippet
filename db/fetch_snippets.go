package db

import (
	"QuickSnip/ui"
	"database/sql"
)

func FetchSnippets(conn *sql.DB) ([]ui.Snippet, error) {
	rows, err := conn.Query("SELECT id, title, body FROM snippets ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	var snippets []ui.Snippet
	for rows.Next() {
		var s ui.Snippet
		if err := rows.Scan(&s.ID, &s.Title, &s.Body); err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil

}
