package service

import (
	"simpleRestApi/internal/domain"
	"simpleRestApi/internal/repository"
)

type MiddlewareService struct {
	repo repository.MiddleWare
}

func NewMiddlewareService(repo repository.MiddleWare) *MiddlewareService {
	return &MiddlewareService{repo: repo}
}

func (s *MiddlewareService) GetUserById(userId int) (domain.UserGet, error) {

	return s.repo.GetUserById(userId)
}
