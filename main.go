package main

import (
	"embed"
	"fmt"
	"log"
	"runtime"

	"fdf/render/ebitenrenderer"
	"fdf/render/pngrenderer"
)

//go:embed maps/*.fdf
var mapData embed.FS

func NewFdf() (*Fdf, error) {
	buf, err := mapData.ReadFile("maps/42.fdf")
	if err != nil {
		return nil, fmt.Errorf("fs readfile: %w", err)
	}

	m, err := newFdf(buf)
	if err != nil {
		return nil, fmt.Errorf("loadMap: %w", err)
	}
	return m, nil
}

func main() {
	println("start")

	g, err := NewFdf()
	if err != nil {
		log.Fatalf("NewGame: %s.", err)
	}

	if runtime.GOOS != "js" {
		if err := pngrenderer.New("foo.png", 2050, 1100).Run(g); err != nil {
			log.Fatal(err)
		}
	}

	if err := ebitenrenderer.New(300, 300).Run(g); err != nil {
		log.Fatal(err)
	}
}
