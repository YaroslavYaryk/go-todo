package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"simpleRestApi/internal/domain"
	"simpleRestApi/internal/repository/hepler"
	"simpleRestApi/pkg/psql"
	"strings"
	"time"
)

type TodoItemPostgres struct {
	db *sqlx.DB
}

func NewTodoItemPostgres(db *sqlx.DB) *TodoItemPostgres {
	return &TodoItemPostgres{
		db: db,
	}
}

func (r *TodoItemPostgres) Create(listId int, item domain.TodoItem, user domain.UserGet) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var userNotification *time.Time
	if user.IsPaidMember {
		userNotification = item.NotificationTime
	} else {
		userNotification = nil
	}

	category := item.CategoryId
	if item.CategoryId == nil {
		categoryObj, err := hepler.GetCategoryByName(tx, "Empty")
		if err != nil {
			_ = tx.Rollback()
			return 0, err
		}
		category = &categoryObj.Id
	}

	var id int
	createItemQuery := fmt.Sprintf("INSERT INTO %s (title, description, done, created_at, updated_at, note, notification_time, predicted_time_to_spend, is_deleted, priority, category, parent_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id", psql.TodoItemTable)

	row := tx.QueryRow(createItemQuery, item.Title, item.Description, false, time.Now(), time.Now(), item.Note, userNotification, item.PredictedTimeToSpend, false, item.Priority, category, item.ParentID)
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

func (r *TodoItemPostgres) GetAll(listId int, userId int) ([]domain.TodoItemGet, error) {
	var items []domain.TodoItemGet

	query := fmt.Sprintf(`SELECT 
        ti.id, ti.title, ti.description, ti.done, ti.created_at as CreatedAt, ti.updated_at as UpdatedAt, 
        ti.note as Note, ti.notification_time as NotificationTime, ti.predicted_time_to_spend as PredictedTimeToSpend, ti.priority as Priority, 
        ti.is_deleted as IsDeleted, ti.parent_id as ParentId,
        c.id AS "category.id", c.name AS "category.name", c.icon_name as "category.iconname"
    FROM 
        %s ti
    JOIN 
        %s c ON ti.category = c.id
	JOIN 
	   %s li on ti.id = li.item_id
	JOIN 
		%s ul on li.list_id = ul.list_id
	WHERE
	    ul.user_id = $1 and li.list_id = $2 and is_deleted = False
	ORDER BY
	    ti.priority
	`, psql.TodoItemTable, psql.CategoryTable, psql.ListItemTable, psql.UsersListsTable)

	err := r.db.Select(&items, query, userId, listId)

	hierarchy := hepler.BuildHierarchy(items)

	return hierarchy, err

}

func (r *TodoItemPostgres) GetById(listId int, userId int, itemId int) (domain.TodoItemGet, error) {
	var items []domain.TodoItemGet

	query := fmt.Sprintf(`SELECT 
        ti.id, ti.title, ti.description, ti.done, ti.created_at as CreatedAt, ti.updated_at as UpdatedAt, 
        ti.note as Note, ti.notification_time as NotificationTime, ti.predicted_time_to_spend as PredictedTimeToSpend, ti.priority as Priority, 
        ti.is_deleted as IsDeleted, ti.parent_id as ParentId,
        c.id AS "category.id", c.name AS "category.name", c.icon_name as "category.iconname"
    FROM 
        %s ti
    LEFT JOIN 
        %s c ON ti.category = c.id
	JOIN 
	   %s li on ti.id = li.item_id
	JOIN 
		%s ul on li.list_id = ul.list_id
	WHERE
	    ul.user_id = $1 and li.list_id = $2 and (ti.id = $3 or ti.parent_id = $3)
	`, psql.TodoItemTable, psql.CategoryTable, psql.ListItemTable, psql.UsersListsTable)
	err := r.db.Select(&items, query, userId, listId, itemId)

	if len(items) > 0 {
		hierarchy := hepler.BuildHierarchy(items)
		return hierarchy[0], err
	}

	return items[0], err
}

func (r *TodoItemPostgres) Update(itemId int, input domain.TodoItem, userId int) (bool, error) {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1
	updatedRate := false

	tx, err := r.db.Begin()
	if err != nil {
		return false, err
	}

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

	if input.Priority != nil {
		setValues = append(setValues, fmt.Sprintf("priority=$%d", argId))
		args = append(args, &input.Priority)
		argId++
	}

	if input.Note != nil {
		setValues = append(setValues, fmt.Sprintf("note=$%d", argId))
		args = append(args, &input.Note)
		argId++
	}

	if input.NotificationTime != nil {
		setValues = append(setValues, fmt.Sprintf("notification_time=$%d", argId))
		args = append(args, &input.NotificationTime)
		argId++
	}

	if input.PredictedTimeToSpend != nil {
		setValues = append(setValues, fmt.Sprintf("predicted_time_to_spend=$%d", argId))
		args = append(args, &input.PredictedTimeToSpend)
		argId++
	}

	if input.CategoryId != nil {
		setValues = append(setValues, fmt.Sprintf("category=$%d", argId))
		args = append(args, &input.CategoryId)
		argId++
	}

	// change updated_at
	setValues = append(setValues, fmt.Sprintf("updated_at=$%d", argId))
	args = append(args, time.Now())
	argId++

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`UPDATE %s SET %s WHERE id = $%d`, psql.TodoItemTable, setQuery, argId)
	args = append(args, itemId)

	_, err = tx.Exec(query, args...)
	if err != nil {

		_ = tx.Rollback()
		return false, err
	}

	if input.Done != nil {
		updatedRate, err = hepler.RecalculateUserRate(tx, userId)
		if err != nil {
			_ = tx.Rollback()
			return false, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return false, err
	}

	logrus.Debugf("updated todo item with id %d", itemId)

	return updatedRate, err
}

func (r *TodoItemPostgres) Delete(itemId int, userId int) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	//query := fmt.Sprintf(`DELETE FROM %s WHERE item_id = $1`, psql.ListItemTable)
	//_, err = tx.Exec(query, itemId)
	//if err != nil {
	//	_ = tx.Rollback()
	//	return 0, err
	//}

	query := fmt.Sprintf(`UPDATE %s SET is_deleted = True WHERE id = $1`, psql.TodoItemTable)
	_, err = tx.Exec(query, itemId)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	_, err = hepler.RecalculateUserRate(tx, userId)
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
