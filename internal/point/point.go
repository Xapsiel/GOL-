package point

type Point struct {
	X         int
	Y         int
	Status    string
	NewStatus string
	LiveCount int
	DeadCount int
}

func NewPoint(x, y int) *Point {
	return &Point{X: x, Y: y, Status: "live"}
}
func (p *Point) Check(n int, area [][]Point) {
	p.LiveCount = 0
	p.DeadCount = 0
	for x := p.X - 1; x <= p.X+1; x++ {
		for y := p.Y - 1; y <= p.Y+1; y++ {
			if x == p.X && y == p.Y {
				continue
			}
			if area[x][y].Status == "border" {
				continue
			} else if area[x][y].Status == "live" {
				p.LiveCount++
			} else if area[x][y].Status == "dead" {
				p.DeadCount++
			}
		}
	}

}
