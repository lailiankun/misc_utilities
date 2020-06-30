package main

import (
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

var (
	window      *sdl.Window
	renderer    *sdl.Renderer
	particle    f64.Vec3
	nearDivider f64.Vec3
	nearPlane   [3]f64.Vec3
	farDist     float64
	eye, center f64.Vec3
	zbuffer     []float64
	fovy        float64
)

func main() {
	runtime.LockOSThread()
	log.SetFlags(0)
	log.SetPrefix("frustum-view: ")
	initSDL()
	reset()
	for {
		event()
		update()
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

	w, h := 1280, 800
	wflag := sdl.WINDOW_RESIZABLE
	window, renderer, err = sdl.CreateWindowAndRenderer(w, h, wflag)
	ck(err)

	zbuffer = make([]float64, w*h)

	sdlgfx.SetFont(sdlgfx.Font10x20, 10, 20)
	window.SetTitle("Frustum View")
}

func reset() {
	particle = f64.Vec3{}
	nearDivider = f64.Vec3{-.7, 0, 1}
	nearPlane = [3]f64.Vec3{
		{0, -1, -3},
		{-1, 1, -3},
		{1, 1, -3},
	}
	eye = f64.Vec3{0, 0, -2}
	center = f64.Vec3{0, 0, -1}
	farDist = 5
	fovy = math.Pi / 2
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
			case sdl.K_SPACE:
				reset()
			case sdl.K_1:
				fovy -= 0.1
			case sdl.K_2:
				fovy += 0.1
			}
		case sdl.WindowEvent:
			if ev.Event == sdl.WINDOWEVENT_RESIZED {
				zbuffer = make([]float64, ev.Data[0]*ev.Data[1])
			}
		case sdl.MouseMotionEvent:
			view, _, _ := getCurrentViewport()
			rx := f64.Clamp(float64(ev.Xrel), -0.1, 0.1)
			ry := f64.Clamp(float64(ev.Yrel), -0.1, 0.1)

			if view == 2 {
				if ev.State&sdl.BUTTON_LMASK != 0 {
					eye.X += rx
					eye.Y += ry

					center.X += rx
					center.Y += ry
				}
				if ev.State&sdl.BUTTON_RMASK != 0 {
					c := polygonCentroid(nearPlane[:])
					B := planeBasis(nearPlane[:])
					X, Y, Z, _ := B.Basis3()

					var QX, QY, QZ f64.Quat
					keys := sdl.GetKeyboardState()
					if keys[sdl.SCANCODE_LSHIFT] != 0 {
						QX = QX.FromAxisAngle(X, rx)

					} else {
						QZ = QZ.FromAxisAngle(Z, rx)

					}
					QY = QY.FromAxisAngle(Y, ry)
					for i := range nearPlane {
						p := nearPlane[i].Sub(c)
						p = QX.Transform3(p)
						p = QY.Transform3(p)
						p = QZ.Transform3(p)
						nearPlane[i] = p.Add(c)
					}
				}
			}
		case sdl.MouseWheelEvent:
			view, _, _ := getCurrentViewport()
			if view == 2 {
				eye.Z += f64.Clamp(float64(ev.Y), -0.1, 0.1)
				if math.Abs(eye.Z) <= 1.2 {
					eye.Z = f64.Sign(eye.Z) * 1.2
				}
			}
		}
	}
}

func update() {
	particle = getParticlePos()
}

func blit() {
	for i := range zbuffer {
		zbuffer[i] = math.MaxFloat32
	}
	renderer.Clear()
	for i := 0; i < 3; i++ {
		blitOutline(i)
	}
	for i := 0; i < 2; i++ {
		blitSelectAxis(i)
	}
	blitFrustum()
	renderer.Present()
}

func blitPlaneAxis(p []f64.Vec3) {
	var C, M, N, P f64.Mat4
	N = ndcToScreen(2)
	C.LookAt(eye, center, f64.Vec3{0, 1, 0})
	P.Perspective(fovy, aspectRatio(), .1, 1000)
	M.Mul(&P, &C)
	M.Mul(&N, &M)

	B := planeBasis(p)
	O := f64.Vec3{}
	X := f64.Vec3{.5, 0, 0}
	Y := f64.Vec3{0, .5, 0}
	Z := f64.Vec3{0, 0, .5}

	var T, U f64.Mat4
	T.Translate(2.3, 1.5, 0)
	U.Identity()
	U.Mul(&B, &U)
	U.Mul(&T, &U)
	U.Mul(&C, &U)
	U.Mul(&P, &U)
	U.Mul(&N, &U)
	a := U.Transform3(O)
	b := U.Transform3(X)
	blitLine(2, a, b, color.RGBA{255, 0, 0, 255})
	b = U.Transform3(Y)
	blitLine(2, a, b, color.RGBA{0, 255, 0, 255})
	b = U.Transform3(Z)
	blitLine(2, a, b, color.RGBA{0, 0, 255, 255})
}

