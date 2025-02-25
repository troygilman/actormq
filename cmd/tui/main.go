package main

import (
	"log/slog"
	"os"

	"github.com/anthdm/hollywood/actor"
	"github.com/troygilman/actormq/tui"
)

func main() {
	logFile, err := os.OpenFile("tmp/application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	slog.SetDefault(slog.New(slog.NewTextHandler(logFile, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	if err := tui.Run(tui.Config{
		Pods: []*actor.PID{
			actor.NewPID("127.0.0.1:8080", "pod/A"),
			actor.NewPID("127.0.0.1:8080", "pod/B"),
			actor.NewPID("127.0.0.1:8080", "pod/C"),
		},
	}); err != nil {
		panic(err)
	}
}
