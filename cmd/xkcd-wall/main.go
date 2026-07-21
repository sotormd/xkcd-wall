package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"
	"xkcd-wall/xkcd"
)

func main() {
	defaultComicType := "today"
	comicTypeArg := flag.String("t", defaultComicType, "today, random, or <number>")

	defaultDimensions := "1920x1200"
	dimensionsArg := flag.String("d", defaultDimensions, "dimensions, eg. 1920x1080")

	defaultBackground := "2e3440"
	backgroundArg := flag.String("b", defaultBackground, "background, eg. ffffff")

	defaultForeground := "d8dee9"
	foregroundArg := flag.String("f", defaultForeground, "foreground, eg. 000000")

	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Error: expected output path")
		os.Exit(1)
	}

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not get cache dir: %v\n", err)
		os.Exit(1)
	}

	comicType := *comicTypeArg
	dimensions := *dimensionsArg
	bg := *backgroundArg
	fg := *foregroundArg

	path, err := xkcd.Get(comicType, cacheDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not fetch comic: %v\n", err)
		os.Exit(1)
	}

	colored, err := xkcd.Colorize(path, bg, fg, cacheDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not colorize comic: %v\n", err)
		os.Exit(1)
	}

	background, err := xkcd.MakeBackground(dimensions, bg, cacheDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not create background: %v\n", err)
		os.Exit(1)
	}

	final, err := xkcd.CompositeCenter(colored, background, cacheDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not create final image: %v\n", err)
		os.Exit(1)
	}

	if err := xkcd.CopyFile(final, args[0]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not copy image to target: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(args[0])
}
