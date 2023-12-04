// Package main is the entrypoint.
package main

import (
	"embed"
	"log"
	"runtime"

	"go.creack.net/fdf/render/ebitenrenderer"
	"go.creack.net/fdf/render/pngrenderer"
)

//go:embed maps/*.fdf
var mapData embed.FS

func main() {
	println("start")

	g, err := NewFdf()
	if err != nil {
		log.Fatalf("NewFdf: %s.", err)
	}

	if runtime.GOOS != "js" {
		if err := pngrenderer.New("docs/42.png", 2050, 1100).Run(g); err != nil {
			log.Fatal(err)
		}
	}

	if err := ebitenrenderer.New(300, 300).Run(g); err != nil {
		log.Fatal(err)
	}
}
