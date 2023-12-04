package render

import (
	"image"

	"github.com/creack/fdf/projection"
)

// Renderer defines what a renderer can do.
type Renderer interface {
	Run(Engine) error
}

// Engine defines the possible methods to interract from the Renderer to the Engine.
type Engine interface {
	SetProjection(projection.Projection) image.Rectangle
	GetProjection() projection.Projection

	GetHeightFactor() float64
	SetHeightFactor(float64)

	Draw() image.Image
}

// Iso is a helper function to initialize an isomorphic projection.
func Iso(fdf Engine, screenWidth, screenHeight int) image.Rectangle {
	// Initialize a projection with scale 1 to get base boundaries.
	p := projection.NewIsomorphic(1)

	// Set the projection on the engine, returns the initial bounds.
	bounds := fdf.SetProjection(p)

	// Compute the scale from the initial bounds.
	scale := projection.GetScale(screenHeight, screenWidth, bounds)

	// Create the final projection with the new scale.
	p = projection.NewIsomorphic(scale)

	// Set the projection on the engine, returns the final bounds.
	bounds = fdf.SetProjection(p)

	return bounds
}
