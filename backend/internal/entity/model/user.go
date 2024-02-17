package model

import (
	"gorm.io/gorm"
	"scoreboardpro/internal/entity/model/football"
)

type User struct {
	gorm.Model

	UserRegister

	FavoriteTables []football.SeasonTable `gorm:"many2many:user_seasontables;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type UserLogin struct {
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null" json:"Password,omitempty"`
}

type UserRegister struct {
	UserLogin

	FirstName string
	LastName  string
}

func (u *User) OmitPassword() {
	u.Password = ""
}
