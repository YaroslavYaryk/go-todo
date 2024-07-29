package service

import (
	"simpleRestApi/internal/domain"
	"simpleRestApi/internal/repository"
)

type Authorization interface {
	CreateUser(user domain.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
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
	Create(userId int, item domain.TodoItem, user domain.UserGet) (int, error)
	GetAll(listId int, userId int) ([]domain.TodoItemGet, error)
	GetById(listId int, userId int, itemId int) (domain.TodoItemGet, error)
	Update(itemId int, input domain.TodoItem, userId int) (bool, error)
	Delete(listId int, userId int) (int, error)
}

type MiddleWare interface {
	GetUserById(userId int) (domain.UserGet, error)
}

type Note interface {
	Create(userId int, note domain.Note) (int, error)
	GetAll(userId int) ([]domain.Note, error)
	GetById(userId int, noteId int) (domain.Note, error)
	Update(noteId int, input domain.Note) error
	Delete(noteId int, userId int) (int, error)
	ShareNote(userId int, noteId int, sharedUserId int) error
}

type Service struct {
	Authorization
	TodoList
	TodoItem
	MiddleWare
	Note
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		TodoList:      newTodoListService(repos.TodoList),
		TodoItem:      newTodoItemService(repos.TodoItem),
		MiddleWare:    NewMiddlewareService(repos.MiddleWare),
		Note:          NewNoteService(repos.Note),
	}
}
