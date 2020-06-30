// barycentric coordinates can be represented as a matrix form just like with
// euclidean/cartesian space, the difference is that the coordinate basis points
// and not vectors, so we need to represent the basis vector as following
// |v0x v1x v2x|
// |v0y v1y v2y|
// |1   1   1  |

// this is for 2D, for more points, use a bigger square matrix that can hold all of them
// can't use a non-square matrix because they don't have inverses

// v0, v1, and v2 represents the vertices of the simplex
// this matrix is the basis matrix to generate points, if you want to generate
// points inside the simplex (point, line, triangle, etc), multiply the matrix with a point in the range [0, 1]
// so B*p where p is inside [0, 1]
// VertexPoint=Binv*weightVector

// If you invert this matrix, then Binv*P where P is in vertex coordinate space, the result will
// give you the weights of the coefficients for each vertices, so
// weightVector=Binv*VertexPoint

// If you want to represent a line or points, just ignore the other vertex (test the ignored vertex to be near to 0)
// since if you move along 2 vertices, only the weights for those 2 vertex will change, the other will be 0
// for a point, test if both ignored vertices are near 0

package main

import (
	"flag"
	"fmt"
	"image/color"
	"log"
	"math"
	"os"
	"runtime"

	"github.com/qeedquan/go-media/math/f64"
	"github.com/qeedquan/go-media/sdl"
	"github.com/qeedquan/go-media/sdl/sdlgfx"
	"github.com/qeedquan/go-media/sdl/sdlimage/sdlcolor"
)

type Simplex struct {
	mode int
	vert [3]f64.Vec2
}

var (
	window   *sdl.Window
	renderer *sdl.Renderer
	simplex  Simplex
	grab     int
)

var conf struct {
	width  int
	height int
	radius int
	aspect float64
}

func main() {
	runtime.LockOSThread()
	log.SetFlags(0)
	parseFlags()
	initSDL()
	loop()
}

func ck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func parseFlags() {
	flag.IntVar(&conf.width, "width", 1024, "window width")
	flag.IntVar(&conf.height, "height", 768, "window height")
	flag.Parse()
	conf.radius = 10
	conf.aspect = float64(conf.width) / float64(conf.height)
}

func initSDL() {
	err := sdl.Init(sdl.INIT_VIDEO)
	ck(err)

	window, renderer, err = sdl.CreateWindowAndRenderer(conf.width, conf.height, sdl.WINDOW_RESIZABLE)
	ck(err)

	window.SetTitle("Barycentric Basis")
	renderer.SetLogicalSize(conf.width, conf.height)

	sdlgfx.SetFont(sdlgfx.Font10x20, 10, 20)
}

func loop() {
	reset()
	for {
		event()
		blit()
	}
}

func reset() {
	grab = -1
	simplex = Simplex{
		mode: 3,
		vert: [3]f64.Vec2{
			{0.7, 0.5},
			{0.7 - 0.3, 0.5 + 0.3},
			{0.7 + 0.3, 0.5 + 0.3},
		},
	}
}

func event() {
	for {
		ev := sdl.PollEvent()
		if ev == nil {
			break
		}
		switch ev := ev.(type) {
		case sdl.QuitEvent:
			os.Exit(0)
		case sdl.KeyDownEvent:
			switch ev.Sym {
			case sdl.K_ESCAPE:
				os.Exit(0)
			case sdl.K_r:
				reset()
			case sdl.K_1:
				simplex.mode = 1
			case sdl.K_2:
				simplex.mode = 2
			case sdl.K_3:
				simplex.mode = 3
			}
		case sdl.MouseMotionEvent:
			if grab >= 0 {
				simplex.vert[grab] = f64.Vec2{
					float64(ev.X) / float64(conf.width) * conf.aspect,
					float64(ev.Y) / float64(conf.height),
				}
			}
		case sdl.MouseButtonDownEvent:
			if ev.Button == 1 {
				grip(&simplex, int(ev.X), int(ev.Y))
			}
		case sdl.MouseButtonUpEvent:
			if ev.Button == 1 {
				grab = -1
			}
		}
	}
}

func grip(t *Simplex, mx, my int) {
	gx := float64(mx) / float64(conf.width)
	gy := float64(my) / float64(conf.height)
	rd := float64(conf.radius) / (float64(conf.width) / conf.aspect)
	for i, v := range t.vert {
		rx := v.X / conf.aspect
		ry := v.Y
		if (gx-rx)*(gx-rx)+(gy-ry)*(gy-ry) < rd*rd {
			grab = i
			return
		}
	}
}

func blit() {
	renderer.SetDrawColor(sdl.Color{88, 88, 88, 255})
	renderer.Clear()
	simplex.Draw()
	renderer.Present()
}

func (t *Simplex) Draw() {
	x1 := math.MaxFloat64
	y1 := math.MaxFloat64
	x2 := -math.MaxFloat64
	y2 := -math.MaxFloat64
	for _, v := range t.vert {
		x1 = math.Min(x1, v.X)
		x2 = math.Max(x2, v.X)
		y1 = math.Min(y1, v.Y)
		y2 = math.Max(y2, v.Y)
	}
	x1 *= float64(conf.width) / conf.aspect
	y1 *= float64(conf.height)
	x2 *= float64(conf.width) / conf.aspect
	y2 *= float64(conf.height)

	for y := y1; y <= y2; y++ {
		for x := x1; x <= x2; x++ {
			px := x / (float64(conf.width) / conf.aspect)
			py := y / float64(conf.height)
			b := t.Coord(f64.Vec2{px, py})
			// mode 1 = point
			// mode 2 = line
			// mode 3 = triangle
			switch {
			case t.mode == 1 && math.Abs(b.Y) <= 1e-3 && math.Abs(b.Z) <= 1e-3,
				t.mode == 2 && 0 <= b.X && b.X <= 1 && 0 <= b.Y && b.Y <= 1 && math.Abs(b.Z) <= 1e-2,
				t.mode == 3 && b.X >= 0 && b.Y >= 0 && b.Z >= 0:
				renderer.SetDrawColor(sdl.Color{
					uint8(255 * b.X),
					uint8(255 * b.Y),
					uint8(255 * b.Z),
					255,
				})
				renderer.DrawPoint(int(x), int(y))
			}
		}
	}
	c := []sdl.Color{
		sdlcolor.Red,
		sdlcolor.Blue,
		sdlcolor.Green,
	}
	for i := 0; i < t.mode; i++ {
		v := t.vert[:]
		x := v[i].X * float64(conf.width) / conf.aspect
		y := v[i].Y * float64(conf.height)
		sdlgfx.FilledCircle(renderer, int(x), int(y), conf.radius, c[i])
	}

	y := 20
	for i := 0; i < t.mode; i++ {
		str := fmt.Sprintf("%v", t.vert[i])
		sdlgfx.String(renderer, 10, y, color.RGBA{255, 255, 255, 255}, str)
		y += 20
	}
}

func (t *Simplex) Coord(p f64.Vec2) f64.Vec3 {
	v := t.vert[:]
	B := f64.Mat3{
		{v[0].X, v[1].X, v[2].X},
		{v[0].Y, v[1].Y, v[2].Y},
		{1, 1, 1},
	}
	B.Inverse()
	return B.Transform(f64.Vec3{p.X, p.Y, 1})
}
