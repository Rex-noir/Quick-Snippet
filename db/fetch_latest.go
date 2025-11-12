package db

import (
	"QuickSnip/db/models"
	"database/sql"
)

func FetchLatest(db *sql.DB, n *int) ([]models.Snippet, error) {
	var number int
	if n == nil {
		number = 1
	} else {
		number = *n
	}

	rows, err := db.Query("SELECT id,title,body FROM snippets ORDER BY id DESC LIMIT ?", number)
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
