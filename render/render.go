package render

import "image"

// Engine defines the possible methods to interract from the Renderer to the Engine.
type Engine interface {
	Draw() image.Image

	IncScale(n float64)
	SetScale(n float64)
}

// StubEngine implements the Renderer interface with no-op methods.
type StubEngine struct{}

// Draw implements the interface.
func (StubEngine) Draw() image.Image { return nil }

// IncScale implements the interface.
func (StubEngine) IncScale(float64) {}

// SetScale implements the interface.
func (StubEngine) SetScale(float64) {}
