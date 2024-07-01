package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"simpleRestApi/internal/domain"
	"simpleRestApi/pkg/psql"
	"strings"
)

type TodoItemPostgres struct {
	db *sqlx.DB
}

func NewTodoItemPostgres(db *sqlx.DB) *TodoItemPostgres {
	return &TodoItemPostgres{
		db: db,
	}
}

func (r *TodoItemPostgres) Create(listId int, item domain.TodoItem) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	createItemQuery := fmt.Sprintf("INSERT INTO %s (title, description, done) VALUES ($1, $2, $3) RETURNING id", psql.TodoItemTable)
	row := tx.QueryRow(createItemQuery, item.Title, item.Description, false)
	if err := row.Scan(&id); err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	createListItemQuery := fmt.Sprintf("INSERT INTO %s (list_id, item_id) VALUES ($1, $2)", psql.ListItemTable)
	_, err = tx.Exec(createListItemQuery, listId, id)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()

}

func (r *TodoItemPostgres) GetAll(listId int, userId int) ([]domain.TodoItem, error) {
	var items []domain.TodoItem

	query := fmt.Sprintf(`select ti.id, ti.title, ti.description, ti.done from %s ti
    							join %s li on ti.id = li.item_id    
    							join %s ul on li.list_id = ul.list_id        
                                    where ul.user_id = $1 and li.list_id = $2`,
		psql.TodoItemTable, psql.ListItemTable, psql.UsersListsTable)
	err := r.db.Select(&items, query, userId, listId)
	return items, err
}

func (r *TodoItemPostgres) GetById(listId int, userId int, itemId int) (domain.TodoItem, error) {
	var item domain.TodoItem

	query := fmt.Sprintf(`select ti.id, ti.title, ti.description, ti.done from %s ti
    							join %s li on ti.id = li.item_id    
    							join %s ul on li.list_id = ul.list_id        
                                    where ul.user_id = $1 and li.list_id = $2 and ti.id = $3`,
		psql.TodoItemTable, psql.ListItemTable, psql.UsersListsTable)
	err := r.db.Get(&item, query, userId, listId, itemId)
	return item, err
}

func (r *TodoItemPostgres) Update(itemId int, input domain.TodoItem) error {
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

	if input.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done=$%d", argId))
		args = append(args, &input.Done)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`UPDATE %s SET %s WHERE id = $%d`, psql.TodoItemTable, setQuery, argId)
	args = append(args, itemId)

	_, err := r.db.Exec(query, args...)

	logrus.Debugf("updated todo item with id %d", itemId)

	return err
}

func (r *TodoItemPostgres) Delete(itemId int) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	query := fmt.Sprintf(`DELETE FROM %s WHERE item_id = $1`, psql.ListItemTable)
	_, err = tx.Exec(query, itemId)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	query = fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, psql.TodoItemTable)
	_, err = tx.Exec(query, itemId)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return itemId, nil
}
