package app

import (
	"log/slog"
)

type (
	App struct {
		log *slog.Logger

		// dependencies...
	}
)

func New(log *slog.Logger) *App {
	return &App{
		log: log,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		a.log.Error("failed to run", err)
		panic(err.Error())
	}
}

func (a *App) Run() error {
	a.log.Info("Hello World")

	return nil
}
