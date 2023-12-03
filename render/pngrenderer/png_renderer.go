package pngrenderer

import (
	"bytes"
	"image/png"
	"log"
	"os"

	"fdf/projection"
	"fdf/render"
)

type renderer struct {
	fileName string

	width  int
	height int
}

func New(fileName string, width, height int) render.Renderer {
	return &renderer{
		fileName: fileName,
		width:    width,
		height:   height,
	}
}

func (r *renderer) Run(g render.Engine) error {
	g.SetProjection(projection.NewIsomorphic2(projection.GetScale(r.width, r.height, g.SetProjection(projection.NewIsomorphic2(1)))))

	fdfImg := g.Draw()

	buf := bytes.NewBuffer(nil)
	if err := png.Encode(buf, fdfImg); err != nil {
		log.Fatalf("Encode png: %s.", err)
	}
	if err := os.WriteFile(r.fileName, buf.Bytes(), 0o644); err != nil {
		log.Fatalf("Write file: %s.", err)
	}

	return nil
}
