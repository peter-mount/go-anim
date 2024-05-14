package tween

type Point struct {
	X, Y float64
}

func (p Point) Translate(dx, dy float64) Point {
	return Point{p.X + dx, p.Y + dy}
}

func (p Point) Add(q Point) Point {
	return Point{p.X + q.X, p.Y + q.Y}
}

func (p Point) Sub(q Point) Point {
	return Point{p.X - q.X, p.Y - q.Y}
}

type Line struct {
	Start, End Point
	delta      Point
	d          float64
}

func (l Line) Translate(frameLength int, t float64) Point {
	if l.d == 0 {
		l.d = t / float64(frameLength)
		l.delta = l.End.Sub(l.Start)
	}
	return l.Start.Translate(l.delta.X*l.d, l.delta.Y*l.d)
}
