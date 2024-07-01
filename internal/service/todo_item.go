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

func (s *TodoItemService) Create(listId int, item domain.TodoItem) (int, error) {
	return s.repo.Create(listId, item)
}

func (s *TodoItemService) GetAll(listId int, userId int) ([]domain.TodoItem, error) {
	return s.repo.GetAll(listId, userId)
}

func (s *TodoItemService) GetById(listId int, userId int, itemId int) (domain.TodoItem, error) {
	return s.repo.GetById(listId, userId, itemId)
}

func (s *TodoItemService) Update(itemId int, input domain.TodoItem) error {
	return s.repo.Update(itemId, input)
}

func (s *TodoItemService) Delete(itemId int) (int, error) {
	return s.repo.Delete(itemId)
}
