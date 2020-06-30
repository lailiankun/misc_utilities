// measure rate the length of the shape on each scanline to the next
// by using implicit functions to count
package main

import (
	"fmt"
)

func main() {
	cx, cy := 0, 0
	for i := 0; i < 6; i++ {
		measure(cx, cy, ipow10(i), incircle)
	}

	measure(cx, cy, 1000, insquare)
}

// for circles growth rate between 2 scanlines increases
// in a triangle ramp (linearly goes up, then down)
// the largest scanline value is around radius_of_circle = (largest_scanline_count/2)^2
func incircle(x, y, cx, cy, r int) bool {
	return (x-cx)*(x-cx)+(y-cy)*(y-cy) <= r
}

func insquare(x, y, cx, cy, r int) bool {
	x0 := cx - r
	x1 := cx + r
	y0 := cy - r
	y1 := cy + r
	return x0 <= x && x <= x1 && y0 <= y && y <= y1
}

func measure(cx, cy, r int, in func(x, y, cx, cy, r int) bool) {
	x0 := cx - r
	x1 := cx + r
	y0 := cy - r
	y1 := cy + r

	for y := y0; y < y1; y++ {
		n := 0
		for x := x0; x < x1; x++ {
			if in(x, y, cx, cy, r) {
				n++
			}
		}
		if n > 0 {
			fmt.Printf("%d %d: %d\n", y, r, n)
		}
	}
	fmt.Println()
}

func ipow10(p int) int {
	r := 1
	for i := 1; i < p; i++ {
		r *= 10
	}
	return r
}
