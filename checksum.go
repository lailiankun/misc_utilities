// computes various checksums for data
package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"flag"
	"fmt"
	"hash"
	"hash/adler32"
	"hash/crc32"
	"hash/crc64"
	"hash/fnv"
	"io"
	"log"
	"os"
	"strings"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/blake2s"
	"golang.org/x/crypto/md4"
	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"
)

func main() {
	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() < 1 {
		usage()
	}

	h := makeHash(flag.Arg(0))
	if h == nil {
		log.Fatalf("unsupported checksum %q", flag.Arg(0))
	}

	if flag.NArg() < 2 {
		sum(os.Stdin, "-", h)
	} else {
		for i := 1; i < flag.NArg(); i++ {
			name := flag.Arg(i)
			f, err := os.Open(name)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}

			err = sum(f, name, h)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}

			f.Close()
			h.Reset()
		}
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: checksum file ...\n")
	fmt.Fprintln(os.Stderr, "supported checksums:")
	fmt.Fprintln(os.Stderr, "adler32 blake2b-256 blake2b-384 blake2b-512 blakes-256 crc8 crc16-ccitt crc32-ieee crc64-iso crc64-ecma fletch16 fnv32 fnv32a fnv64 fnv64a fnv128 fnv128a")
	fmt.Fprintln(os.Stderr, "lunh md4 md5 ripemd160 sha1 sha256 sha384 sha3-224 sha3-256 sha3-384 sha3-512 sha512 sha512-224 sha512-256 sum16 sum32 sum64 xor8")
	os.Exit(1)
}

func mustMakeHash(h hash.Hash, err error) hash.Hash {
	if err != nil {
		log.Fatal("error making hash: ", err)
	}
	return h
}

func makeHash(name string) hash.Hash {
	switch strings.ToLower(name) {
	case "blake2b-256":
		return mustMakeHash(blake2b.New256(nil))
	case "blake2b-384":
		return mustMakeHash(blake2b.New384(nil))
	case "blake2b-512":
		return mustMakeHash(blake2b.New512(nil))
	case "blake2s-256":
		return mustMakeHash(blake2s.New256(nil))
	case "ripemd160":
		return ripemd160.New()
	case "md4":
		return md4.New()
	case "md5":
		return md5.New()
	case "sha1":
		return sha1.New()
	case "sha256":
		return sha256.New()
	case "sha384":
		return sha512.New384()
	case "sha3-224":
		return sha3.New224()
	case "sha3-256":
		return sha3.New256()
	case "sha3-384":
		return sha3.New384()
	case "sha3-512":
		return sha3.New512()
	case "sha512":
		return sha512.New()
	case "sha512-224":
		return sha512.New512_224()
	case "sha512-256":
		return sha512.New512_256()
	case "crc32-ieee":
		return crc32.NewIEEE()
	case "crc64-iso":
		return crc64.New(crc64.MakeTable(crc64.ISO))
	case "crc64-ecma":
		return crc64.New(crc64.MakeTable(crc64.ECMA))
	case "adler32":
		return adler32.New()
	case "fnv32":
		return fnv.New32()
	case "fnv32a":
		return fnv.New32a()
	case "fnv64":
		return fnv.New64()
	case "fnv64a":
		return fnv.New64a()
	case "fnv128":
		return fnv.New128()
	case "fnv128a":
		return fnv.New128a()
	case "xor8":
		return new(xor8)
	case "fletch16":
		return &fletch16{}
	case "luhn":
		return new(luhn)
	case "sum16":
		return new(sum16)
	case "sum32":
		return new(sum32)
	case "sum64":
		return new(sum64)
	case "crc8":
		return new(crc8)
	case "crc16-ccitt":
		c := new(crc16ccitt)
		c.Reset()
		return c
	}
	return nil
}

func sum(r io.Reader, name string, h hash.Hash) error {
	_, err := io.Copy(h, r)
	if err != nil {
		return err
	}
	digest := h.Sum(nil)
	fmt.Printf("%x %s\n", digest, name)
	return nil
}

type fletch16 struct {
	sum1, sum2 uint16
}

func (d *fletch16) Reset()      { d.sum1, d.sum2 = 0, 0 }
func (fletch16) Size() int      { return 2 }
func (fletch16) BlockSize() int { return 1 }

