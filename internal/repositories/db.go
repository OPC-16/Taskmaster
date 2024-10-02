package repositories

import (
	"database/sql"
	"fmt"

   _ "github.com/mattn/go-sqlite3"
)

func NewSQLiteDB(dbPath string) (*sql.DB, error) {
   db, err := sql.Open("sqlite3", dbPath)
   if err != nil {
      return nil, fmt.Errorf("failed to open database: %w", err)
   }

   if err := db.Ping(); err != nil {
      return nil, fmt.Errorf("failed to ping the database: %w", err)
   }

   if err := createTables(db); err != nil {
      return nil, fmt.Errorf("failed to create tables: %w", err)
   }

   return db, nil
}

func createTables(db *sql.DB) error {
   userTable := `CREATE TABLE IF NOT EXISTS users (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      username TEXT NOT NULL UNIQUE,
      email TEXT NOT NULL UNIQUE,
      password TEXT NOT NULL
   );`

   taskTable := `create table if not exists tasks (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      user_id INTEGER NOT NULL,
      title TEXT NOT NULL,
      description TEXT,
      category TEXT,
      deadline DATETIME,
      priority INTEGER,
      status TEXT,
      created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
      updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
      FOREIGN KEY (user_id) REFERENCES users (id)
   );`

   for _, table := range []string{userTable, taskTable} {
      _, err := db.Exec(table)
      if err != nil {
         return err
      }
   }

   return nil
}
