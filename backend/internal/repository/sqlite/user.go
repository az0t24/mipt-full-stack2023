package sqlite

import (
	"errors"
	"gorm.io/gorm"
	"scoreboardpro/internal/entity"
	"scoreboardpro/internal/entity/model"
)

type UserSQLite struct {
	db *gorm.DB
}

func NewUserSQLite(db *gorm.DB) *UserSQLite {
	return &UserSQLite{db: db}
}

func (r *UserSQLite) GetAll() (*[]model.User, error) {
	var users []model.User

	if result := r.db.Find(&users); result.Error != nil {
		return nil, result.Error
	} else {
		return &users, nil
	}
}

func (r *UserSQLite) Get(id uint) (*model.User, error) {
	var user model.User

	if result := r.db.Where("id = ?", id).Preload("LastVisit").First(&user); result.Error == nil {
		return &user, nil
	} else if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &user, entity.ErrUserNotFound
	} else {
		return &user, result.Error
	}
}

func (r *UserSQLite) GetByEmail(email string) (*model.User, error) {
	var user model.User

	if result := r.db.Where("email = ?", email).First(&user); result.Error == nil {
		return &user, nil
	} else if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &user, entity.ErrUserNotFound
	} else {
		return &user, result.Error
	}
}

func (r *UserSQLite) Create(user *model.User) error {
	if result := r.db.Create(user); result.Error != nil {
		switch {
		case errors.Is(result.Error, gorm.ErrDuplicatedKey):
			return entity.ErrUserExists
		default:
			return result.Error
		}
	} else {
		return nil
	}
}

func (r *UserSQLite) Update(user *model.User) error {
	result := r.db.Model(user).Updates(user)
	if result.Error != nil {
		return result.Error
	} else {
		return nil
	}
}

func (r *UserSQLite) Delete(id uint) error {
	result := r.db.Delete(&model.User{}, id)
	if result.Error != nil {
		return result.Error
	} else {
		return nil
	}
}