// to draw near and far plane
// we transform the points of the square to the basis
// of the plane to orient it in space correctly,
// then for far plane, we move the z coordinate far distance
// away and scale accordingly, since near and far have the same
// basis and far is only z distance away, we only need to scale
// and translate far to match near
func blitPlane(p []f64.Vec3, isNear bool, drawParticle bool) []f64.Vec3 {
	col := color.RGBA{255, 255, 255, 255}

	var C, M, N, P f64.Mat4
	N = ndcToScreen(2)
	C.LookAt(eye, center, f64.Vec3{0, 1, 0})
	P.Perspective(fovy, aspectRatio(), .1, 1000)
	M.Mul(&P, &C)
	M.Mul(&N, &M)

	B := planeBasis(p)

	sq := []f64.Vec3{
		{-1, -1, 0},
		{1, -1, 0},
		{1, 1, 0},
		{-1, 1, 0},
	}
	var nsq []f64.Vec3
	var S f64.Mat4
	nsq = append(nsq, sq...)
	S.Scale(0.25, 0.25, 0.25)
	for i := range sq {
		if !isNear {
			sq[i].X *= farDist
			sq[i].Y *= farDist
			sq[i].Z -= farDist
		}
		sq[i] = S.Transform3(sq[i])
	}

	var sqp []f64.Vec3
	var U f64.Mat4
	for i := range sq {
		U.Identity()
		U.Mul(&B, &U)
		U.Mul(&C, &U)
		U.Mul(&P, &U)
		U.Mul(&N, &U)
		j := (i + 1) % len(sq)
		a := U.Transform3(sq[i])
		b := U.Transform3(sq[j])
		blitLine(2, a, b, col)
		sqp = append(sqp, a)
	}

	if drawParticle {
		U.Identity()
		U.Mul(&B, &U)
		U.Mul(&C, &U)
		U.Mul(&P, &U)
		U.Mul(&N, &U)

		p := particle
		p.Z = f64.LinearRemap(p.Z, -1, 1, 0, -farDist)
		p = S.Transform3(p)

		a := U.Transform3(p)
		blitPoint(int(a.X), int(a.Y), int(a.Z), color.RGBA{26, 24, 32, 255})

		np := p
		np.X /= p.Z
		np.Y /= p.Z
		np.Z = 0
		a = U.Transform3(np)
		blitPoint(int(a.X), int(a.Y), int(a.Z), color.RGBA{26, 34, 255, 255})
	}

	return sqp
}

func blitFrustum() {
	viewport := getViewport(2)

	blitPlaneAxis(nearPlane[:])
	s1 := blitPlane(nearPlane[:], true, false)
	s2 := blitPlane(nearPlane[:], false, true)
	for i := range s1 {
		blitLine(2, s1[i], s2[i], color.RGBA{35, 25, 61, 255})
	}

	x := int(viewport.Min.X) + 10
	y := int(viewport.Min.Y)
	str := fmt.Sprintf("eye(%.2f, %.2f, %.2f)", eye.X, eye.Y, eye.Z)
	sdlgfx.String(renderer, x, y, sdlcolor.White, str)

	y += 20
	str = fmt.Sprintf("center(%.2f, %.2f, %.2f)", center.X, center.Y, center.Z)
	sdlgfx.String(renderer, x, y, sdlcolor.White, str)
}

