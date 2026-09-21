package main

import (
	"math"
	"math/rand"
)

const W = 96
const H = 72

type point struct{ X, Y int }
type vec struct{ X, Y float64 }
type spark struct {
	P     point
	Dir   int
	Clock float64
	Turn  int
}
type mote struct {
	P, V      vec
	Life, Max float64
	Color     uint32
}
type game struct {
	Board                                                   [H][W]uint8
	Mark                                                    [H][W]bool
	Player, Anchor                                          point
	Trail                                                   []point
	Q                                                       [4]vec
	Velocity                                                [4]vec
	History                                                 [][4]vec
	Sparks                                                  []spark
	Particles                                               []mote
	Rng                                                     *rand.Rand
	Score, Best, Lives, Level                               int
	QCount                                                  int
	Percent                                                 float64
	Mode                                                    string
	Slow, Muted                                             bool
	Clock, MoveClock, Idle, Fuse, LevelTime, Flash, Claimed float64
	FuseOn                                                  bool
	Keys                                                    [256]bool
	Events                                                  []string
}

var dirs = []point{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}

func newGame() *game {
	g := &game{Rng: rand.New(rand.NewSource(42)), Mode: "ready", Lives: 3, Level: 1}
	g.level()
	return g
}
func inside(p point) bool { return p.X >= 0 && p.X < W && p.Y >= 0 && p.Y < H }
func (g *game) level() {
	g.QCount = min(4, g.Level+1)
	g.Board = [H][W]uint8{}
	g.Mark = [H][W]bool{}
	g.Trail = nil
	g.FuseOn = false
	g.Idle = 0
	g.Fuse = 0
	g.Percent = 0
	g.History = nil
	g.LevelTime = 0
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			if x == 0 || x == W-1 || y == 0 || y == H-1 {
				g.Board[y][x] = 1
			}
		}
	}
	g.Player = point{W / 2, H - 1}
	g.Anchor = g.Player
	for i := 0; i < g.QCount; i++ {
		g.Q[i] = vec{float64(W)/2 + float64(i*2-3), float64(H)/2 + float64(i%2*6-3)}
		g.Velocity[i] = vec{float64(9 + i*2), float64(8 - i*3)}
	}
	g.Sparks = []spark{{point{0, 0}, 0, 0, 1}, {point{W - 1, 0}, 2, 0, -1}}
}
func (g *game) restart() {
	g.Score = 0
	g.Lives = 3
	g.Level = 1
	g.Mode = "playing"
	g.level()
	g.event("start")
}
func (g *game) event(s string) {
	if len(g.Events) < 16 {
		g.Events = append(g.Events, s)
	}
}
func (g *game) boundary(p point) bool {
	if !inside(p) || g.Board[p.Y][p.X] == 0 {
		return false
	}
	if p.X == 0 || p.Y == 0 || p.X == W-1 || p.Y == H-1 {
		return true
	}
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if g.Board[p.Y+dy][p.X+dx] == 0 {
				return true
			}
		}
	}
	return false
}
func (g *game) nearest(p point) point {
	for radius := 0; radius < W+H; radius++ {
		for y := max(0, p.Y-radius); y <= min(H-1, p.Y+radius); y++ {
			for x := max(0, p.X-radius); x <= min(W-1, p.X+radius); x++ {
				q := point{x, y}
				if abs(x-p.X)+abs(y-p.Y) == radius && g.boundary(q) {
					return q
				}
			}
		}
	}
	return point{0, 0}
}
func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
func (g *game) direction() int {
	if g.Keys[37] || g.Keys['A'] {
		return 2
	}
	if g.Keys[39] || g.Keys['D'] {
		return 0
	}
	if g.Keys[38] || g.Keys['W'] {
		return 3
	}
	if g.Keys[40] || g.Keys['S'] {
		return 1
	}
	return -1
}
func (g *game) burst(p point, c uint32) {
	for i := 0; i < 32; i++ {
		a := g.Rng.Float64() * 2 * math.Pi
		s := g.Rng.Float64()*14 + 3
		g.Particles = append(g.Particles, mote{vec{float64(p.X), float64(p.Y)}, vec{math.Cos(a) * s, math.Sin(a) * s}, .65, .65, c})
	}
}
func (g *game) die() {
	if g.Mode != "playing" {
		return
	}
	for _, p := range g.Trail {
		g.Mark[p.Y][p.X] = false
	}
	g.Trail = nil
	g.FuseOn = false
	g.Fuse = 0
	g.Idle = 0
	g.Lives--
	g.burst(g.Player, 0xff658b)
	g.Flash = .35
	g.event("lose")
	g.Player = g.nearest(g.Anchor)
	g.Sparks = []spark{{point{0, 0}, 0, 0, 1}, {point{W - 1, 0}, 2, 0, -1}}
	if g.Lives <= 0 {
		g.Mode = "over"
	} else {
		g.Mode = "hurt"
		g.Claimed = 1.0
	}
}
func (g *game) move(d int) {
	if d < 0 {
		return
	}
	n := point{g.Player.X + dirs[d].X, g.Player.Y + dirs[d].Y}
	if !inside(n) {
		return
	}
	if len(g.Trail) == 0 {
		if g.boundary(n) {
			g.Player = n
			g.Anchor = n
			return
		}
		if g.Board[n.Y][n.X] != 0 {
			return
		}
		g.Anchor = g.Player
		g.Slow = g.Keys[16]
		g.FuseOn = false
		g.Fuse = 0
	} else {
		if !g.Keys[16] {
			g.Slow = false
		}
		if g.Mark[n.Y][n.X] {
			g.die()
			return
		}
	}
	g.Player = n
	g.Idle = 0
	if g.Board[n.Y][n.X] != 0 {
		g.capture()
		return
	}
	g.Trail = append(g.Trail, n)
	g.Mark[n.Y][n.X] = true
}
func sampleLine(a, b vec, fn func(point) bool) bool {
	n := int(math.Max(math.Abs(b.X-a.X), math.Abs(b.Y-a.Y))*4) + 1
	for i := 0; i <= n; i++ {
		t := float64(i) / float64(n)
		p := point{int(math.Round(a.X + (b.X-a.X)*t)), int(math.Round(a.Y + (b.Y-a.Y)*t))}
		if !fn(p) {
			return false
		}
	}
	return true
}
func (g *game) qValid(q [4]vec) bool {
	for i := 0; i < g.QCount-1; i++ {
		if !sampleLine(q[i], q[i+1], func(p point) bool { return inside(p) && g.Board[p.Y][p.X] == 0 }) {
			return false
		}
	}
	return true
}
func (g *game) hitsTrail() bool {
	hit := false
	for i := 0; i < g.QCount-1; i++ {
		sampleLine(g.Q[i], g.Q[i+1], func(p point) bool {
			if inside(p) && g.Mark[p.Y][p.X] {
				hit = true
			}
			return true
		})
	}
	return hit
}
func (g *game) capture() {
	if g.hitsTrail() {
		g.die()
		return
	}
	tint := uint8(2)
	mult := 10
	if g.Slow {
		tint = 3
		mult = 20
	}
	fresh := 0
	for _, p := range g.Trail {
		g.Board[p.Y][p.X] = tint
		g.Mark[p.Y][p.X] = false
		fresh++
	}
	g.Trail = nil
	g.FuseOn = false
	g.Idle = 0
	g.Fuse = 0
	// Keep every region touched by the live Qix. Fill every other connected region.
	seen := [H][W]bool{}
	queue := []point{}
	seed := func(p point) bool {
		if inside(p) && g.Board[p.Y][p.X] == 0 && !seen[p.Y][p.X] {
			seen[p.Y][p.X] = true
			queue = append(queue, p)
		}
		return true
	}
	for i := 0; i < g.QCount-1; i++ {
		sampleLine(g.Q[i], g.Q[i+1], seed)
	}
	for head := 0; head < len(queue); head++ {
		p := queue[head]
		for _, d := range dirs {
			seed(point{p.X + d.X, p.Y + d.Y})
		}
	}
	count := 0
	for y := 1; y < H-1; y++ {
		for x := 1; x < W-1; x++ {
			if g.Board[y][x] == 0 && !seen[y][x] {
				g.Board[y][x] = tint
				fresh++
			}
			if g.Board[y][x] != 0 {
				count++
			}
		}
	}
	g.Score += fresh * mult
	g.Best = max(g.Best, g.Score)
	g.Percent = float64(count) * 100 / float64((W-2)*(H-2))
	g.Flash = .25
	g.burst(g.Player, 0x58e8f2)
	g.event("clear")
	g.Player = g.nearest(g.Player)
	g.Anchor = g.Player
	for i := range g.Sparks {
		g.Sparks[i].P = g.nearest(g.Sparks[i].P)
	}
	if g.Percent >= 75 {
		g.Score += g.Level*1000 + int(g.Percent-75)*100
		g.Best = max(g.Best, g.Score)
		g.Mode = "cleared"
		g.Claimed = 1.8
		g.event("win")
	}
}
func (g *game) tick(dt float64) {
	g.Clock += dt
	g.Flash = math.Max(0, g.Flash-dt)
	alive := g.Particles[:0]
	for _, m := range g.Particles {
		m.Life -= dt
		m.P.X += m.V.X * dt
		m.P.Y += m.V.Y * dt
		if m.Life > 0 {
			alive = append(alive, m)
		}
	}
	g.Particles = alive
	if g.Mode == "hurt" || g.Mode == "cleared" {
		g.Claimed -= dt
		if g.Claimed <= 0 {
			if g.Mode == "cleared" {
				g.Level++
				g.level()
			}
			g.Mode = "playing"
		}
		return
	}
	if g.Mode != "playing" {
		return
	}
	g.LevelTime += dt
	g.MoveClock += dt

	speed := .021
	if g.Keys[16] {
		speed = .048
	}
	for g.MoveClock >= speed {
		g.MoveClock -= speed
		g.move(g.direction())
		if g.Mode != "playing" {
			return
		}
	}
	if len(g.Trail) > 0 {
		g.Idle += dt
		if g.Idle > .6 {
			g.FuseOn = true
		}
		if g.FuseOn {
			g.Fuse += dt * (17 + float64(g.Level)*2)
			if int(g.Fuse) >= len(g.Trail) {
				g.die()
				return
			}
		}
	}
	qspeed := 1 + float64(g.Level-1)*.1
	for i := 0; i < g.QCount; i++ {
		for axis := 0; axis < 2; axis++ {
			candidate := g.Q
			if axis == 0 {
				candidate[i].X += g.Velocity[i].X * dt * qspeed
			} else {
				candidate[i].Y += g.Velocity[i].Y * dt * qspeed
			}
			if g.qValid(candidate) {
				g.Q = candidate
			} else if axis == 0 {
				g.Velocity[i].X = -g.Velocity[i].X
			} else {
				g.Velocity[i].Y = -g.Velocity[i].Y
			}
		}
		if g.Rng.Float64() < dt*.7 {
			g.Velocity[i].X += (g.Rng.Float64() - .5) * 5
			g.Velocity[i].Y += (g.Rng.Float64() - .5) * 5
		}
	}
	if g.hitsTrail() {
		g.die()
		return
	}
	if len(g.History) == 0 || int(g.Clock*45) != int((g.Clock-dt)*45) {
		g.History = append(g.History, g.Q)
		if len(g.History) > 22 {
			g.History = g.History[1:]
		}
	}
	if g.LevelTime > 30 && len(g.Sparks) == 2 {
		g.Sparks = append(g.Sparks, spark{point{W - 1, H - 1}, 2, 0, 1}, spark{point{0, H - 1}, 0, 0, -1})
		g.event("power")
	}
	for i := range g.Sparks {
		s := &g.Sparks[i]
		s.Clock += dt
		interval := math.Max(.022, .085-float64(g.Level)*.005)
		for s.Clock >= interval {
			s.Clock -= interval
			order := []int{(s.Dir + 4 + s.Turn) % 4, s.Dir, (s.Dir + 4 - s.Turn) % 4, (s.Dir + 2) % 4}
			for _, d := range order {
				n := point{s.P.X + dirs[d].X, s.P.Y + dirs[d].Y}
				if g.boundary(n) {
					s.P = n
					s.Dir = d
					break
				}
			}
			if s.P == g.Player {
				g.die()
				return
			}
		}
	}
}
