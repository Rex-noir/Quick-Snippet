package db

import "database/sql"

func CreateSnippet(db *sql.DB, title, body string) (int64, error) {
	result, err := db.Exec(`INSERT INTO snippets(title, body) VALUES (?, ?)`, title, body)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil

}
