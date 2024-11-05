package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"simpleRestApi/internal/domain"
	"simpleRestApi/pkg/psql"
	"strings"
)

type UserPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{
		db: db,
	}
}

func (r *UserPostgres) GetUserById(userId int) (domain.UserGet, error) {

	var user domain.UserGet

	query := fmt.Sprintf(`select id, name, username, theme, timezone, is_paid_member as IsPaidMember, users.rate as RateId, task_color as TaskColor from %s where id = $1`,
		psql.UsersTable)

	err := r.db.Get(&user, query, userId)

	return user, err
}

func (r *UserPostgres) UpdateUser(input domain.UserGet, userId int) error {

	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Name != "" {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argId))
		args = append(args, input.Name)
		argId++
	}

	if input.Theme != "" {
		setValues = append(setValues, fmt.Sprintf("theme=$%d", argId))
		args = append(args, input.Theme)
		argId++
	}

	if input.Timezone != nil {
		setValues = append(setValues, fmt.Sprintf("timezone=$%d", argId))
		args = append(args, input.Timezone)
		argId++
	}

	if input.IsPaidMember {
		setValues = append(setValues, fmt.Sprintf("is_paid_member=$%d", argId))
		args = append(args, input.IsPaidMember)
		argId++
	}

	if input.RateId != nil {
		setValues = append(setValues, fmt.Sprintf("rate_id=$%d", argId))
		args = append(args, input.RateId)
		argId++
	}

	if input.TaskColor != nil {
		setValues = append(setValues, fmt.Sprintf("task_color=$%d", argId))
		args = append(args, input.TaskColor)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`UPDATE %s SET %s WHERE id = $%d`, psql.UsersTable, setQuery, argId)
	args = append(args, userId)

	_, err := r.db.Exec(query, args...)

	logrus.Debugf("updated todo list with id %d", userId)

	return err
}
