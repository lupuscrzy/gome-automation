package dbrepo

import "gorm.io/gorm"

/*
database table models
new models are registered in db.go for migrations
*/

type Temperature struct {
	gorm.Model
	Sensor      Sensor  `gorm:"not null"`
	Temperature float32 `gorm:"not null"`
}
