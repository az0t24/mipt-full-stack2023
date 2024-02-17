package entity

import (
	"scoreboardpro/internal/entity/model/football"
)

type TableRepository interface {
	Create(table *football.SeasonTable) error
	GetAll() (*[]football.SeasonTable, error)
	Get(id uint) (*football.SeasonTable, error)
	GetBySeasonAndName(name, season string) (*football.SeasonTable, error)
	Update(*football.SeasonTable) error
	Delete(id uint) error
}

type TableService interface {
	Get(id uint) (*football.SeasonTable, error)
	GetAll() (*[]football.SeasonTable, error)
	GetByNameAndSeason(name, season string) (*football.SeasonTable, error)
}
