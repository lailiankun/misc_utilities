/*

 https://en.wikipedia.org/wiki/Rauzy_fractal
 https://bitbucket.org/snippets/NeuralOutlet/8qEa/rauzy-fractal-problem
 https://math.stackexchange.com/questions/1264647/what-method-is-used-for-projecting-the-rauzy-fractal

 the substitution: 1 -> 12, 2 -> 13, 3 -> 1
 has the incidence matrix:

   1 1 0
   1 0 1
   1 0 0

 which has the approx eigenvector:

   ω = ⟨ω1,ω2,ω3⟩ = ⟨-0.412+0.6i, 0.223-1.12i, 1⟩

  and the mapping function is the count of amount of instances of character N in the substitution string U:

  sum foreach n: |U|_n * ω_n

  the resulting value is a single complex number.

*/

package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/color"
	"log"
	"os"

	"github.com/qeedquan/go-media/image/imageutil"
	"github.com/qeedquan/go-media/math/f64"
)

var (
	width   = flag.Int("w", 1024, "image width")
	height  = flag.Int("h", 1024, "image height")
	iters   = flag.Int("i", 0, "number of iterations")
	variant = flag.Int("v", 0, "use rule variant")
	spacing = flag.Float64("s", 100, "spacing")
	radius  = flag.Float64("r", 4, "dot radius")
	outfile = flag.String("o", "rauzy.png", "output file")
)

func main() {
	flag.Usage = usage
	flag.Parse()

	s := "12131211213121213121121312131211213121213121121312112131212131211213121312112131212131"
	if flag.NArg() >= 1 {
		s = flag.Arg(0)
	}
	s = expand(s, *iters, *variant)

	r := image.Rect(0, 0, *width, *height)
	m := image.NewRGBA(r)
	render(m, s, *spacing, *radius)

	err := imageutil.WriteRGBAFile(*outfile, m)
	if err != nil {
		log.Fatal(err)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: [options] string")
	flag.PrintDefaults()
	os.Exit(2)
}

func render(img *image.RGBA, str string, spacing, radius float64) {
	pal := []color.RGBA{
		{255, 0, 0, 255},
		{0, 255, 0, 255},
		{0, 0, 255, 255},
	}

	b := img.Bounds()
	fwc := float64(b.Dx()) / 2
	fhc := float64(b.Dy()) / 2

	M := f64.Mat4{}
	M.RotateZ(f64.Deg2Rad(90))

	for i := range str {
		pt := evalpt(str[i:])

		tp := f64.Vec4{real(pt), imag(pt), 0, 1}
		tp = M.Transform(tp)

		fx := tp.X*spacing + fwc
		fy := tp.Y*spacing + fhc
		px := int(fx)
		py := int(fy)

		if !('1' <= str[i] && str[i] <= '3') {
			continue
		}

		pi := str[i] - '1'
		imageutil.FilledCircle(img, px, py, f64.Iround(radius), pal[pi])
	}
}

func evalpt(str string) complex128 {
	eig := [...]complex128{
		-0.412 + 0.61i,
		0.223 - 1.12i,
		1,
	}

	var cnt [3]float64
	for _, r := range str {
		if !('1' <= r && r <= '3') {
			continue
		}
		cnt[r-'1']++
	}

	sum := 0.0 + 0i
	for i := range cnt {
		sum += complex(cnt[i], 0) * eig[i]
	}
	return sum
}

func expand(str string, iters, variant int) string {
	tab := [][3]string{
		{"12", "13", "1"},
		{"12", "31", "1"},
		{"12", "23", "312"},
		{"123", "1", "31"},
		{"123", "1", "1132"},
	}
	if !(0 <= variant && variant < len(tab)) {
		variant = 0
	}

	for i := 0; i < iters; i++ {
		w := new(bytes.Buffer)
		for _, r := range str {
			c := r - '0'
			if !(1 <= c && c <= 3) {
				continue
			}
			w.WriteString(tab[variant][c-1])
		}
		str = w.String()
	}
	return str
}
