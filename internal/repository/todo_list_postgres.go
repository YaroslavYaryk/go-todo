package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"simpleRestApi/internal/domain"
	"simpleRestApi/pkg/psql"
	"strings"
)

type TodoListPostgres struct {
	db *sqlx.DB
}

func NewTodoListPostgres(db *sqlx.DB) *TodoListPostgres {
	return &TodoListPostgres{
		db: db,
	}
}

func (r *TodoListPostgres) Create(userId int, list domain.TodoList) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	createListQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", psql.TodoListTable)
	row := tx.QueryRow(createListQuery, list.Title, list.Description)
	if err := row.Scan(&id); err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	createUsersListQuery := fmt.Sprintf("INSERT INTO %s (user_id, list_id) VALUES ($1, $2)", psql.UsersListsTable)
	_, err = tx.Exec(createUsersListQuery, userId, id)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()

}

func (r *TodoListPostgres) GetAll(userId int) ([]domain.TodoListExtended, error) {
	var lists []domain.TodoListExtended

	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description, ul.user_id as UserId FROM %s tl INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1",
		psql.TodoListTable, psql.UsersListsTable)
	err := r.db.Select(&lists, query, userId)
	return lists, err
}

func (r *TodoListPostgres) GetById(userId, listId int) (domain.TodoList, error) {
	var list domain.TodoList

	query := fmt.Sprintf(`SELECT tl.id, tl.title, tl.description FROM %s tl
								INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1 AND ul.list_id = $2`,
		psql.TodoListTable, psql.UsersListsTable)
	err := r.db.Get(&list, query, userId, listId)

	return list, err
}

func (r *TodoListPostgres) Update(listId int, input domain.TodoList) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != "" {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, input.Title)
		argId++
	}

	if input.Description != "" {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, input.Description)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`UPDATE %s SET %s WHERE id = $%d`, psql.TodoListTable, setQuery, argId)
	args = append(args, listId)

	_, err := r.db.Exec(query, args...)

	logrus.Debugf("updated todo list with id %d", listId)

	return err
}

func (r *TodoListPostgres) Delete(listId int) (int, error) {

	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, psql.TodoListTable)
	_, err = tx.Exec(query, listId)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	query = fmt.Sprintf(`DELETE FROM %s WHERE list_id = $1`, psql.UsersListsTable)
	_, err = tx.Exec(query, listId)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return listId, nil
}

func (r *TodoListPostgres) IsUserAuthorizedToUpdateList(listId int, userId int) (bool, error) {
	var count int
	query := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE user_id = $1 AND list_id = $2`, psql.UsersListsTable)
	err := r.db.QueryRow(query, userId, listId).Scan(&count)

	if err != nil {
		return false, nil
	}
	return count > 0, nil
}
