package dbrepo

/*
functions to interact with a sqlite database
docs: https://gorm.io/docs/
*/

import (
	"log/slog"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Setup() *gorm.DB {
	slog.Info("initializing database")
	db, err := gorm.Open(sqlite.Open("gome.db"), &gorm.Config{}) // pull "gome.db" from config
	if err != nil {
		msg := "failed to connect to database"
		slog.Error(msg)
		panic(msg)
	}
	db.AutoMigrate(&Temperature{}) // add new tables here, separated by commas
	slog.Info("connected to database")
	return db
}

func InsertTemp(db *gorm.DB, sensor Sensor, temp float32) error {
	slog.Debug("inserting into temperature table", "sensor", sensor, "temp", temp)
	err := db.Create(&Temperature{
		Sensor:      sensor,
		Temperature: temp,
	}).Error
	if err != nil {
		slog.Error("error inserting row in temperature table", "sensor", sensor, "temp", temp, "err", err)
	} else {
		slog.Debug("inserted into temperature table", "sensor", sensor, "temp", temp)
	}
	return err
}

func ReadRecentTemps(db *gorm.DB, minutesToInclude int) ([]Temperature, error) {
	slog.Debug("reading recent temperatures", "minutes", minutesToInclude)
	var results []Temperature
	startTime := time.Now().Add(-1 * time.Duration(minutesToInclude) * time.Minute)
	err := db.Where("created_at > ?", startTime).Find(&results).Error
	if err != nil {
		slog.Error("error querying recent temperatures", "minutes", minutesToInclude, "err", err)
	}
	slog.Debug("returning recent temperatures", "minutes", minutesToInclude, "records", len(results))
	return results, err
}
