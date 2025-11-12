package db

import (
	"QuickSnip/db/models"
	"database/sql"
)

func FetchFirst(db *sql.DB, n int) ([]models.Snippet, error) {

	rows, err := db.Query("SELECT id, title, body FROM snippets ORDER BY id ASC LIMIT ?", n)
	if err != nil {
		return nil, err
	}

	var snippets []models.Snippet
	for rows.Next() {
		var s models.Snippet
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
