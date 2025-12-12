package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func NewDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS participants (
        chat_id INTEGER,
        user_id INTEGER,
        username TEXT,
        PRIMARY KEY(chat_id, user_id)
    );

    CREATE TABLE IF NOT EXISTS assignments (
        chat_id INTEGER,
        giver_id INTEGER,
        receiver_id INTEGER
    );
    `)
	if err != nil {
		return nil, err
	}

	return db, nil
}
