package service

import (
	"simpleRestApi/internal/domain"
	"simpleRestApi/internal/repository"
)

type TodoListService struct {
	repo repository.TodoList
}

func newTodoListService(repo repository.TodoList) *TodoListService {
	return &TodoListService{
		repo: repo,
	}
}

func (s *TodoListService) Create(userId int, list domain.TodoList) (int, error) {
	return s.repo.Create(userId, list)
}

func (s *TodoListService) GetAll(userId int) ([]domain.TodoListExtended, error) {
	return s.repo.GetAll(userId)
}

func (s *TodoListService) GetById(userId, listId int) (domain.TodoList, error) {
	return s.repo.GetById(userId, listId)
}

func (s *TodoListService) Update(listId int, input domain.TodoList) error {
	return s.repo.Update(listId, input)
}

func (s *TodoListService) Delete(listId int) (int, error) {
	return s.repo.Delete(listId)
}

func (s *TodoListService) IsUserAuthorizedToUpdateList(userId int, listId int) (bool, error) {
	return s.repo.IsUserAuthorizedToUpdateList(userId, listId)
}
