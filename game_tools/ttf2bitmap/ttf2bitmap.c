#define _GNU_SOURCE
#include <errno.h>
#include <unistd.h>
#include <limits.h>
#include <getopt.h>

#define STB_DEFINE
#define STB_TRUETYPE_IMPLEMENTATION
#define STB_IMAGE_WRITE_IMPLEMENTATION
#include "stb.h"
#include "stb_truetype.h"
#include "stb_image_write.h"

struct {
	int width;
	int height;
	int start;
	int end;
	float size;
	ulong fg;
	ulong bg;
} opt = {
	.width = 256,
	.height = 256,
	.start = 0,
	.end = 128,
	.size = 16.0f,
	.fg = 0xffffffff,
	.bg = 0xff000000,
};

void
usage(void)
{
	fprintf(stderr, "usage: [options] <ttf> <file>\n");
	fprintf(stderr, "    -b: background color (default: %#lx)\n", opt.bg);
	fprintf(stderr, "    -f: foreground color (default: %#lx)\n", opt.fg);
	fprintf(stderr, "    -h: image height (default: %d)\n", opt.height);
	fprintf(stderr, "    -r: character range (default: %d-%d)\n", opt.start, opt.end);
	fprintf(stderr, "    -s: font size (default: %f)\n", opt.size);
	fprintf(stderr, "    -w: image width (default: %d)\n", opt.width);
	exit(1);
}

void *
ecalloc(size_t nmemb, size_t size)
{
	void *ptr;

	if (nmemb == 0)
		nmemb = 1;
	if (size == 0)
		size = 1;
	ptr = calloc(nmemb, size);
	if (!ptr)
		stb_fatal("failed to allocate memory");
	return ptr;
}

void
put4(uchar *p, uint32 v)
{
	p[0] = v;
	p[1] = (v >> 8) & 0xff;
	p[2] = (v >> 16) & 0xff;
	p[3] = (v >> 24) & 0xff;
}

void
scale(stbtt_aligned_quad *q, int w, int h)
{
	q->s0 *= w;
	q->s1 *= w;
	q->t0 *= h;
	q->t1 *= h;
}

void
parse_options(int *argc, char ***argv)
{
	char *endptr;
	int c;

	while ((c = getopt(*argc, *argv, "b:f:w:h:s:r:?")) != -1) {
		switch (c) {
		case 'b':
			opt.bg = strtoul(optarg, &endptr, 0);
			break;
		case 'f':
			opt.fg = strtoul(optarg, &endptr, 0);
			break;
		case 'w':
			opt.width = atoi(optarg);
			break;
		case 'h':
			opt.height = atoi(optarg);
			break;
		case 's':
			opt.size = atof(optarg);
			break;
		case 'r':
			sscanf(optarg, "%d-%d", &opt.start, &opt.end);
			break;
		case '?':
		default:
			usage();
		}
	}

	if (opt.width <= 0 || opt.height <= 0 || opt.width > 65535 || opt.height > 65535)
		stb_fatal("invalid image dimension specified");

	if (opt.end < opt.start)
		stb_swap(&opt.start, &opt.end, sizeof(opt.start));

	if (opt.start < 0 || opt.end < 0 || opt.start >= 255 || opt.end >= 255)
		stb_fatal("unsupported character range");

	if (opt.size <= 0)
		stb_fatal("invalid font size");

	*argc -= optind;
	*argv += optind;
}

void
make_char_table(char *name, stbtt_bakedchar *cdata)
{
	stbtt_aligned_quad q;
	FILE *fp;
	float x, y;
	int i;

	fp = fopen(name, "wt");
	if (!fp)
		stb_fatal("failed to open character table output file %s: %s", strerror(errno));

	fprintf(fp, "# Generated by ttf2bitmap\n");
	fprintf(fp, "# image_width image_height font_size\n");
	fprintf(fp, "# character xoff yoff xadvance texcoord_min_x texcoord_min_y texcoord_max_x texcoord_max_y\n");
	fprintf(fp, "%d %d %f\n", opt.width, opt.height, opt.size);
	x = y = 0;
	for (i = 0; i < opt.end - opt.start + 1; i++) {
		stbtt_bakedchar *b;

		stbtt_GetBakedQuad(cdata, opt.width, opt.height, i, &x, &y, &q, 1);
		scale(&q, opt.width, opt.height);

		b = cdata + i;
		fprintf(fp, "%d %f %f %f %f %f %f %f\n",
		        i, b->xoff, b->yoff, b->xadvance, q.s0, q.t0, q.s1, q.t1);
	}
	fclose(fp);
}

void
make_bitmap(char *ttf, char *img)
{
	stbtt_bakedchar *cdata;
	uchar *bitmap, *rgba, *font, *p;
	char *chartab;
	int x, y, rv;
	size_t length;

	font = stb_file(ttf, &length);
	if (!font)
		stb_fatal("failed to load font %s: %s", ttf, strerror(errno));

	bitmap = ecalloc(opt.width, opt.height);
	rgba = ecalloc(opt.width * 4, opt.height);
	cdata = ecalloc(opt.end - opt.start + 1, sizeof(*cdata));
	if ((rv = stbtt_BakeFontBitmap(font, 0, opt.size, bitmap, opt.width, opt.height, opt.start, opt.end - opt.start + 1, cdata)) < 0)
		stb_fatal("failed to bake font %s", ttf);

	p = rgba;
	for (y = 0; y < opt.height; y++) {
		for (x = 0; x < opt.width; x++) {
			if (bitmap[y * opt.width + x])
				put4(p, opt.fg);
			else
				put4(p, opt.bg);
			p += 4;
		}
	}

	if (!img) {
		img = ecalloc(strlen(ttf) + 4, 1);
		stb_replacedir(img, ttf, ".");
		stb_replaceext(img, img, "png");
	}

	if (stb_suffixi(img, "bmp"))
		rv = stbi_write_bmp(img, opt.width, opt.height, 4, rgba);
	else if (stb_suffixi(img, "tga"))
		rv = stbi_write_tga(img, opt.width, opt.height, 4, rgba);
	else
		rv = stbi_write_png(img, opt.width, opt.height, 4, rgba, opt.width * 4);

	if (rv == 0)
		stb_fatal("failed to write image to file %s", img);

	chartab = ecalloc(strlen(img) + 4, 1);
	stb_replaceext(chartab, img, "tbl");
	make_char_table(chartab, cdata);
}

int
main(int argc, char *argv[])
{
	parse_options(&argc, &argv);
	if (argc < 1 || argc > 2)
		usage();

	if (argc < 2)
		make_bitmap(argv[0], NULL);
	else
		make_bitmap(argv[0], argv[1]);

	return 0;
}
