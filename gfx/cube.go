package main

import (
	"log"
	"math"
	"runtime"

	"github.com/qeedquan/go-media/sdl"
	"github.com/qeedquan/go-media/sdl/sdlimage/sdlcolor"
)

const (
	W = 1280
	H = 800
)

type Display struct {
	*sdl.Window
	*sdl.Renderer
}

func newDisplay(w, h int, wflag sdl.WindowFlags) (*Display, error) {
	window, renderer, err := sdl.CreateWindowAndRenderer(w, h, wflag)
	return &Display{window, renderer}, err
}

var (
	camera Vec3
)

func main() {
	runtime.LockOSThread()
	log.SetFlags(0)
	err := sdl.Init(sdl.INIT_EVERYTHING)
	ck(err)

	screen, err := newDisplay(W, H, sdl.WINDOW_RESIZABLE)
	ck(err)

	screen.SetLogicalSize(W, H)
	screen.SetTitle("Cube")

	camera = Vec3{4, 4, 4}

	for {
		for {
			ev := sdl.PollEvent()
			if ev == nil {
				break
			}
			switch ev := ev.(type) {
			case sdl.QuitEvent:
				return
			case sdl.KeyDownEvent:
				switch ev.Sym {
				case sdl.K_ESCAPE:
					return
				case sdl.K_LEFT:
					camera.X -= 0.5
				case sdl.K_RIGHT:
					camera.X += 0.5
				case sdl.K_UP:
					camera.Y += 0.5
				case sdl.K_DOWN:
					camera.Y -= 0.5
				case sdl.K_1:
					camera.Z += 0.5
				case sdl.K_2:
					camera.Z -= 0.5
				}
			}
		}

		screen.SetDrawColor(sdlcolor.Black)
		screen.Clear()

		axis := translate(Vec3{-3, 0, 1})

		r0 := Vec4{0, 0, 0, 1}
		r1 := Vec4{1, 0, 0, 1}
		r0 = axis.Mulv(r0)
		r1 = axis.Mulv(r1)
		line3(screen, r0, r1, sdlcolor.Red)

		r0 = Vec4{0, 0, 0, 1}
		r1 = Vec4{0, 1, 0, 1}
		r0 = axis.Mulv(r0)
		r1 = axis.Mulv(r1)
		line3(screen, r0, r1, sdlcolor.Green)

		r0 = Vec4{0, 0, 0, 1}
		r1 = Vec4{0, 0, 1, 1}
		r0 = axis.Mulv(r0)
		r1 = axis.Mulv(r1)
		line3(screen, r0, r1, sdlcolor.Blue)
		r0 = Vec4{0, 1, 0, 1}
		r1 = Vec4{0, 1, 1, 1}
		line3(screen, r0, r1, sdl.Color{40, 40, 61, 255})

		r0 = Vec4{0, 1, 1, 1}
		r1 = Vec4{0, 0, 1, 1}
		line3(screen, r0, r1, sdlcolor.Yellow)

		r0 = Vec4{1, 0, 0, 1}
		r1 = Vec4{1, 0, 1, 1}
		line3(screen, r0, r1, sdlcolor.White)

		r0 = Vec4{0, 0, 1, 1}
		r1 = Vec4{1, 0, 1, 1}
		line3(screen, r0, r1, sdl.Color{100, 56, 40, 255})

		r0 = Vec4{1, 0, 1, 1}
		r1 = Vec4{1, 1, 1, 1}
		line3(screen, r0, r1, sdl.Color{205, 14, 15, 255})

		r0 = Vec4{1, 0, 0, 1}
		r1 = Vec4{1, 1, 0, 1}
		line3(screen, r0, r1, sdl.Color{100, 100, 100, 255})

		r0 = Vec4{0, 1, 1, 1}
		r1 = Vec4{1, 1, 1, 1}
		line3(screen, r0, r1, sdl.Color{100, 100, 100, 255})

		r0 = Vec4{0, 1, 0, 1}
		r1 = Vec4{1, 1, 0, 1}
		line3(screen, r0, r1, sdl.Color{100, 100, 100, 255})

		r0 = Vec4{1, 1, 0, 1}
		r1 = Vec4{1, 1, 1, 1}
		line3(screen, r0, r1, sdl.Color{100, 100, 100, 255})

		r0 = Vec4{0, 0, 0, 1}
		r1 = Vec4{0, 0, 1, 1}
		line3(screen, r0, r1, sdl.Color{100, 100, 100, 255})

		r0 = Vec4{0, 0, 0, 1}
		r1 = Vec4{1, 0, 0, 1}
		line3(screen, r0, r1, sdl.Color{100, 100, 100, 255})

		r0 = Vec4{0, 0, 0, 1}
		r1 = Vec4{1, 0, 0, 1}
		line3(screen, r0, r1, sdl.Color{100, 100, 100, 255})

		r0 = Vec4{0, 0, 0, 1}
		r1 = Vec4{0, 1, 0, 1}
		line3(screen, r0, r1, sdl.Color{100, 100, 100, 255})

		screen.Present()
	}
}

