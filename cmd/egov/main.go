package main

import (
	"log/slog"
	"os"
	"os/signal"

	"github.com/IceMAN2377/kaspitest/app"
	"github.com/IceMAN2377/kaspitest/internal/config"
)

func main() {
	config := config.NewConfig()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	app := app.NewApp(config, logger)

	go func() {
		logger.Info("Starting application", slog.Int("port", config.HttpPort))
		if err := app.Run(); err != nil {
			panic("failed to start http server: " + err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit
	slog.Info("Stopping application")
}
