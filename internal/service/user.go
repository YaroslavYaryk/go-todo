package service

import (
	"simpleRestApi/internal/domain"
	"simpleRestApi/internal/repository"
)

type UserService struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUserById(userId int) (domain.UserGet, error) {

	return s.repo.GetUserById(userId)
}

func (s *UserService) UpdateUser(input domain.UserGet, userId int) error {

	return s.repo.UpdateUser(input, userId)
}
