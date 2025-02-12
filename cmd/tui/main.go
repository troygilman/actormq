package main

import (
	"log/slog"
	"os"

	"github.com/troygilman/actormq/tui"
)

func main() {
	logFile, err := os.OpenFile("tmp/application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	slog.SetDefault(slog.New(slog.NewTextHandler(logFile, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	if err := tui.Run(); err != nil {
		panic(err)
	}
}
