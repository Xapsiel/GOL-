package game

import (
	"crypton/internal/point"
	"fmt"
	"time"
)

type GOL struct {
	n        int
	Area     [][]point.Point
	GameMap  *string
	finished bool
}

func NewGIL(n int) *GOL {
	area := generateArea(n)
	return &GOL{n, area, new(string), false}
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

func (g *GOL) Run(speed time.Duration) {
	g.Print()
	fmt.Print(g.GameMap)

	for {
		g.Print()
		time.Sleep(speed)
		fmt.Print(g.GameMap)
	}
}
func (g *GOL) Step() {
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

}
func (g *GOL) commit() {
	for x := 1; x < g.n-1; x++ {
		for y := 1; y < g.n-1; y++ {
			g.Area[x][y].Status = g.Area[x][y].NewStatus
		}
	}
}
func (g *GOL) Print() {
	*(g.GameMap) = ""
	for x := 0; x <= g.n-1; x++ {
		for y := 0; y <= g.n-1; y++ {
			if g.Area[x][y].Status == "live" || g.Area[x][y].Status == "busy" {
				g.Area[x][y].Status = "live"
				*g.GameMap += fmt.Sprintf("%c ", rune(9634))
			} else if g.Area[x][y].Status == "dead" {
				*g.GameMap += fmt.Sprintf("%c ", rune(9635))
			} else if g.Area[x][y].Status == "border" {
				*g.GameMap += fmt.Sprintf("%c ", rune(9673))
			} else if g.Area[x][y].Status == "active" {
				*g.GameMap += fmt.Sprintf("%c ", rune(9668))
			}

		}
		*g.GameMap += "\n"
	}

}
