package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"simpleRestApi/internal/domain"
	"simpleRestApi/pkg/psql"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func (r *AuthPostgres) CreateUser(user domain.User) (int, error) {

	var id int
	query := fmt.Sprintf("INSERT INTO %s (name, username, password_hash) VALUES ($1, $2, $3) RETURNING id", psql.UsersTable)
	row := r.db.QueryRow(query, user.Name, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *AuthPostgres) GetUser(username, password string) (domain.User, error) {

	var user domain.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE username=$1 AND password_hash=$2", psql.UsersTable)
	err := r.db.Get(&user, query, username, password)

	return user, err
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}
