package service

import (
	"scoreboardpro/internal/entity"
	"scoreboardpro/internal/entity/model/football"
)

type TableService struct {
	repo entity.TableRepository
}

func NewTableService(tableRepo entity.TableRepository) *TableService {
	return &TableService{
		repo: tableRepo,
	}
}

func (s *TableService) Get(id uint) (*football.SeasonTable, error) {
	table, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	} else {
		return table, nil
	}
}

func (s *TableService) GetByNameAndSeason(name, season string) (*football.SeasonTable, error) {
	table, err := s.repo.GetBySeasonAndName(name, season)
	if err != nil {
		return nil, err
	} else {
		return table, nil
	}
}
func (s *TableService) GetAll() (*[]football.SeasonTable, error) {
	tablesDB, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	return tablesDB, nil
}
