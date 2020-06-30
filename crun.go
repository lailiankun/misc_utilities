package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var (
	buildOnly  = flag.Bool("b", false, "build only")
	outputFile = flag.String("o", "", "output file")
	format     = flag.Bool("f", false, "format source file")
	pkgcfg     = flag.String("p", "", "use pkg-config for package flags")
	exflags    = flag.String("x", "", "pass extra flags")
	verbose    = flag.Bool("v", false, "be verbose")
	isCpp      = false
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("crun: ")
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() < 1 {
		usage()
	}

	rand.Seed(time.Now().UnixNano())
	files, args := splitArgs(flag.Args())
	if len(files) == 0 {
		log.Fatal("no C source file specified")
	}
	crun(files, args)
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: [options] file ... [args]")
	flag.PrintDefaults()
	os.Exit(2)
}

func splitArgs(args []string) ([]string, []string) {
	for i := range args {
		arg := strings.ToLower(args[i])
		if strings.HasSuffix(arg, ".cpp") || strings.HasSuffix(arg, ".cc") {
			isCpp = true
		}
		if !strings.HasSuffix(arg, ".c") && !strings.HasSuffix(arg, ".cpp") && !strings.HasSuffix(arg, ".cc") {
			return args[:i], args[i:]
		}
	}
	return args, nil
}

func crun(files, args []string) {
	cc := os.Getenv("CC")
	if cc == "" {
		cc = "cc"
	}
	if isCpp {
		cc = os.Getenv("CPP")
		if cc == "" {
			cc = "c++"
		}
	}
	sharedFlags := "-Wall -Wextra -pedantic -ggdb -g3 -lm -lpthread"
	if runtime.GOOS == "linux" && (runtime.GOARCH == "amd64" || runtime.GOARCH == "386") {
		sharedFlags += " -fsanitize=address"
	}
	if *pkgcfg != "" {
		pw := new(bytes.Buffer)
		cmd := exec.Command("pkg-config", "--cflags", "--libs", *pkgcfg)
		cmd.Stderr = pw
		cmd.Stdout = pw
		err := cmd.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", pw.Bytes())
			log.Fatal(err)
		}
		sharedFlags += " " + strings.TrimSpace(pw.String())
	}

	cflags := os.Getenv("CFLAGS")
	if cflags == "" {
		cflags = "-std=c11 " + sharedFlags
	}
	if isCpp {
		cflags = os.Getenv("CPPFLAGS")
		if cflags == "" {
			cflags = "-std=c++14 " + sharedFlags
		}
	}

	binName := tempName()
	wd, _ := os.Getwd()

	if *format {
		var fargs []string
		fargs = append(fargs, "-i")
		fargs = append(fargs, files...)
		cmd := exec.Command("clang-format", fargs...)
		cmd.Dir = wd
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}

	if *outputFile != "" {
		binName = *outputFile
	}

	var cargs []string
	cargs = append(cargs, "-o")
	cargs = append(cargs, binName)
	cargs = append(cargs, files...)
	cargs = append(cargs, strings.Split(cflags, " ")...)
	if *exflags != "" {
		cargs = append(cargs, *exflags)
	}
	if *verbose {
		fmt.Printf("%s %v\n", cc, cargs)
	}

	cmd := exec.Command(cc, cargs...)
	cmd.Dir = wd
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	if *buildOnly {
		return
	}

	donech := make(chan error)
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch)

	exePath := filepath.Join(wd, binName)
	exe := filepath.Join(exePath)
	cmd = exec.Command(exe, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	go func() {
		donech <- cmd.Run()
	}()

	select {
	case <-sigch:
	case err = <-donech:
	}
	os.Remove(exePath)
	if err != nil {
		log.Fatal(err)
	}
}

func tempName() string {
	return fmt.Sprintf("crun-%d", rand.Int())
}
