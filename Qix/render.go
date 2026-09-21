package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

const SW = 1280
const SH = 864
const BX = 64
const BY = 164
const Cell = 8

type canvas struct{ Pix, Backdrop []uint32 }

func newCanvas(bg image.Image) *canvas {
	c := &canvas{Pix: make([]uint32, SW*SH), Backdrop: make([]uint32, SW*SH)}
	b := bg.Bounds()
	for y := 0; y < SH; y++ {
		for x := 0; x < SW; x++ {
			r, g, v, _ := bg.At(b.Min.X+x*b.Dx()/SW, b.Min.Y+y*b.Dy()/SH).RGBA()
			c.Backdrop[y*SW+x] = uint32(r>>8)<<16 | uint32(g>>8)<<8 | uint32(v>>8)
		}
	}
	return c
}
func (c *canvas) rect(x, y, w, h int, v uint32) {
	for yy := max(0, y); yy < min(SH, y+h); yy++ {
		for xx := max(0, x); xx < min(SW, x+w); xx++ {
			c.Pix[yy*SW+xx] = v
		}
	}
}
func (c *canvas) blend(x, y int, v uint32, a float64) {
	if x < 0 || x >= SW || y < 0 || y >= SH {
		return
	}
	a = math.Max(0, math.Min(1, a))
	o := c.Pix[y*SW+x]
	r := float64(o>>16&255)*(1-a) + float64(v>>16&255)*a
	g := float64(o>>8&255)*(1-a) + float64(v>>8&255)*a
	b := float64(o&255)*(1-a) + float64(v&255)*a
	c.Pix[y*SW+x] = uint32(r)<<16 | uint32(g)<<8 | uint32(b)
}
func (c *canvas) line(x1, y1, x2, y2 float64, v uint32, w, a float64) {
	dx, dy := x2-x1, y2-y1
	l := dx*dx + dy*dy
	for y := max(0, int(math.Min(y1, y2)-w-1)); y <= min(SH-1, int(math.Max(y1, y2)+w+1)); y++ {
		for x := max(0, int(math.Min(x1, x2)-w-1)); x <= min(SW-1, int(math.Max(x1, x2)+w+1)); x++ {
			t := 0.0
			if l > 0 {
				t = math.Max(0, math.Min(1, ((float64(x)-x1)*dx+(float64(y)-y1)*dy)/l))
			}
			d := math.Hypot(float64(x)-x1-t*dx, float64(y)-y1-t*dy)
			if d < w {
				c.blend(x, y, v, a*(1-d/w))
			}
		}
	}
}
func (c *canvas) glow(p vec, v uint32, r float64) {
	x, y := float64(BX)+p.X*Cell, float64(BY)+p.Y*Cell
	c.line(x, y, x, y, v, r, .55)
	c.line(x, y, x, y, 0xfaffff, 3, 1)
}
func (c *canvas) draw(g *game) {
	copy(c.Pix, c.Backdrop)
	c.rect(BX-5, BY-5, (W-1)*Cell+11, (H-1)*Cell+11, 0x345a79)
	c.rect(BX-3, BY-3, (W-1)*Cell+7, (H-1)*Cell+7, 0x070d20)
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			state := g.Board[y][x]
			v := uint32(0x070d20)
			if state == 2 {
				v = 0x123e55
			}
			if state == 3 {
				v = 0x4a2858
			}
			if state == 1 {
				v = 0x24536c
			}
			if state > 0 {
				c.rect(BX+x*Cell-3, BY+y*Cell-3, Cell, Cell, v)
			}
			if state == 0 && x%4 == 0 && y%4 == 0 {
				c.blend(BX+x*Cell, BY+y*Cell, 0x43658b, .35)
			}
			if state > 0 && g.boundary(point{x, y}) {
				c.rect(BX+x*Cell-1, BY+y*Cell-1, 3, 3, 0x4a94ad)
			}
		}
	}
	for h, q := range g.History {
		a := float64(h+1) / float64(max(1, len(g.History)))
		for i := 0; i < g.QCount-1; i++ {
			v := uint32(0xac5cff)
			if (h+i)%3 == 0 {
				v = 0x43eeea
			}
			x1, y1 := float64(BX)+q[i].X*Cell, float64(BY)+q[i].Y*Cell
			x2, y2 := float64(BX)+q[i+1].X*Cell, float64(BY)+q[i+1].Y*Cell
			c.line(x1, y1, x2, y2, v, 5, .10*a)
			c.line(x1, y1, x2, y2, v, 1.5, .7*a)
		}
	}
	for i := 0; i < g.QCount-1; i++ {
		a, b := g.Q[i], g.Q[i+1]
		c.line(float64(BX)+a.X*Cell, float64(BY)+a.Y*Cell, float64(BX)+b.X*Cell, float64(BY)+b.Y*Cell, 0xf8ddff, 1.4, 1)
	}
	prior := g.Anchor
	trailColor := uint32(0x41e8f0)
	if g.Slow {
		trailColor = 0xff7dbb
	}
	for _, p := range g.Trail {
		c.line(float64(BX+prior.X*Cell), float64(BY+prior.Y*Cell), float64(BX+p.X*Cell), float64(BY+p.Y*Cell), trailColor, 4, .8)
		prior = p
	}
	if g.FuseOn && len(g.Trail) > 0 {
		p := g.Trail[min(len(g.Trail)-1, int(g.Fuse))]
		c.glow(vec{float64(p.X), float64(p.Y)}, 0xff6b35, 13)
	}
	for _, s := range g.Sparks {
		p := vec{float64(s.P.X), float64(s.P.Y)}
		c.glow(p, 0xffb454, 14)
		x, y := float64(BX)+p.X*Cell, float64(BY)+p.Y*Cell
		r := 8.0
		angle := g.Clock * 6
		c.line(x-math.Cos(angle)*r, y-math.Sin(angle)*r, x+math.Cos(angle)*r, y+math.Sin(angle)*r, 0xffdb7b, 2, 1)
	}
	for _, m := range g.Particles {
		c.glow(m.P, m.Color, 4*m.Life/m.Max)
	}
	p := g.Player
	c.glow(vec{float64(p.X), float64(p.Y)}, 0x67fcff, 14)
	x, y := float64(BX+p.X*Cell), float64(BY+p.Y*Cell)
	for i := 0; i < 4; i++ {
		a := float64(i) * math.Pi / 2
		b := a + math.Pi/2
		c.line(x+math.Cos(a)*6, y+math.Sin(a)*6, x+math.Cos(b)*6, y+math.Sin(b)*6, 0xffffff, 1.4, 1)
	}
	c.rect(880, 164, 332, 576, 0x0b142b)
	c.rect(880, 164, 3, 576, 0x477692)
	c.rect(906, 369, 278, 8, 0x22354d)
	c.rect(906, 369, int(278*g.Percent/100), 8, 0x52dedf)
	c.rect(906+208, 365, 2, 16, 0xffcd91)
	if g.Mode != "playing" && g.Mode != "hurt" {
		c.rect(201, 372, 486, 151, 0x15203b)
		c.rect(201, 372, 486, 2, 0x7ab6d1)
	}
}
func (c *canvas) save(path string) error {
	im := image.NewRGBA(image.Rect(0, 0, SW, SH))
	for i, v := range c.Pix {
		im.SetRGBA(i%SW, i/SW, color.RGBA{uint8(v >> 16), uint8(v >> 8), uint8(v), 255})
	}
	f, e := os.Create(path)
	if e != nil {
		return e
	}
	defer f.Close()
	return png.Encode(f, im)
}
