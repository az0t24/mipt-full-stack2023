package service

import (
	"scoreboardpro/internal/entity"
	"scoreboardpro/internal/repository"
)

type Service struct {
	User  entity.UserService
	Table entity.TableService
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		User:  NewUserService(repo.User, repo.Table),
		Table: NewTableService(repo.Table),
	}
}
