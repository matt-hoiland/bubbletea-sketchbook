package main

import (
	"log/slog"
	"os"

	tea "charm.land/bubbletea/v2"

	"github.com/matt-hoiland/bubbletea-sketchbook/tools/basic-todo/window"
)

func main() {
	if _, err := tea.NewProgram(window.New(
		window.WithTitle("basic-todo"),
	)).Run(); err != nil {
		slog.Error("program failed to start", slog.String("err", err.Error()))
		os.Exit(1)
	}
}
