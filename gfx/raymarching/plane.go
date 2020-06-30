// based on http://fabiensanglard.net/rayTracing_back_of_business_card/
package main

import (
	"image"
	"image/png"
	"log"
	"math"
	"math/rand"
	"os"

	"github.com/qeedquan/go-media/math/f64"
)

func main() {
	log.SetPrefix("plane: ")
	log.SetFlags(0)

	m := image.NewRGBA(image.Rect(0, 0, 512, 512))
	raytrace(m)

	f, err := os.Create("plane.png")
	ck(err)
	ck(png.Encode(f, m))
	ck(f.Close())
}

func ck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func raytrace(m *image.RGBA) {
	r := m.Bounds()
	g := f64.Vec3{-6, -16, 0}.Normalize()

	a := f64.Vec3{0, 0, 1}.CrossNormalize(g)
	a = a.Scale(0.002)

	b := g.CrossNormalize(a)
	b = b.Scale(0.002)

	c := a.Add(b)
	c = c.Scale(-256)
	c = c.Add(g)

	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			p := f64.Vec3{13, 13, 13}
			for r := 0; r <= 64; r++ {
				u := (rand.Float64() - .5) * 99
				v := (rand.Float64() - .5) * 99
				au := a.Scale(u)
				bv := b.Scale(v)
				t := au.Add(bv)

				o := f64.Vec3{17, 16, 8}.Add(t)

				rx := rand.Float64() + float64(x)
				ry := rand.Float64() + float64(y)
				ax := a.Scale(rx)
				by := b.Scale(ry)
				abc := ax.Add(by).Add(c)
				abc = abc.Scale(16)

				tm := t.Neg()
				tm = tm.Add(abc)
				d := tm.Normalize()

				p = sample(o, d).Scale(3.5).Add(p)
			}
			p = p.Scale(1 / 255.0)
			m.Set(r.Max.X-1-x, r.Max.Y-1-y, p)
		}
	}
}

func sample(o, d f64.Vec3) f64.Vec3 {
	t, n, m := trace(o, d)
	if m == 0 {
		c := f64.Vec3{.5, .5, .5}
		p := math.Pow(1-d.Z, 4)
		return c.Scale(p)
	}

	h := o.AddScale(d, t)
	u := f64.Vec3{
		9 + rand.Float64(),
		9 + rand.Float64(),
		16,
	}
	l := u.Add(h.Neg()).Normalize()
	r := d.Add(n.Scale(n.Dot(d) * -2))
	b := l.Dot(n)

	_, _, xm := trace(h, l)
	if b < 0 || xm != 0 {
		b = 0
	}

	nb := 0.0
	if b > 0 {
		nb = 1
	}
	p := math.Pow(l.Dot(r)*nb, 99)

	h = h.Scale(.2)
	if m&1 != 0 {
		v := f64.Vec3{}
		hx := math.Ceil(h.X)
		hy := math.Ceil(h.Y)
		if int(hx+hy)&1 != 0 {
			v = f64.Vec3{0, 0, 0}
		} else {
			v = f64.Vec3{3, 3, 3}
		}
		return v.Scale(b*.2 + .1)
	}
	return f64.Vec3{p, p, p}.Add(sample(h, r).Scale(.5))
}

func trace(o, d f64.Vec3) (t float64, n f64.Vec3, m int) {
	t = 1e9
	m = 0
	p := -o.Z / d.Z
	if p >= 0.01 {
		t = p
		n = f64.Vec3{0, 0, 1}
		m = 1
	}

	return
}
