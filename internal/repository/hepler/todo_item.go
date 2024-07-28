package hepler

import (
	"database/sql"
	"fmt"
	"simpleRestApi/internal/domain"
	"simpleRestApi/pkg/psql"
)

func CalculateCleanerCompletedTasks(tx *sql.Tx, userId int) (int, error) {
	query := fmt.Sprintf(
		`select Count(ti.id) from %s ti 
					join %s li ON li.item_id = ti.id 
					join %s ul on li.list_id = ul.list_id
					where user_id = $1 and ti.is_deleted = False and ti.done = true`, psql.TodoItemTable, psql.ListItemTable, psql.UsersListsTable)

	var count int
	err := tx.QueryRow(query, userId).Scan(&count)
	if err != nil {

		_ = tx.Rollback()
		return 0, err
	}

	return count, nil
}

func GetUserRate(tx *sql.Tx, userId int) (sql.NullInt64, error) {
	var oldRate sql.NullInt64
	query := fmt.Sprintf(`select rate from %s where id = $1`, psql.UsersTable)
	err := tx.QueryRow(query, userId).Scan(&oldRate)
	if err != nil {

		_ = tx.Rollback()
		return sql.NullInt64{Valid: true}, err
	}

	return oldRate, nil
}

func RecalculateUserRate(tx *sql.Tx, userId int) (bool, error) {

	updatedRate := false

	//	calculate user completed tasks
	count, err := CalculateCleanerCompletedTasks(tx, userId)
	if err != nil {
		return false, err
	}

	fmt.Println(count)

	//	get user old rate
	oldRate, err := GetUserRate(tx, userId)
	if err != nil {
		return false, err
	}

	// recalculate cleaner rate
	var newRate sql.NullInt64
	query := fmt.Sprintf(`select id from %s where task_completed <= $1`, psql.RateTable)
	err = tx.QueryRow(query, count).Scan(&newRate)
	if err != nil {
		return false, err
	}

	if oldRate.Int64 != newRate.Int64 {
		updatedRate = true
		query = fmt.Sprintf(`UPDATE %s SET rate = $1 WHERE id = $2`, psql.UsersTable)
		_, err = tx.Exec(query, newRate.Int64, userId)
	}

	return updatedRate, nil
}

func GetCategoryByName(tx *sql.Tx, name string) (domain.Category, error) {
	var category domain.Category
	query := fmt.Sprintf(`select id, name, icon_name from %s where name = $1`, psql.CategoryTable)

	err := tx.QueryRow(query, name).Scan(&category.Id, &category.Name, &category.IconName)

	if err != nil {
		return domain.Category{}, err
	}

	return category, nil
}

func BuildHierarchy(items []domain.TodoItemGet) []domain.TodoItemGet {
	itemMap := make(map[int]*domain.TodoItemGet)
	var roots []domain.TodoItemGet

	for i := range items {
		item := &items[i]
		itemMap[item.Id] = item
	}

	for i := range items {
		item := &items[i]
		if item.ParentID != nil {
			parent := itemMap[*item.ParentID]
			parent.Children = append(parent.Children, *item)
		}
	}

	for i := range items {
		item := &items[i]
		if item.ParentID == nil {
			roots = append(roots, *item)
		}
	}

	return roots
}
