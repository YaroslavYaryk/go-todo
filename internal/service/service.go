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
	Create(userId int, item domain.TodoItem) (int, error)
	GetAll(listId int, userId int) ([]domain.TodoItem, error)
	GetById(listId int, userId int, itemId int) (domain.TodoItem, error)
	Update(itemId int, input domain.TodoItem) error
	Delete(listId int) (int, error)
}

type Service struct {
	Authorization
	TodoList
	TodoItem
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		TodoList:      newTodoListService(repos.TodoList),
		TodoItem:      newTodoItemService(repos.TodoItem),
	}
}
