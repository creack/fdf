// Package main is the entrypoint.
package main

import (
	"embed"
	"flag"
	"log"
	"runtime"

	"go.creack.net/fdf/render/ebitenrenderer"
	"go.creack.net/fdf/render/pngrenderer"
)

//go:embed maps/*.fdf
var mapData embed.FS

func main() {
	var renderer, filePath, source string
	flag.StringVar(&renderer, "r", "ebitengine", "Renderer: 'png' or 'ebitengine'. Always 'ebitengine' for WASM.")
	flag.StringVar(&filePath, "f", "./fdf.png", "Only for 'png' renderer: path where to create the image.")
	flag.StringVar(&source, "s", "maps/42.fdf", "Source .fdf map file.")
	flag.Parse()

	if runtime.GOOS == "js" {
		renderer = "ebitengine"
	}

	g, err := NewFdf(source)
	if err != nil {
		log.Fatalf("NewFdf: %s.", err)
	}

	switch renderer {
	case "png":
		if err := pngrenderer.New(filePath, 2050, 1100).Run(g); err != nil {
			log.Fatal(err)
		}
	case "ebitengine":
		println("Starting ebitengine.")
		if err := ebitenrenderer.New(300, 300).Run(g); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("Invalid renderer %q.", renderer)
	}
}
