package sqlite

import (
	"gorm.io/gorm"
	"scoreboardpro/internal/entity/model/football"
)

type TableSQLite struct {
	db *gorm.DB
}

func NewTableSQLite(db *gorm.DB) *TableSQLite {
	//return &TableSQLite{db: db}

	//fixme test only:
	t := &TableSQLite{db: db}
	t.Create(&football.SeasonTable{
		ChampionshipName: "Bundesliga",
		Season:           "2018-2019",
		Teams: []football.Result{
			{
				TeamName: "Dortmund",
				Wins:     15,
				Draws:    0,
				Loses:    8,
			},
			{
				TeamName: "Munich",
				Wins:     11,
				Draws:    1,
				Loses:    5,
			},
		},
	})
	t.Create(&football.SeasonTable{
		ChampionshipName: "Bundesliga",
		Season:           "2017-2018",
		Teams: []football.Result{
			{
				TeamName: "Bayern04",
				Wins:     5,
				Draws:    0,
				Loses:    0,
			},
			{
				TeamName: "RB",
				Wins:     0,
				Draws:    1,
				Loses:    1,
			},
		},
	})
	return t
}

func (r *TableSQLite) Create(result *football.SeasonTable) error {
	return r.db.Create(result).Error
}

func (r *TableSQLite) Update(result *football.SeasonTable) error {
	return r.db.Model(result).Updates(result).Error
}

func (r *TableSQLite) Delete(id uint) error {
	result := r.db.Delete(&football.SeasonTable{}, id)
	if result.Error != nil {
		return result.Error
	} else {
		return nil
	}
}
func (r *TableSQLite) GetAll() (*[]football.SeasonTable, error) {
	var res []football.SeasonTable

	if result := r.db.Preload("Teams").Find(&res); result.Error != nil {
		return nil, result.Error
	} else {
		return &res, nil
	}
}

func (r *TableSQLite) Get(id uint) (*football.SeasonTable, error) {
	var res football.SeasonTable

	if result := r.db.Preload("Teams").Where("id = ?", id).First(&res); result.Error == nil {
		return &res, nil
	} else {
		return &res, result.Error
	}
}

func (r *TableSQLite) GetBySeasonAndName(name, season string) (*football.SeasonTable, error) {
	var res football.SeasonTable

	if result := r.db.Preload("Teams").Where("championship_name = ? AND season = ?", name, season).First(&res); result.Error == nil {
		return &res, nil
	} else {
		return &res, result.Error
	}
}
