package football

import "gorm.io/gorm"

type Result struct {
	gorm.Model
	SeasonTableID uint

	TeamName     string
	Wins         uint32
	Draws        uint32
	Loses        uint32
	GoalsFor     uint32
	GoalsAgainst uint32
	GoalsDiff    uint32
	Points       uint32
}
