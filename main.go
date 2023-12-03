package main

import (
	"embed"
	_ "embed"
	"fmt"
	_ "image/png"
	"log"

	"fdf/render/ebitenrenderer"
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
	g, err := NewFdf()
	if err != nil {
		log.Fatalf("NewGame: %s.", err)
	}

	// if err := pngrenderer.New("foo.png", 2050, 1100).Run(g); err != nil {
	// 	log.Fatal(err)
	// }

	if err := ebitenrenderer.New(1024, 1024).Run(g); err != nil {
		log.Fatal(err)
	}
}
