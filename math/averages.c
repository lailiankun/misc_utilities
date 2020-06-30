#include <stdio.h>
#include <math.h>

double
arithmetic(double *v, size_t n)
{
	size_t i;
	double m;

	if (n == 0)
		return 0;

	m = 0;
	for (i = 0; i < n; i++)
		m += v[i];
	return m / n;
}

double
geometric(double *v, size_t n)
{
	size_t i;
	double m;

	if (n == 0)
		return 0;

	m = 1;
	for (i = 0; i < n; i++)
		m *= v[i];
	return pow(m, 1.0 / n);
}

// https://www.cut-the-knot.org/pythagoras/corollary.shtml
// for positive (a,b) (a+b)/2 >= sqrt(ab)
// (a+b)/2 is known as the arithmetic mean of the numbers a and b; ab−−√ is their geometric mean also known as the mean proportional because if k=ab−−√ then a/k=k/b and vice versa.
void
test_inequalities(void)
{
	size_t i, j;

	for (i = 1; i < 1000; i++) {
		for (j = 1; j < 1000; j++) {
			double v[] = {i, j};
			printf("%f %f\n", arithmetic(v, 2), geometric(v, 2));
		}
	}
}

int
main(void)
{
	test_inequalities();
	return 0;
}
