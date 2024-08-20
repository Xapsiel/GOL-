package main

import (
	"crypton/internal/app"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(app.InitialModel(4))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
	// points := []point.Point{*point.NewPoint(2, 2), *point.NewPoint(3, 2), *point.NewPoint(4, 2)}
	// gil := game.NewGIL(20, points...)
	// gil.Run(1000 * time.Millisecond)
}
