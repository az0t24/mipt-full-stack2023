package sqlite

import (
	"errors"
	"gorm.io/gorm"
	"scoreboardpro/internal/entity"
	"scoreboardpro/internal/entity/model/football"
)

type ResultSQLite struct {
	db *gorm.DB
}

func NewResultSQLite(db *gorm.DB) *ResultSQLite {
	return &ResultSQLite{db: db}

	//fixme: test only
	//r := &ResultSQLite{db: db}
	//r.Create(&football.Result{
	//	SeasonTableID:
	//})
}

func (r *ResultSQLite) Create(result *football.Result) error {
	return r.db.Create(result).Error
}

func (r *ResultSQLite) Update(result *football.Result) error {
	return r.db.Model(result).Updates(result).Error
}

func (r *ResultSQLite) Delete(id uint) error {
	result := r.db.Delete(&football.Result{}, id)
	if result.Error != nil {
		return result.Error
	} else {
		return nil
	}
}
func (r *ResultSQLite) GetAll() (*[]football.Result, error) {
	var res []football.Result

	if result := r.db.Find(&res); result.Error != nil {
		return nil, result.Error
	} else {
		return &res, nil
	}
}

func (r *ResultSQLite) Get(id uint) (*football.Result, error) {
	var res football.Result

	if result := r.db.Where("id = ?", id).First(&res); result.Error == nil {
		return &res, nil
	} else if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &res, entity.ErrUserNotFound
	} else {
		return &res, result.Error
	}
}

func (r *ResultSQLite) GetByTeam(team string) (*football.Result, error) {
	var res football.Result

	if result := r.db.Where("teamName = ?", team).First(&res); result.Error == nil {
		return &res, nil
	} else if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &res, entity.ErrUserNotFound
	} else {
		return &res, result.Error
	}
}
