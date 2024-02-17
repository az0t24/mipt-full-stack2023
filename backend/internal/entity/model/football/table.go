package football

import "gorm.io/gorm"

type SeasonTable struct {
	gorm.Model
	Season           string
	ChampionshipName string
	Teams            []Result `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
