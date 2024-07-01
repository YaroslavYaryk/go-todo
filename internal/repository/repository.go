package repository

import (
	"github.com/jmoiron/sqlx"
	"simpleRestApi/internal/domain"
)

type Authorization interface {
	CreateUser(user domain.User) (int, error)
	GetUser(username, password string) (domain.User, error)
}

type TodoList interface {
	Create(userId int, list domain.TodoList) (int, error)
	GetAll(userId int) ([]domain.TodoListExtended, error)
	GetById(userId, listId int) (domain.TodoList, error)
	Update(listId int, input domain.TodoList) error
	Delete(listId int) (int, error)
	IsUserAuthorizedToUpdateList(listId int, userId int) (bool, error)
}
type TodoItem interface {
	Create(listId int, item domain.TodoItem) (int, error)
	GetAll(listId int, userId int) ([]domain.TodoItem, error)
	GetById(listId int, userId int, itemId int) (domain.TodoItem, error)
	Update(itemId int, input domain.TodoItem) error
	Delete(itemId int) (int, error)
}

type Repository struct {
	Authorization
	TodoList
	TodoItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		TodoList:      NewTodoListPostgres(db),
		TodoItem:      NewTodoItemPostgres(db),
	}
}
