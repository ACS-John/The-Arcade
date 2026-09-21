package main

import "testing"

func cut(g *game, x int) {
	g.Player = point{x, 0}
	g.Anchor = g.Player
	for y := 1; y < H; y++ {
		g.move(1)
	}
}
func TestCaptureAndScoring(t *testing.T) {
	a := newGame()
	a.Mode = "playing"
	cut(a, 25)
	if a.Percent < 24 || a.Percent > 28 {
		t.Fatalf("capture %v", a.Percent)
	}
	if a.Board[30][12] != 2 || a.Board[30][70] != 0 {
		t.Fatal("wrong region filled")
	}
	if len(a.Trail) != 0 || a.Mark[20][25] {
		t.Fatal("trail not cleared")
	}
	if a.Score != 25*(H-2)*10 {
		t.Fatalf("capture score %d", a.Score)
	}
}
func TestQixProtectsItsRegion(t *testing.T) {
	g := newGame()
	g.Mode = "playing"
	for i := range g.Q {
		g.Q[i].X = 10 + float64(i)
	}
	cut(g, 70)
	if g.Board[30][12] != 0 || g.Board[30][85] != 2 {
		t.Fatal("filled Qix region")
	}
}
func TestTrailDeathAndFuse(t *testing.T) {
	g := newGame()
	g.Mode = "playing"
	g.move(3)
	g.move(3)
	g.move(1)
	if g.Lives != 2 || len(g.Trail) != 0 {
		t.Fatal("self crossing not fatal")
	}
	g.Mode = "playing"
	g.move(3)
	for i := 0; i < 200; i++ {
		g.tick(.01)
	}
	if g.Lives != 1 {
		t.Fatalf("fuse did not catch stopped player: %d", g.Lives)
	}
}
func TestBoundaryAndPause(t *testing.T) {
	g := newGame()
	before := g.Player
	g.Mode = "playing"
	g.move(3)
	if g.Player == before || len(g.Trail) != 1 {
		t.Fatal("movement did not automatically draw")
	}
	g.Mode = "paused"
	q := g.Q
	g.tick(1)
	if g.Q != q {
		t.Fatal("pause advanced Qix")
	}
}
func TestCompletion(t *testing.T) {
	g := newGame()
	g.Mode = "playing"
	for i := range g.Q {
		g.Q[i].X = 8 + float64(i)
	}
	cut(g, 20)
	if g.Mode != "cleared" || g.Percent < 75 {
		t.Fatal("level not completed")
	}
	g.tick(2)
	if g.Level != 2 || g.Percent != 0 || g.Lives != 3 {
		t.Fatal("next level failed")
	}
}
func TestQixLineHitsTrail(t *testing.T) {
	g := newGame()
	g.Mode = "playing"
	p := point{48, 36}
	g.Mark[p.Y][p.X] = true
	g.Q = [4]vec{{40, 36}, {56, 36}, {0, 0}, {0, 0}}
	if !g.hitsTrail() {
		t.Fatal("line interior collision missed")
	}
}
func TestSimulation(t *testing.T) {
	g := newGame()
	g.restart()
	for i := 0; i < 30000; i++ {
		g.Keys = [256]bool{}
		g.Keys[37+i/170%4] = true
		g.tick(1.0 / 120)
		if !inside(g.Player) {
			t.Fatal("player out of bounds")
		}
		if g.Mode == "over" {
			g.restart()
		}
		if !g.qValid(g.Q) {
			t.Fatal("Qix crossed claimed territory")
		}
		g.Events = nil
	}
}

func TestQixComplexityByLevel(t *testing.T) {
	g := newGame()
	for level := 1; level <= 6; level++ {
		g.Level = level
		g.level()
		if g.QCount != min(4, level+1) {
			t.Fatalf("level %d has %d points", level, g.QCount)
		}
		if !g.qValid(g.Q) {
			t.Fatal("invalid starting Qix")
		}
	}
}
func TestArrowOnlyDrawing(t *testing.T) {
	for _, key := range []int{38, 'W'} {
		g := newGame()
		g.restart()
		g.Keys[key] = true
		g.tick(.03)
		if len(g.Trail) != 1 {
			t.Fatalf("key %d did not draw automatically", key)
		}
	}
}
func TestInactiveQixPointsIgnored(t *testing.T) {
	g := newGame()
	g.Q[2] = vec{0, 0}
	g.Q[3] = vec{0, 0}
	if !g.qValid(g.Q) {
		t.Fatal("inactive points affect collision")
	}
	g.Mark[0][0] = true
	if g.hitsTrail() {
		t.Fatal("inactive points hit trail")
	}
}

func TestShiftSlowDraw(t *testing.T) {
	fast := newGame()
	fast.restart()
	fast.Keys[38] = true
	fast.tick(.1)
	slow := newGame()
	slow.restart()
	slow.Keys[38] = true
	slow.Keys[16] = true
	slow.tick(.1)
	if len(slow.Trail) >= len(fast.Trail) || len(slow.Trail) == 0 {
		t.Fatal("Shift did not slow automatic drawing")
	}
	a := newGame()
	a.restart()
	cut(a, 25)
	b := newGame()
	b.restart()
	b.Keys[16] = true
	cut(b, 25)
	if b.Score != 2*a.Score || b.Board[30][12] != 3 {
		t.Fatal("full slow capture did not earn double")
	}
	c := newGame()
	c.restart()
	c.Player = point{25, 0}
	c.Anchor = c.Player
	c.Keys[16] = true
	c.move(1)
	c.Keys[16] = false
	for y := 2; y < H; y++ {
		c.move(1)
	}
	if c.Score != a.Score {
		t.Fatal("mixed-speed capture received slow bonus")
	}
}
