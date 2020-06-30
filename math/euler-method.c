// https://en.wikipedia.org/wiki/Euler_method
#include <stdio.h>
#include <math.h>

/*

Approximate first order differential equation with a given initial condition

dy/dt = F(t, y(t))
y(t0) = y0

y(t+h) ~ y(t) + h*f(t, y(t))
hence iterative solution is
y_n+1 = y_n + hf(t_n, y_n)

f: derivative of F(t, y)
y0: initial condition
a: t0
b: t1
h: step size

We need to start at the initial condition for a and iterate to b using step size h
The accuracy is determined by the step size.

*/
void
fweuler(double (*f)(double, double), double y0, double a, double b, double h, void (*cb)(double, double, void *), void *ud)
{
	double y, t;

	y = y0;
	for (t = a; t <= b; t += h) {
		cb(t, y, ud);
		y += h * f(t, y);
	}
}

// http://rosettacode.org/wiki/Euler_method
// F(t) = (T-T0) * exp(-k*t)
double
f1(double y, double t)
{
	(void)y;
	return -0.07 * (t - 20);
}

// F(t) = exp(t)
double
f2(double y, double t)
{
	(void)y;
	return t;
}

void
print(double t, double y, void *u)
{
	(void)u;
	printf("%.6f %.6f\n", t, y);
}

int
main(void)
{
	fweuler(f1, 100, 0, 256, 0.0001, print, NULL);
	fweuler(f2, 1, 0, 4, 0.00001, print, NULL);
	return 0;
}
