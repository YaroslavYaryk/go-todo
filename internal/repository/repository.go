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
	Create(listId int, item domain.TodoItem, user domain.UserGet) (int, error)
	GetAll(listId int, userId int) ([]domain.TodoItemGet, error)
	GetById(listId int, userId int, itemId int) (domain.TodoItemGet, error)
	Update(itemId int, input domain.TodoItem, userId int) (bool, error)
	Delete(itemId int, userId int) (int, error)
}

type MiddleWare interface {
}

type Note interface {
	Create(userId int, note domain.Note) (int, error)
	GetAll(userId int) ([]domain.Note, error)
	GetById(userId int, noteId int) (domain.Note, error)
	Update(noteId int, input domain.Note) error
	Delete(noteId int, userId int) (int, error)
	ShareNote(userId int, noteId int, sharedUserId int) error
}

type User interface {
	GetUserById(userId int) (domain.UserGet, error)
	UpdateUser(input domain.UserGet, userId int) error
}

type Repository struct {
	Authorization
	TodoList
	TodoItem
	MiddleWare
	Note
	User
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		TodoList:      NewTodoListPostgres(db),
		TodoItem:      NewTodoItemPostgres(db),
		MiddleWare:    NewMiddlewarePostgres(db),
		Note:          NewNotePostgres(db),
		User:          NewUserPostgres(db),
	}
}
