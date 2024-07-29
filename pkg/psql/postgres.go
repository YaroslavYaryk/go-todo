package psql

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

const (
	UsersTable      = "users"
	TodoListTable   = "todo_lists"
	UsersListsTable = "users_lists"
	TodoItemTable   = "todo_items"
	ListItemTable   = "lists_items"
	CategoryTable   = "category"
	RateTable       = "rate"
	NoteTable       = "notes"
	UserNoteTable   = "user_note"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", cfg.Host, cfg.Port, cfg.Username, cfg.Database, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
