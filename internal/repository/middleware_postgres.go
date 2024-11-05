package repository

import (
	"github.com/jmoiron/sqlx"
)

type MiddlewarePostgres struct {
	db *sqlx.DB
}

func NewMiddlewarePostgres(db *sqlx.DB) *MiddlewarePostgres {
	return &MiddlewarePostgres{
		db: db,
	}
}
