package main

import (
	"crypton/internal/app"
)

func main() {
	app.Run(20)

	// g := game.NewGIL(20)
	// g.Area[3][3].Status="live"
	// g.Area[4][3].Status="live"
	// g.Area[5][3].Status="live"
	// g.Step()
	// fmt.Println(g.Print())
	// g.Step()
	// fmt.Println(g.Print())
	// points := []point.Point{*point.NewPoint(2, 2), *point.NewPoint(3, 2), *point.NewPoint(4, 2)}
	// gil := game.NewGIL(20, points...)
	// gil.Run(1000 * time.Millisecond)
}
