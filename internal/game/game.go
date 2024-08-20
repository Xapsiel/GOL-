package game

import (
	"crypton/internal/point"
	"fmt"
	"time"
)

type GIL struct {
	n        int
	Area     [][]point.Point
	finished bool
}

func NewGIL(n int) *GIL {
	area := generateArea(n)
	return &GIL{n, area, false}
}

func generateArea(n int) [][]point.Point {
	area := make([][]point.Point, n)
	for x := 0; x <= n-1; x++ {

		area[x] = make([]point.Point, n)
		for y := 0; y <= n-1; y++ {
			if x == 0 || x == n-1 || y == 0 || y == n-1 {
				area[x][y] = point.Point{X: x, Y: y, Status: "border"}

			} else {
				area[x][y] = point.Point{X: x, Y: y, Status: "dead"}

			}
		}
	}
	return area
}

func (g *GIL) Run(speed time.Duration) {
	fmt.Print(g.Print())

	for {
		for x := 1; x < g.n-1; x++ {
			for y := 1; y < g.n-1; y++ {
				g.Area[x][y].Check(g.n, g.Area)
				g.Area[x][y].NewStatus = g.Area[x][y].Status
				if g.Area[x][y].Status == "dead" && g.Area[x][y].LiveCount == 3 {
					g.Area[x][y].NewStatus = "live"
				} else if g.Area[x][y].Status == "live" && !(g.Area[x][y].LiveCount == 2 || g.Area[x][y].LiveCount == 3) {
					g.Area[x][y].NewStatus = "dead"
				}

			}
		}
		g.commit()
		time.Sleep(speed)
		fmt.Print(g.Print())
	}
}
func (g *GIL) commit() {
	for x := 1; x < g.n-1; x++ {
		for y := 1; y < g.n-1; y++ {
			g.Area[x][y].Status = g.Area[x][y].NewStatus
		}
	}
}
func (g *GIL) Print() string {
	gameMap := ""
	for x := 0; x <= g.n-1; x++ {
		for y := 0; y <= g.n-1; y++ {
			if g.Area[x][y].Status == "live" {
				gameMap += fmt.Sprintf("%c ", rune(9634))
			} else if g.Area[x][y].Status == "dead" {
				gameMap += fmt.Sprintf("%c ", rune(9635))
			} else if g.Area[x][y].Status == "border" {
				gameMap += fmt.Sprintf("%c ", rune(9673))
			} else if g.Area[x][y].Status == "active" {
				gameMap += fmt.Sprintf("%c ", rune(9668))
			} else {
				gameMap += fmt.Sprintf("%c ", rune(9671))

			}
		}
		gameMap += "\n"
	}
	return fmt.Sprintf("%s", string(gameMap))
}
