package service

import (
	"apigateway/pkg/repository"

	"go.uber.org/fx"
)

type service struct {
	repo repository.IRepository
}

// NewService ...
func NewService(repo repository.IRepository) IService {
	return &service{
		repo: repo,
	}
}

// IService ...
type IService interface {
}

// Module Export service module
var Module = fx.Options(
	fx.Provide(NewService),
)
