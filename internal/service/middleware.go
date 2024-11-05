package service

import (
	"simpleRestApi/internal/repository"
)

type MiddlewareService struct {
	repo repository.MiddleWare
}

func NewMiddlewareService(repo repository.MiddleWare) *MiddlewareService {
	return &MiddlewareService{repo: repo}
}
