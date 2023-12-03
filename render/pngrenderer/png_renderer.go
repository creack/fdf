package pngrenderer

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math"
	"os"

	"fdf/math3"
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
	canvas := image.NewRGBA(image.Rect(0, 0, r.width, r.height))

	g.SetProjection(projection.NewIsomorphic(1, image.Point{}, math3.Vec3{
		X: math.Atan(math.Sqrt2),
		Z: 45,
	}))
	bounds := g.Draw().Bounds()
	fmt.Printf("Initial bounds: %v\n", bounds)
	s := projection.GetScale(canvas.Bounds().Dx(), canvas.Bounds().Dy(), bounds)

	g.SetProjection(projection.NewIsomorphic(s, image.Point{}, math3.Vec3{
		X: math.Atan(math.Sqrt2),
		Z: 45,
	}))
	fmt.Printf("s: %v\n", s)
	fmt.Println(">>", g.Draw().Bounds())
	fmt.Println(">>1", canvas.Bounds().Dx(), canvas.Bounds().Dy())

	fdfImg := g.Draw()
	_ = fdfImg
	offset := projection.GetOffsetCenter(canvas.Bounds().Dx(), canvas.Bounds().Dy(), bounds)
	_ = offset

	// Draw the fdf in the canvas.
	draw.Draw(canvas, canvas.Bounds(), fdfImg, offset.Mul(-1), draw.Src)

	center := offset
	for i := 0; i < r.height; i++ {
		canvas.Set(center.X, i, color.RGBA{A: 255, B: 255})
	}
	for i := 0; i < r.width; i++ {
		canvas.Set(i, center.Y, color.RGBA{A: 255, B: 255})
	}

	for i := -5; i < 5; i++ {
		for j := -5; j < 5; j++ {
			canvas.Set(center.X+i, center.Y+j, color.RGBA{A: 255, R: 255})
		}
	}

	buf := bytes.NewBuffer(nil)
	if err := png.Encode(buf, canvas); err != nil {
		log.Fatalf("Encode png: %s.", err)
	}
	if err := os.WriteFile(r.fileName, buf.Bytes(), 0o644); err != nil {
		log.Fatalf("Write file: %s.", err)
	}

	return nil
}
