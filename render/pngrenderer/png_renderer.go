// Package pngrenderer provides a stdlib image/png impementation of the renderer.
package pngrenderer

import (
	"bytes"
	"fmt"
	"image/png"
	"log"
	"os"

	"go.creack.net/fdf/render"
)

type renderer struct {
	fileName string

	width  int
	height int
}

// New creates the renderer.
func New(fileName string, width, height int) render.Renderer {
	return &renderer{
		fileName: fileName,
		width:    width,
		height:   height,
	}
}

// Run implements the interface.
func (r *renderer) Run(g render.Engine) error {
	// Set the projection to isometric.
	render.Iso(g, r.width, r.height)

	// Render the wireframe.
	fdfImg := g.Draw()

	buf := bytes.NewBuffer(nil)
	if err := png.Encode(buf, fdfImg); err != nil {
		return fmt.Errorf("png.Encode: %w", err)
	}
	if err := os.WriteFile(r.fileName, buf.Bytes(), 0o600); err != nil {
		log.Fatalf("Write file: %s.", err)
	}

	return nil
}
