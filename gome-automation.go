package main

import (
	"log/slog"
	"os"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/lupuscrzy/gome-automation/internal/dbrepo"
)

var conf = koanf.New(".")

func main() {
	loggerSetup()
	slog.Info("starting gome-automation")

	if err := conf.Load(file.Provider("internal/config/config.yml"), yaml.Parser()); err != nil {
		slog.Info("error loading config", "err", err)
	}
	envPrefix := "GOME_"
	err := conf.Load(env.Provider(envPrefix, ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, envPrefix)), "_", ".", -1)
	}), nil)
	if err != nil {
		slog.Error("error loading config", "err", err)
	}
	slog.Info("Database Filename", "name", conf.String("DatabaseName"))
	slog.Info("Minutes To Pull", "min", conf.Int("Minutes"))

	db := dbrepo.Setup()
	// _ := dbrepo.InsertTemp(db, dbrepo.Thermostat, 69.420)
	readings, _ := dbrepo.ReadRecentTemps(db, 1)
	for _, reading := range readings {
		println(reading.Temperature)
	}
	slog.Info("stopping gome-automation")
}

func loggerSetup() {
	var logLevel = new(slog.LevelVar)
	logLevel.Set(slog.LevelInfo)
	textHandler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: logLevel})
	logger := slog.New(textHandler)
	slog.SetDefault(logger)
}