func (d *fletch16) Write(b []byte) (int, error) {
	for i := range b {
		d.sum1 = (d.sum1 + uint16(b[i])) % 255
		d.sum2 = (d.sum2 + d.sum1) % 255
	}
	return len(b), nil
}

func (d *fletch16) Sum(b []byte) []byte {
	return append(b, byte(d.sum2), byte(d.sum1))
}

type xor8 uint8

func (d *xor8) Reset()      { *d = 0 }
func (xor8) Size() int      { return 1 }
func (xor8) BlockSize() int { return 1 }

func (d *xor8) Write(b []byte) (int, error) {
	for i := range b {
		*d ^= xor8(b[i])
	}
	return len(b), nil
}

func (d xor8) Sum(b []byte) []byte {
	return append(b, uint8(d))
}

type luhn uint64

func (d *luhn) Reset()      { *d = 0 }
func (luhn) Size() int      { return 8 }
func (luhn) BlockSize() int { return 1 }

func (d *luhn) Write(b []byte) (int, error) {
	tab := [...]uint64{0, 2, 4, 6, 8, 1, 3, 5, 7, 9}
	odd := len(b) & 1

	var sum uint64
	for i, c := range b {
		if c < '0' || c > '9' {
			c %= 10
		} else {
			c -= '0'
		}
		if i&1 == odd {
			sum += tab[c]
		} else {
			sum += uint64(c)
		}
	}
	*d = luhn(sum)

	return len(b), nil
}

func (d luhn) Sum(b []byte) []byte {
	return append(b, byte(d), byte(d>>8), byte(d>>16), byte(d>>24),
		byte(d>>32), byte(d>>40), byte(d>>48), byte(d>>56))
}

type sum16 uint16

func (d *sum16) Reset()      { *d = 0 }
func (sum16) Size() int      { return 2 }
func (sum16) BlockSize() int { return 1 }

func (d *sum16) Write(b []byte) (int, error) {
	for i := range b {
		*d += sum16(b[i])
	}
	return len(b), nil
}

func (d sum16) Sum(b []byte) []byte {
	return append(b, byte(d&0xff), byte(d>>8))
}

type sum32 uint32

func (d *sum32) Reset()      { *d = 0 }
func (sum32) Size() int      { return 4 }
func (sum32) BlockSize() int { return 1 }

func (d *sum32) Write(b []byte) (int, error) {
	for i := range b {
		*d += sum32(b[i])
	}
	return len(b), nil
}

func (d sum32) Sum(b []byte) []byte {
	return append(b, byte(d&0xff), byte(d>>8), byte(d>>16), byte(d>>24))
}

type sum64 uint64

func (d *sum64) Reset()      { *d = 0 }
func (sum64) Size() int      { return 8 }
func (sum64) BlockSize() int { return 1 }

func (d *sum64) Write(b []byte) (int, error) {
	for i := range b {
		*d += sum64(b[i])
	}
	return len(b), nil
}

func (d sum64) Sum(b []byte) []byte {
	return append(b, byte(d&0xff), byte(d>>8), byte(d>>16), byte(d>>24), byte(d>>32), byte(d>>40), byte(d>>48), byte(d>>56))
}

type crc8 uint8

func (c *crc8) Reset()      { *c = 0 }
func (crc8) Size() int      { return 1 }
func (crc8) BlockSize() int { return 1 }

func (c crc8) calc(b byte) crc8 {
	d := c ^ crc8(b)
	for i := 0; i < 8; i++ {
		if d&0x80 != 0 {
			d <<= 1
			d ^= 0x07
		} else {
			d <<= 1
		}
	}
	return d
}

func (c *crc8) Write(b []byte) (int, error) {
	for i := range b {
		*c = c.calc(b[i])
	}
	return len(b), nil
}

func (c crc8) Sum(b []byte) []byte {
	return append(b, byte(c))
}

type crc16ccitt uint16

func (c *crc16ccitt) Reset()      { *c = 0xffff }
func (crc16ccitt) Size() int      { return 2 }
func (crc16ccitt) BlockSize() int { return 1 }

func (c *crc16ccitt) Write(b []byte) (int, error) {
	for i := range b {
		x := *c>>8 ^ crc16ccitt(b[i])
		x ^= x >> 4
		*c = (*c << 8) ^ crc16ccitt((x<<12)^(x<<5)^x)
	}
	return len(b), nil
}

func (c crc16ccitt) Sum(b []byte) []byte {
	return append(b, uint8(c>>8), uint8(c))
}
