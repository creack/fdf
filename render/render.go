package render

import (
	"image"

	"fdf/projection"
)

type Renderer interface {
	Run(Engine) error
}

// Engine defines the possible methods to interract from the Renderer to the Engine.
type Engine interface {
	SetProjection(projection.Projection) image.Rectangle
	GetProjection() projection.Projection
	Draw() image.Image
}

func Iso(fdf Engine, screenWidth, screenHeight int) image.Rectangle {
	return fdf.SetProjection(
		projection.NewIsomorphic(
			projection.GetScale(
				screenHeight,
				screenWidth,
				fdf.SetProjection(
					projection.NewIsomorphic(1),
				),
			),
		),
	)
}

// // StubEngine implements the Renderer interface with no-op methods.
// type StubEngine struct{}

// // Draw implements the interface.
// func (StubEngine) Draw() image.Image { return nil }

// // IncScale implements the interface.
// func (StubEngine) IncScale(float64) {}

// // SetScale implements the interface.
// func (StubEngine) SetScale(float64) {}
