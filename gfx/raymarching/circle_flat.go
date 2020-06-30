// based on https://www.shadertoy.com/view/llt3R4
package main

import (
	"image"
	"image/color"
	"image/draw"
	"log"
	"math"
	"os"
	"runtime"

	"github.com/qeedquan/go-media/math/f64"
	"github.com/qeedquan/go-media/sdl"
	"github.com/qeedquan/go-media/sdl/sdlimage/sdlcolor"
)

const (
	RAD   = math.Pi / 180
	MIN   = 0.0
	MAX   = 100.0
	STEPS = 255
	EPS   = 1e-4
)

var (
	window   *sdl.Window
	renderer *sdl.Renderer
	texture  *sdl.Texture
	canvas   *image.RGBA
	eye      = f64.Vec3{0, 0, 5}
)

func main() {
	runtime.LockOSThread()

	log.SetFlags(0)
	log.SetPrefix("")

	initSDL()
	for {
		event()
		blit()
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

	texture, err = renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, w, h)
	ck(err)

	canvas = image.NewRGBA(image.Rect(0, 0, w, h))

	window.SetTitle("Raymarching Sphere with Flat Shading")

	renderer.SetLogicalSize(w, h)
}

func event() {
	const S = 0.25
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
			case sdl.K_UP:
				eye.Y -= S
			case sdl.K_DOWN:
				eye.Y += S
			case sdl.K_LEFT:
				eye.X -= S
			case sdl.K_RIGHT:
				eye.X += S
			case sdl.K_z:
				eye.Z -= S
			case sdl.K_x:
				eye.Z += S
			}
		}
	}
}

func blit() {
	renderer.SetDrawColor(sdlcolor.Black)
	renderer.Clear()
	draw.Draw(canvas, canvas.Bounds(), image.NewUniform(color.Black), image.ZP, draw.Src)
	raymarch()
	texture.Update(nil, canvas.Pix, canvas.Stride)
	renderer.Copy(texture, nil, nil)
	renderer.Present()
}

func raymarch() {
	w, h := renderer.LogicalSize()

	L := f64.Vec3{0, 0, 1}
	r := f64.Vec2{float64(w), float64(h)}
	p := f64.Vec2{0, 0}
	for p.Y = 0; p.Y < float64(h); p.Y++ {
		for p.X = 0; p.X < float64(w); p.X++ {
			dir := raydir(45, r, p)
			dist := distance(eye, dir, MIN, MAX)
			if dist > MAX-EPS {
				canvas.SetRGBA(int(p.X), int(p.Y), color.RGBA{0, 0, 0, 255})
			} else {
				I := L.Dot(eye.AddScale(dir, dist))
				I = f64.Clamp(I*255, 0, 255)
				c := color.RGBA{
					0, 0, uint8(I), 255,
				}
				canvas.SetRGBA(int(p.X), int(p.Y), c)
			}
		}
	}
}

func raydir(fovy float64, size, fragCoord f64.Vec2) f64.Vec3 {
	x := fragCoord.X - size.X/2
	y := fragCoord.Y - size.Y/2
	z := size.Y / math.Tan(fovy*RAD/2)
	return f64.Vec3{x, y, -z}.Normalize()
}

func distance(eye, dir f64.Vec3, start, end float64) float64 {
	depth := start
	for i := 0; i < STEPS; i++ {
		dist := sdf(eye.AddScale(dir, depth))
		if dist < EPS {
			return depth
		}
		depth += dist
		if depth >= end {
			return end
		}
	}
	return end
}

func sdf(p f64.Vec3) float64 {
	return sphere(p)
}

func sphere(p f64.Vec3) float64 {
	return p.Len() - 1
}
