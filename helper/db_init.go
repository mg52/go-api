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

	sqlStatementCreateTable2 := `
CREATE TABLE IF NOT EXISTS "todos" (
 id SERIAL PRIMARY KEY,
 user_id INT NOT NULL,
 content VARCHAR (255) NOT NULL,
    CONSTRAINT FK_user_id FOREIGN KEY(user_id)
        REFERENCES users(id)
);`

	_, err = db.Exec(sqlStatementCreateTable2)
	if err != nil {
		return err
	}

	return nil
}
