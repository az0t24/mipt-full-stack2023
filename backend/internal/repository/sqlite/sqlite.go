package sqlite

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"scoreboardpro/internal/entity/model"
	"scoreboardpro/internal/entity/model/football"
)

func NewSQLiteDB(dbUri string) (*gorm.DB, error) {
	db, err := gorm.Open(
		sqlite.Open(dbUri),
		&gorm.Config{},
	)
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&football.Result{},
		&football.SeasonTable{},
		&model.User{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}
