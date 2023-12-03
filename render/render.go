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
	Draw() image.Image
}

// // StubEngine implements the Renderer interface with no-op methods.
// type StubEngine struct{}

// // Draw implements the interface.
// func (StubEngine) Draw() image.Image { return nil }

// // IncScale implements the interface.
// func (StubEngine) IncScale(float64) {}

// // SetScale implements the interface.
// func (StubEngine) SetScale(float64) {}
