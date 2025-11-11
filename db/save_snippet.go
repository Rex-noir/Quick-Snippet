package db

import (
	"QuickSnip/db/models"
	"database/sql"
)

func SaveSnippet(conn *sql.DB, snippet models.Snippet) (int64, error) {
	result, err := conn.Exec(`UPDATE snippets SET title = ?, body = ? WHERE id =?`, snippet.Title, snippet.Body, snippet.ID)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}
