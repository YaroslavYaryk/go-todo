package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"simpleRestApi/internal/domain"
	"simpleRestApi/pkg/psql"
)

type MiddlewarePostgres struct {
	db *sqlx.DB
}

func NewMiddlewarePostgres(db *sqlx.DB) *MiddlewarePostgres {
	return &MiddlewarePostgres{
		db: db,
	}
}

func (r *MiddlewarePostgres) GetUserById(userId int) (domain.UserGet, error) {

	var user domain.UserGet

	query := fmt.Sprintf(`select id, username, theme, timezone, is_paid_member as IsPaidMember, users.rate as RateId from %s where id = $1`,
		psql.UsersTable)

	err := r.db.Get(&user, query, userId)

	return user, err
}
