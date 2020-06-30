// http://www.songho.ca/opengl/gl_projectionmatrix.html
package main

import (
	"fmt"
	"math"
)

func main() {
	run(1, 20)
	run(-1, -1500)
	run(-5, 5)
}

func run(n, f float64) {
	// the projection maps [-n,-f] to [-1,1]
	// test points are n -> -1
	// -b/a -> 0 (since z=-a-b/z -> -b/a=z
	// f -> -1
	fmt.Println(project(-n, n, f))
	a := -(f + n) / (f - n)
	b := -2 * f * n / (f - n)
	fmt.Println(project(-b/a, n, f))
	fmt.Println(project(-f, n, f))

	fmt.Println()
	// this is not a linear map, so the densities
	// need to be 1/d^2 so we can see values
	// near the far planes
	for i := 0.0; i < 1.0; i += 1 / (f * f) {
		a := linear(i, n, f)

		// project takes negative coordinate because opengl
		// uses right handed coordinate where -z is forward
		b := project(-a, n, f)
		b = 0.5*b + 0.5

		// project and linear are inverses of each other
		// so we should get back the same values up until
		// some epsilon due to floating point errors
		fmt.Println(i, a, b, math.Abs(i-b))
	}
	a = linear(1, n, f)
	b = project(-a, n, f)
	b = 0.5*b + 0.5
	fmt.Println(1, a, b, math.Abs(1-b))
}

func linear(d, n, f float64) float64 {
	// depth is from [0, 1]
	// map it to [-1, 1] since opengl
	//
	d = d*2 - 1

	// given d between [-1, 1] in 1/z space
	// need to map it back to linear space
	// this will map values from [0, 1] depth values
	// into the [n,f] range, but with 1/z densities
	// so most values between [0,1] will map nearer to
	// the near plane and only a small fraction of values near 1
	// map more closely to the far values assuming finite precision

	// derivation: solve for the inverse of the project function
	// a = -(f+n) / (f-n)
	// b = -2fn   / (f-n)

	// to simplify notation let
	// d = z_ndc
	// z = z_eye

	// d = (a*z + b) / -z
	// solving for z
	// -dz = az + b
	// -b = az + dz
	// -b/(a+d) = z

	// plug in b
	// -2fn / ((f-n)*(a+d)) = z
	// focusing on the denominator
	// a simplifies to -f+n
	// a = (f-n)*-(f+n)(f-n) = -f-n
	// plug it back in, denominator becomes
	// -f-n+d*(f-n)

	// so we have
	// z = -2fn / (-f-n+d(f-n))
	// if we multiply by -1/-1 we get
	// z = 2fn / (f+n-d(f-n))

	return (2 * n * f) / (f + n - d*(f-n))
}

func project(z, n, f float64) float64 {
	// z coordinates is [n, f], the projection
	// gives a mapping of values between [-1, 1]

	// project does a perspective mapping of z
	// such that near and far are treated unequally
	// if we wanted a linear remap it would've been
	// d = (z-n) / (f-n) where d is the depth but
	// physically near things have more depth than
	// far things and so the depth buffer should logically
	// discard more far away things than near things
	// we need to remap it in terms of 1/z depth,
	// this gives more weight to near things than far things

	// derivation of 1/z depth basically we just solve
	// for z_ndc with perspective division, the perspective
	// division will give us 1/z depth
	// z_ndc = z_clip/w_clip = (Az_eye + Bw_eye)/(-z_eye)
	// we need to solve for constants A and B
	// the division of -z_eye is the perspective divide from the w
	// coordinate (this is the third row of the perspective projection matrix)

	// in eye space, w_eye is equal to 1 so equation
	// becomes
	// z_ndc = z_clip/w_clip = (Az_eye + B)/-z_eye

	// to find A and B we can use the fact that
	// we want to map [-n,-f] to [-1, 1] so we set
	// z_eye = -n and z_eye = -f and get
	// -1 = (-An + B)/n
	//  1 = (-Af + B)/f

	// -An + B = -n
	// -Af + B = f
	// use first equation
	// B = An - n
	// plug B into second equation
	// -Af + An - n = f
	// A(-f + n) = f + n
	// A = (f+n)/(-f+n) move -1 to the top
	// A = -(f+n) / (f-n)

	// B = -(nf+nn)/(f-n) - n
	// B = -(nf+nn+nf-nn)/(f-n)
	// B = -2nf/(f-n)

	// projective equation becomes
	// z_ndc = (A*z_eye + B)/-z_eye
	// can simplify to
	// z_ndc = -A - B/z_eye

	a := -(f + n) / (f - n)
	b := -2 * f * n / (f - n)
	return -a - b/z
}
