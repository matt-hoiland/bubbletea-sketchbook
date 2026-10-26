package main

import (
	"log"
	"os"

	tea "charm.land/bubbletea/v2"
)

func main() {
	m := New(
		WithWidth(30),
	)
	if _, ok := os.LookupEnv("DEBUG"); ok {
		var err error
		if m.dump, err = os.Create("message.log"); err != nil {
			log.Panic(err)
		}
	}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		log.Panic(err)
	}
}
