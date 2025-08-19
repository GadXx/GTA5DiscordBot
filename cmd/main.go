package main

import (
	"discordbot/internal/app"
	"log/slog"
)

func main() {
	slog.Info("Starting bot")

	application, err := app.NewApp()
	if err != nil {
		slog.Error("failed to init app", slog.String("error", err.Error()))
		return
	}

	if err := application.Run(); err != nil {
		slog.Error("failed to run bot", slog.String("error", err.Error()))
	}

	<-make(chan struct{})
	application.Close()
}
