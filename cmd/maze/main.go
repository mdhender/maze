// Copyright (c) 2024 Michael D Henderson. All rights reserved.

// Package main implements a command line application to generate mazes using Wilson's algorithm
package main

import (
	"flag"
	"github.com/mdhender/maze"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {
	var testSeed int64
	flag.Int64Var(&testSeed, "seed", testSeed, "generate maze from seed")
	height := 125
	flag.IntVar(&height, "height", height, "height of maze (in cells)")
	width := 125
	flag.IntVar(&width, "width", width, "width of maze (in cells)")
	scale := 20
	flag.IntVar(&scale, "scale", scale, "width of cells in rendered maze")
	var pngFile, pngSolvedFile string
	flag.StringVar(&pngFile, "png", pngFile, "optional name of PNG image file to render")
	flag.StringVar(&pngSolvedFile, "png-solved", pngSolvedFile, "optional name of PNG image file with solution")
	var svgFile, svgSolvedFile string
	flag.StringVar(&svgFile, "svg", svgFile, "optional name of SVG image file to render")
	flag.StringVar(&svgSolvedFile, "svg-solved", svgSolvedFile, "optional name of SVG image file with solution")
	var txtFile string
	flag.StringVar(&txtFile, "text", txtFile, "optional name of text file to render")
	var version bool
	flag.BoolVar(&version, "version", version, "print version and exit")

	flag.Parse()

	if version {
		log.Println("maze: 1.0.0")
		return
	}

	// set seed only if we're testing changes
	if testSeed != 0 {
		log.Printf("maze: using seed %d\n", testSeed)
		rand.Seed(testSeed)
	}

	started := time.Now()
	rg, err := maze.RectangleMaze(height, width, false)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("maze: created %5d x %5d maze in %v\n", height, width, time.Now().Sub(started))

	if txtFile != "" {
		started = time.Now()
		w, err := os.OpenFile(txtFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		} else if err = rg.RenderText(w); err != nil {
			log.Fatal(err)
		} else if err = w.Close(); err != nil {
			log.Fatal(err)
		}
		log.Printf("maze: created %s in %v\n", txtFile, time.Now().Sub(started))
	}

	if pngFile != "" {
		started = time.Now()
		w, err := os.OpenFile(pngFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		} else if err = rg.RenderPNG(w, scale); err != nil {
			log.Fatal(err)
		} else if err = w.Close(); err != nil {
			log.Fatal(err)
		}
		log.Printf("maze: created %s in %v\n", pngFile, time.Now().Sub(started))
	}

	if svgFile != "" {
		started = time.Now()
		w, err := os.OpenFile(svgFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		} else if err = rg.RenderSVG(w, scale); err != nil {
			log.Fatal(err)
		} else if err = w.Close(); err != nil {
			log.Fatal(err)
		}
		log.Printf("maze: created %s in %v\n", svgFile, time.Now().Sub(started))
	}

	if pngSolvedFile != "" {
		rg.Solve()
		started = time.Now()
		w, err := os.OpenFile(pngSolvedFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		} else if err = rg.RenderPNG(w, scale); err != nil {
			log.Fatal(err)
		} else if err = w.Close(); err != nil {
			log.Fatal(err)
		}
		log.Printf("maze: created %s in %v\n", pngSolvedFile, time.Now().Sub(started))
	}

	if svgSolvedFile != "" {
		started = time.Now()
		w, err := os.OpenFile(svgSolvedFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		} else if err = rg.RenderSVG(w, scale); err != nil {
			log.Fatal(err)
		} else if err = w.Close(); err != nil {
			log.Fatal(err)
		}
		log.Printf("maze: created %s in %v\n", svgSolvedFile, time.Now().Sub(started))
	}

}
