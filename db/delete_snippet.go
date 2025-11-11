package db

import "database/sql"

func DeleteSnippet(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM snippets WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}
