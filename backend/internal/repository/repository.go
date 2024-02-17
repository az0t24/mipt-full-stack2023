package repository

import (
	"gorm.io/gorm"
	"scoreboardpro/internal/entity"
	"scoreboardpro/internal/repository/sqlite"
)

type Repository struct {
	User   entity.UserRepository
	Table  entity.TableRepository
	Result entity.ResultRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		User:   sqlite.NewUserSQLite(db),
		Table:  sqlite.NewTableSQLite(db),
		Result: sqlite.NewResultSQLite(db),
	}
}
