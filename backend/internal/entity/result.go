package entity

import (
	"scoreboardpro/internal/entity/model/football"
)

type ResultRepository interface {
	Create(result *football.Result) error
	GetAll() (*[]football.Result, error)
	Get(id uint) (*football.Result, error)
	GetByTeam(team string) (*football.Result, error)
	Update(*football.Result) error
	Delete(id uint) error
}

type ResultService interface {
}