func ck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func line3(re *Display, a, b Vec4, c sdl.Color) {
	modelview := lookat(camera, Vec3{0, 0, 0}, Vec3{0, 1, 0})
	projection := perspective(math.Pi/4, W*1.0/H, 0.1, 1000)

	p := modelview.Mulv(a)
	q := modelview.Mulv(b)

	p = projection.Mulv(p)
	q = projection.Mulv(q)

	p = p.Scale(1 / p.W)
	q = q.Scale(1 / q.W)

	dc := screenSpace(W, H)
	p = dc.Mulv(p)
	q = dc.Mulv(q)

	re.SetDrawColor(c)
	re.DrawLine(int(p.X), int(p.Y), int(q.X), int(q.Y))
}

func line(re *Display, x0, y0, x1, y1 int, c sdl.Color) {
	re.SetDrawColor(c)

	if x0 == x1 {
		vline(re, x0, y0, y1)
		return
	}
	if y0 == y1 {
		hline(re, y0, x0, x1)
		return
	}

	dx := abs(x1 - x0)
	dy := abs(y1 - y0)
	sx, sy := -1, -1
	if x0 < x1 {
		sx = 1
	}
	if y0 < y1 {
		sy = 1
	}

	err := -dy / 2
	if dx > dy {
		err = dx / 2
	}

	for {
		re.DrawPoint(x0, y0)
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := err
		if e2 > -dx {
			err -= dy
			x0 += sx
		}
		if e2 < dy {
			err += dx
			y0 += sy
		}
	}
}

func vline(re *Display, x, y0, y1 int) {
	if y1 < y0 {
		y0, y1 = y1, y0
	}

	for y0 < y1 {
		re.DrawPoint(x, y0)
		y0++
	}
}

