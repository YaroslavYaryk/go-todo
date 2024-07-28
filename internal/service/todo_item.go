package service

import (
	"simpleRestApi/internal/domain"
	"simpleRestApi/internal/repository"
)

type TodoItemService struct {
	repo repository.TodoItem
}

func newTodoItemService(repo repository.TodoItem) *TodoItemService {
	return &TodoItemService{
		repo: repo,
	}
}

func (s *TodoItemService) Create(listId int, item domain.TodoItem, user domain.UserGet) (int, error) {
	return s.repo.Create(listId, item, user)
}

func (s *TodoItemService) GetAll(listId int, userId int) ([]domain.TodoItemGet, error) {
	return s.repo.GetAll(listId, userId)
}

func (s *TodoItemService) GetById(listId int, userId int, itemId int) (domain.TodoItemGet, error) {
	return s.repo.GetById(listId, userId, itemId)
}

func (s *TodoItemService) Update(itemId int, input domain.TodoItem, userId int) (bool, error) {
	return s.repo.Update(itemId, input, userId)
}

func (s *TodoItemService) Delete(itemId int, userId int) (int, error) {
	return s.repo.Delete(itemId, userId)
}
