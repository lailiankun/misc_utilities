package main

import (
	"image"
	"image/color"
	"image/draw"
	"log"
	"math"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/qeedquan/go-media/math/f64"
	"github.com/qeedquan/go-media/sdl"
	"github.com/qeedquan/go-media/sdl/sdlgfx"
	"github.com/qeedquan/go-media/sdl/sdlimage/sdlcolor"
)

var (
	window   *sdl.Window
	renderer *sdl.Renderer
	texture  *sdl.Texture
	canvas   *image.RGBA
	zbuffer  [][]float64
	fps      sdlgfx.FPSManager

	at     = f64.Vec3{8, 5, 4}
	eye    = f64.Vec3{0, 0, -2}
	camera f64.Mat4
)

func main() {
	runtime.LockOSThread()
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(0)
	log.SetPrefix("")
	initSDL()

	for {
		event()
		blit()
		fps.Delay()
	}
}

func ck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func initSDL() {
	err := sdl.Init(sdl.INIT_VIDEO)
	ck(err)

	w, h := 800, 600
	window, renderer, err = sdl.CreateWindowAndRenderer(w, h, sdl.WINDOW_RESIZABLE)
	ck(err)

	resize(w, h)

	window.SetTitle("Cube")

	fps.Init()
	fps.SetRate(60)
}

func resize(w, h int) {
	var err error

	if texture != nil {
		texture.Destroy()
	}

	texture, err = renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, w, h)
	ck(err)

	canvas = image.NewRGBA(image.Rect(0, 0, w, h))
}

func event() {
	for {
		ev := sdl.PollEvent()
		if ev == nil {
			break
		}

		const S = 0.25
		switch ev := ev.(type) {
		case sdl.QuitEvent:
			os.Exit(0)
		case sdl.KeyDownEvent:
			switch ev.Sym {
			case sdl.K_ESCAPE:
				os.Exit(0)
			case sdl.K_LEFT:
				at.X -= S
			case sdl.K_RIGHT:
				at.X += S
			case sdl.K_UP:
				at.Y -= S
			case sdl.K_DOWN:
				at.Y += S
			case sdl.K_z:
				at.Z += S
			case sdl.K_x:
				at.Z -= S
			case sdl.K_KP_2:
				eye.Y -= S
			case sdl.K_KP_4:
				eye.X += S
			case sdl.K_KP_6:
				eye.X -= S
			case sdl.K_KP_8:
				eye.Y += S
			case sdl.K_KP_7:
				eye.Z -= S
			case sdl.K_KP_9:
				eye.Z += S
			}
		case sdl.WindowEvent:
			switch ev.Event {
			case sdl.WINDOWEVENT_RESIZED:
				resize(int(ev.Data[0]), int(ev.Data[1]))
			}
		}
	}
}

func blit() {
	camera.LookAt(eye, at, f64.Vec3{0, 1, 0})
	renderer.SetDrawColor(sdlcolor.Black)
	renderer.Clear()
	draw.Draw(canvas, canvas.Bounds(), image.NewUniform(color.Black), image.ZP, draw.Src)
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			for k := 0; k < 4; k++ {
				blitBlock(i, j, k)
			}
		}
	}
	texture.Update(nil, canvas.Pix, canvas.Stride)
	renderer.Copy(texture, nil, nil)
	renderer.Present()
}

func blitBlock(i, j, k int) {
	w, h, err := renderer.OutputSize()
	ck(err)

	const S = 0.10
	l_length := -S
	l_height := S
	l_width := S

	v := []f64.Vec3{
		{l_length, -l_height, -l_width},
		{-l_length, -l_height, -l_width},
		{-l_length, l_height, -l_width},
		{l_length, l_height, -l_width},

		{-l_length, -l_height, l_width},
		{l_length, -l_height, l_width},
		{l_length, l_height, l_width},
		{-l_length, l_height, l_width},

		{l_length, -l_height, l_width},
		{l_length, -l_height, -l_width},
		{l_length, l_height, -l_width},
		{l_length, l_height, l_width},

		{-l_length, -l_height, -l_width},
		{-l_length, -l_height, l_width},
		{-l_length, l_height, l_width},
		{-l_length, l_height, -l_width},

		{-l_length, -l_height, -l_width},
		{l_length, -l_height, -l_width},
		{l_length, -l_height, l_width},
		{-l_length, -l_height, l_width},

		{l_length, l_height, -l_width},
		{-l_length, l_height, -l_width},
		{-l_length, l_height, l_width},
		{l_length, l_height, l_width},
	}

	var modelMatrix, viewMatrix, projMatrix, screenMatrix f64.Mat4
	modelMatrix.Translate(float64(i)*S, float64(j)*S, float64(k)*S)
	viewMatrix = camera
	projMatrix.Perspective(math.Pi/4, float64(w)/float64(h), 1, 100)
	screenMatrix.Viewport(0, 0, float64(w), float64(h))

	for i := range v {
		v[i] = modelMatrix.Transform3(v[i])
		v[i] = viewMatrix.Transform3(v[i])
		v[i] = projMatrix.Transform3(v[i])
		v[i] = screenMatrix.Transform3(v[i])
	}

	for i := 0; i < len(v); i += 4 {
		Line(canvas, int(v[i].X), int(v[i].Y), int(v[i+1].X), int(v[i+1].Y), color.White)
		Line(canvas, int(v[i+1].X), int(v[i+1].Y), int(v[i+2].X), int(v[i+2].Y), color.White)
		Line(canvas, int(v[i+2].X), int(v[i+2].Y), int(v[i+3].X), int(v[i+3].Y), color.White)
		Line(canvas, int(v[i+3].X), int(v[i+3].Y), int(v[i].X), int(v[i].Y), color.White)
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func Line(img draw.Image, x0, y0, x1, y1 int, color color.Color) {
	dx := abs(x1 - x0)
	dy := abs(y1 - y0)

	var sx, sy int
	if x0 < x1 {
		sx = 1
	} else {
		sx = -1
	}
	if y0 < y1 {
		sy = 1
	} else {
		sy = -1
	}
	err := dx - dy

	var e2 int
	for {
		img.Set(x0, y0, color)
		if x0 == x1 && y0 == y1 {
			return
		}
		e2 = 2 * err
		if e2 > -dy {
			err = err - dy
			x0 = x0 + sx
		}
		if e2 < dx {
			err = err + dx
			y0 = y0 + sy
		}
	}
}
