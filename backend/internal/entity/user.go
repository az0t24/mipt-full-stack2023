package entity

import (
	"scoreboardpro/internal/entity/model"
	"scoreboardpro/internal/entity/model/football"
)

type UserRepository interface {
	Create(*model.User) error
	GetAll() (*[]model.User, error)
	Get(id uint) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	Update(*model.User) error
	Delete(id uint) error
}

type UserService interface {
	Get(id uint) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	GetAll() (*[]model.User, error)
	Update(user *model.User) error
	Delete(user *model.User) error

	Register(userReg *model.UserRegister) error
	Login(userLogin *model.UserLogin) (uint, error)

	MarkAsFavorite(id, tableId uint) error
	GetFavorites(id uint) ([]football.SeasonTable, error)
}
