package helper

import "database/sql"

func CreateDatabaseObjects(db *sql.DB) error {
	sqlStatementCreateTable := `
CREATE TABLE IF NOT EXISTS "users" (
 id SERIAL PRIMARY KEY,
 username TEXT UNIQUE NOT NULL,
 password TEXT NOT NULL
);`

	_, err := db.Exec(sqlStatementCreateTable)
	if err != nil {
		return err
	}

	return nil
}
