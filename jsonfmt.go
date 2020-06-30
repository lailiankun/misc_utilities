package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

var (
	write = flag.Bool("w", false, "write to file")
	sep   = flag.String("s", "\t", "separator")
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("jsonfmt: ")
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() == 0 {
		err := format(os.Stdin, os.Stdout, *sep)
		ck(err)
		return
	}

	r, err := os.Open(flag.Arg(0))
	if err != nil {
		if *write {
			r = os.Stdin
		} else {
			log.Fatal(err)
		}
	}

	w := new(bytes.Buffer)
	err = format(r, w, *sep)
	ck(err)

	if !*write {
		os.Stdout.Write(w.Bytes())
	} else {
		err = ioutil.WriteFile(flag.Arg(0), w.Bytes(), 0644)
		ck(err)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: [options] <file>")
	flag.PrintDefaults()
	os.Exit(1)
}

func format(r io.Reader, w io.Writer, sep string) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	p := new(bytes.Buffer)
	err = json.Indent(p, b, "", sep)
	if err != nil {
		return err
	}

	_, err = w.Write(p.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func ck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