func blitSelectAxis(view int) {
	// x coordinates under perspective transform
	// on the same "plane" (z is the same over all shapes)
	// with fovy=90 degree view [tan(fovy/2) = 1]
	// gets mapped to [-aspect, aspect]
	// so we need to adjust remap our [-1,1] coordinates to that range
	// we could've just used a projection matrix that just scales
	// x to preserve aspect though since all we really want
	// is to draw in 2d but preserve aspects ratio for coordinates [-1,1]
	aspect := aspectRatio()
	axisLen := 0.5

	var N, P, M, T f64.Mat4
	N = ndcToScreen(view)
	P.Perspective(math.Pi/2, aspect, 1, 1000)
	P[3][2] = -P[3][2]
	if view == 0 {
		T.Translate(-aspect+axisLen+0.1, .35, 0)
	} else {
		T.Translate(-aspect+axisLen+0.1, .8, 0)
	}
	M.Mul(&P, &T)
	M.Mul(&N, &M)

	p0 := f64.Vec3{0, 0, 1}
	p1 := f64.Vec3{-axisLen, 0, 1}
	p0 = M.Transform3(p0)
	p1 = M.Transform3(p1)
	blitArrow(view, p0, p1, sdl.Color{0, 0, 255, 255})

	var col color.RGBA
	p0 = f64.Vec3{0, 0, 1}
	if view == 0 {
		p1 = f64.Vec3{0, axisLen, 1}
		col = color.RGBA{255, 0, 0, 255}
	} else {
		p1 = f64.Vec3{0, -axisLen, 1}
		col = color.RGBA{0, 255, 0, 255}
	}
	p0 = M.Transform3(p0)
	p1 = M.Transform3(p1)
	blitArrow(view, p0, p1, col)

	M.Mul(&N, &P)
	p0 = nearDivider
	p1 = nearDivider
	p0.Y -= .6
	p1.Y += .6
	p0 = M.Transform3(p0)
	p1 = M.Transform3(p1)
	blitLine(view, p0, p1, sdlcolor.Black)

	viewport := getViewport(view)
	var str string
	if view == 0 {
		str = fmt.Sprintf("(X:%.2f, Z:%.2f)", particle.X, particle.Z)
	} else {
		str = fmt.Sprintf("(Y:%.2f, Z:%.2f)", particle.Y, particle.Z)
	}
	x, y := int(viewport.Min.X), int(viewport.Min.Y)
	sdlgfx.String(renderer, x, y, sdlcolor.White, str)

	z := f64.LinearRemap(particle.Z, -1, 1, nearDivider.X, aspect)
	if view == 0 {
		p0 = f64.Vec3{z, particle.X, 1}
	} else {
		p0 = f64.Vec3{z, particle.Y, 1}
	}
	p0 = M.Transform3(p0)
	blitPoint(int(p0.X), int(p0.Y), int(p0.Z), col)
}

func blitArrow(view int, p0, p1 f64.Vec3, col color.RGBA) {
	blitLine(view, p0, p1, col)

	var R1, R2 f64.Mat4
	theta := math.Atan2(p1.Y-p0.Y, p1.X-p0.X)
	R1.RotateZ(theta - f64.Deg2Rad(45*3))
	R2.RotateZ(theta + f64.Deg2Rad(45*3))

	p2 := f64.Vec3{25, 0, 0}
	p2 = R1.Transform3(p2)
	p2 = p2.Add(p1)
	blitLine(view, p1, p2, col)

	p2 = f64.Vec3{25, 0, 0}
	p2 = R2.Transform3(p2)
	p2 = p2.Add(p1)
	blitLine(view, p1, p2, col)
}

func blitPoint(x, y, z int, col color.RGBA) {
	w, _, _ := renderer.OutputSize()
	idx := y*w + x
	if idx >= 0 && idx < len(zbuffer) && float64(z) <= zbuffer[idx] {
		zbuffer[idx] = float64(z)
		sdlgfx.FilledCircle(renderer, x, y, 5, col)
	}
}

func blitLine(view int, p0, p1 f64.Vec3, col color.RGBA) {
	vp := getViewport(view)
	x0, y0, z0 := int(p0.X), int(p0.Y), int(p0.Z)
	x1, y1, z1 := int(p1.X), int(p1.Y), int(p1.Z)

	dx := abs(x1 - x0)
	dy := abs(y1 - y0)
	dz := abs(z1 - z0)
	sx := 1
	if x0 >= x1 {
		sx = -1
	}
	sy := 1
	if y0 >= y1 {
		sy = -1
	}
	sz := 1
	if z0 >= z1 {
		sz = -1
	}
	dm := max(dx, dy, dz)
	i := dm
	x1 = dm / 2
	y1 = dm / 2
	z1 = dm / 2

	for {
		pt := f64.Vec2{float64(x0), float64(y0)}
		if pt.In(vp) {
			blitPoint(x0, y0, z0, col)
		}
		i -= 1
		if i <= 0 {
			break
		}
		x1 -= dx
		if x1 < 0 {
			x1 += dm
			x0 += sx
		}
		y1 -= dy
		if y1 < 0 {
			y1 += dm
			y0 += sy
		}
		z1 -= dz
		if z1 < 0 {
			z1 += dm
			z0 += sz
		}
	}
}

func blitOutline(view int) {
	viewport := getViewport(view)

	rect := sdlRect(viewport)
	gray := uint8(80 + view*30)
	col := sdl.Color{gray, gray, gray, 255}

	renderer.SetDrawColor(col)
	renderer.FillRect(&rect)

	for i := 1.0; i < 6; i++ {
		rect = sdlRect(viewport.Inset(-i))
		renderer.SetDrawColor(sdlcolor.Black)
		renderer.DrawRect(&rect)
	}
}