func hline(re *Display, y, x0, x1 int) {
	if x1 < x0 {
		x0, x1 = x1, x0
	}

	for x0 < x1 {
		re.DrawPoint(x0, y)
		x0++
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

type Vec4 struct {
	X, Y, Z, W float64
}

func (v Vec4) Add(u Vec4) Vec4 {
	return Vec4{v.X + u.X, v.Y + u.Y, v.Z + u.Z, 1}
}

func (v Vec4) Sub(u Vec4) Vec4 {
	return Vec4{v.X - u.X, v.Y - u.Y, v.Z - u.Z, 1}
}

func (v Vec4) Scale(k float64) Vec4 {
	return Vec4{
		v.X * k,
		v.Y * k,
		v.Z * k,
		v.W * k,
	}
}

func (v Vec4) Dot(u Vec4) float64 {
	return v.X*u.X + v.Y*u.Y + v.Z*u.Z + v.W*u.W
}

func (v Vec4) Length() float64 {
	return math.Sqrt(v.Dot(v))
}

func (v Vec4) Normalize() Vec4 {
	l := v.Length()
	return Vec4{v.X / l, v.Y / l, v.Z / l, v.W / l}
}

type Vec3 struct {
	X, Y, Z float64
}

func (v Vec3) Sub(u Vec3) Vec3 {
	return Vec3{v.X - u.X, v.Y - u.Y, v.Z - u.Z}
}

func (v Vec3) Cross(u Vec3) Vec3 {
	return Vec3{
		v.Y*u.Z - v.Z*u.Y,
		v.Z*u.X - v.X*u.Z,
		v.X*u.Y - v.Y*u.X,
	}
}

func (v Vec3) Dot(u Vec3) float64 {
	return v.X*u.X + v.Y*u.Y + v.Z*u.Z
}

func (v Vec3) Length() float64 {
	return math.Sqrt(v.Dot(v))
}

func (v Vec3) Normalize() Vec3 {
	l := v.Length()
	return Vec3{v.X / l, v.Y / l, v.Z / l}
}

func (v Vec3) Scale(k float64) Vec3 {
	return Vec3{v.X * k, v.Y * k, v.Z * k}
}

type Mat4 [4][4]float64

func translate(v Vec3) Mat4 {
	return Mat4{
		{1, 0, 0, v.X},
		{0, 1, 0, v.Y},
		{0, 0, 1, v.Z},
		{0, 0, 0, 1},
	}
}

func lookat(eye, center, up Vec3) Mat4 {
	f := center.Sub(eye).Normalize()
	s := up.Cross(f)
	u := s.Cross(f)
	m := Mat4{
		{s.X, s.Y, s.Z, 0},
		{u.X, u.Y, u.Z, 0},
		{-f.X, -f.Y, -f.Z, 0},
		{0, 0, 0, 1},
	}
	t := translate(Vec3{-eye.X, -eye.Y, -eye.Z})
	return m.Mul(t)
}

func perspective(fovy, aspect, near, far float64) Mat4 {
	f := math.Tan(fovy / 2)
	z := near - far
	return Mat4{
		{1 / (f * aspect), 0, 0, 0},
		{0, 1 / f, 0, 0},
		{0, 0, (-near - far) / z, 2 * far * near / z},
		{0, 0, 1, 0},
	}
}

func ortho(l, r, b, t, n, f float64) Mat4 {
	sx := 2 / (r - l)
	sy := 2 / (t - b)
	sz := 2 / (f - n)

	tx := -(r + l) / (r - l)
	ty := -(t + b) / (t - b)
	tz := -(f + n) / (f - n)

	return Mat4{
		{sx, 0, 0, tx},
		{0, sy, 0, ty},
		{0, 0, sz, tz},
		{0, 0, 0, 1},
	}
}

func identity() Mat4 {
	return Mat4{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
}

func screenSpace(w, h float64) Mat4 {
	hw := w / 2
	hh := h / 2
	return Mat4{
		{hw, 0, 0, hw},
		{0, -hh, 0, hh},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
}

func (m Mat4) Mul(n Mat4) Mat4 {
	var p Mat4

	for i := range m {
		for j := range m[i] {
			for k := range m[j] {
				p[i][j] += m[i][k] * n[k][j]
			}
		}
	}

	return p
}

func (m Mat4) Mulv(v Vec4) Vec4 {
	return Vec4{
		m[0][0]*v.X + m[0][1]*v.Y + m[0][2]*v.Z + m[0][3]*v.W,
		m[1][0]*v.X + m[1][1]*v.Y + m[1][2]*v.Z + m[1][3]*v.W,
		m[2][0]*v.X + m[2][1]*v.Y + m[2][2]*v.Z + m[2][3]*v.W,
		m[3][0]*v.X + m[3][1]*v.Y + m[3][2]*v.Z + m[3][3]*v.W,
	}
}
