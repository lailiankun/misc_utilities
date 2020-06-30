#include <stdio.h>
#include <stdlib.h>
#include <limits.h>

typedef struct {
	int *data;
	size_t rp, wp;
	size_t len, cap;
} CQ;

int
mkcq(CQ *cq, size_t cap)
{
	cq->data = calloc(sizeof(*cq->data), cap);
	if (cq->data == NULL)
		return -1;
	cq->cap = cap;
	cq->len = 0;
	cq->rp = 0;
	cq->wp = 0;
	return 0;
}

void
cqfree(CQ *cq)
{
	if (cq == NULL)
		return;
	free(cq->data);
}

int
cqenq(CQ *cq, int val)
{
	if (cq->len >= cq->cap)
		return -1;
	cq->data[cq->wp] = val;
	cq->wp = (cq->wp + 1) % cq->cap;
	cq->len++;
	return 0;
}

int
cqdeq(CQ *cq, int *val)
{
	*val = 0;
	if (cq->len == 0)
		return -1;
	*val = cq->data[cq->rp];
	cq->rp = (cq->rp + 1) % cq->cap;
	cq->len--;
	return 0;
}

int
cqpeek(CQ *cq)
{
	if (cq->len == 0)
		return INT_MIN;
	return cq->data[cq->rp];
}

void
cqdump(CQ *cq)
{
	size_t i;

	for (i = 0; i < cq->len; i++)
		printf("%d ", cq->data[(cq->rp + i) % cq->cap]);
	printf("\n");
}

void
lnvenq(CQ *cq, int val)
{
	if (cq->len >= cq->cap) {
		cq->rp = (cq->rp + 1) % cq->cap;
		cq->len--;
	}
	cq->data[cq->wp] = val;
	cq->wp = (cq->wp + 1) % cq->cap;
	cq->len++;
}

// bounded circular queues enqueues up to cap elements
// and dequeues as a regular queue would, but with
// O(1) behavior since we don't need to move the
// array whenever enqueue/dequeue occurs
void
test1(void)
{
	CQ cq;
	int i, val;

	printf("Bounded Circular Queues\n");
	mkcq(&cq, 5);
	cqenq(&cq, 20);
	cqenq(&cq, 30);
	cqenq(&cq, 40);

	cqdeq(&cq, &val);
	printf("%d\n", val);

	cqdeq(&cq, &val);
	printf("%d\n", val);

	for (i = 0; i < 100; i++)
		cqenq(&cq, i);
	for (i = 0; i < 100; i++) {
		if (cqdeq(&cq, &val) < 0)
			break;
		printf("%d %d\n", i, val);
	}
	printf("%zu\n", cq.len);
	cqenq(&cq, 62);
	printf("%zu %d\n", cq.len, cqpeek(&cq));
	cqfree(&cq);

	printf("\n\n");
}

// last nth value circular queue stores the last nth value
// and overwrites oldest entries, useful for rolling statistics
void
test2(void)
{
	CQ cq;
	int i;

	mkcq(&cq, 8);
	printf("Last Nth Value\n");
	for (i = 0; i < 100; i++) {
		lnvenq(&cq, i);
		cqdump(&cq);
	}
	printf("\n\n");
	cqfree(&cq);
}

int
main(void)
{
	test1();
	test2();
	return 0;
}