func getViewport(view int) f64.Rectangle {
	w, h, _ := renderer.OutputSize()
	fw, fh := float64(w), float64(h)
	switch view {
	case 0: // top left split screen
		return f64.Rect(0, 0, fw/2, fh/2)
	case 1: // bottom left split screen
		return f64.Rect(0, fh/2, fw/2, fh)
	case 2: // right split screen
		return f64.Rect(fw/2, 0, fw, fh)
	default:
		panic("unreachable")
	}
}

func getCurrentViewport() (view int, viewport f64.Rectangle, mp f64.Vec2) {
	mx, my, _ := sdl.GetMouseState()
	mp = f64.Vec2{float64(mx), float64(my)}
	view = -1
	for i := 0; i < 3; i++ {
		viewport = getViewport(i)
		if mp.In(viewport) {
			view = i
			break
		}
	}
	return
}

func getParticlePos() f64.Vec3 {
	view, viewport, mp := getCurrentViewport()
	_, _, button := sdl.GetMouseState()
	pos := particle
	if button&sdl.BUTTON_LMASK == 0 {
		return pos
	}

	if view == 0 {
		pos.X = f64.LinearRemap(mp.Y-viewport.Min.Y, 0, viewport.Dy(), -1, 1)
	} else if view == 1 {
		pos.Y = f64.LinearRemap(mp.Y-viewport.Min.Y, 0, viewport.Dy(), -1, 1)
	}

	if view == 0 || view == 1 {
		var N, P, M f64.Mat4
		N = ndcToScreen(view)
		P.Perspective(math.Pi/2, aspectRatio(), 1, 1000)
		P[3][2] = -P[3][2]
		M.Mul(&N, &P)

		near := M.Transform3(nearDivider)
		z := mp.X
		if z < near.X {
			z = near.X
		}
		pos.Z = f64.LinearRemap(z, near.X, viewport.Dx(), -1, 1)
	}

	return pos
}

func planeBasis(p []f64.Vec3) f64.Mat4 {
	X := p[0].Sub(p[1]).Normalize()
	Y := p[0].Sub(p[2]).Normalize()
	Z := X.Cross(Y)
	Y = X.Cross(Z)

	var m f64.Mat4
	m.FromBasis3(X, Y, Z, f64.Vec3{})
	return m
}

func polygonCentroid(p []f64.Vec3) f64.Vec3 {
	A0 := 0.0
	A1 := 0.0
	for i := range p {
		j := (i + 1) % len(p)
		A0 += p[i].X*p[j].Y - p[j].X*p[i].Y
		A1 += p[i].X*p[j].Z - p[j].X*p[i].Z
	}
	A0 *= 0.5
	A1 *= 0.5

	cx, cy, cz := 0.0, 0.0, 0.0
	for i := range p {
		j := (i + 1) % len(p)
		cx += (p[i].X + p[j].X) * (p[i].X*p[j].Y - p[j].X*p[i].Y)
		cy += (p[i].Y + p[j].Y) * (p[i].X*p[j].Y - p[j].X*p[i].Y)
		cz += (p[i].Z + p[j].Z) * (p[i].X*p[j].Z - p[j].X*p[i].Z)
	}

	if A0 != 0 {
		cx *= 1 / (6 * A0)
		cy *= 1 / (6 * A0)
	} else if len(p) > 0 {
		cx = p[0].X
		cy = p[0].Y
	}
	if A1 != 0 {
		cz *= 1 / (6 * A1)
	} else if len(p) > 0 {
		cz = p[0].Z
	}
	return f64.Vec3{cx, cy, cz}
}

func aspectRatio() float64 {
	w, h, _ := renderer.OutputSize()
	return float64(w) / float64(h)
}

func ndcToScreen(view int) f64.Mat4 {
	viewport := getViewport(view)
	x := viewport.Min.X
	y := viewport.Min.Y
	hw := viewport.Dx() / 2
	hh := viewport.Dy() / 2
	n := 1.0
	f := 1000.0
	return f64.Mat4{
		{hw, 0, 0, hw + x},
		{0, hh, 0, hh + y},
		{0, 0, (f - n) / 2, (f - n) / 2},
		{0, 0, 0, 1},
	}
}

func sdlRect(r f64.Rectangle) sdl.Rect {
	return sdl.Rect{
		int32(r.Min.X),
		int32(r.Min.Y),
		int32(r.Dx()),
		int32(r.Dy()),
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func max(x ...int) int {
	v := x[0]
	for i := range x[1:] {
		if v < x[i] {
			v = x[i]
		}
	}
	return v
}
