// Package projection provides the projection interface and implementations.
package projection

import (
	"image"
	"math"

	"go.creack.net/fdf/math3"
)

const (
	defaultXDeg = 1.
	defaultZDeg = 0.5
)

//nolint:gochecknoglobals // Expected "readonly" globals.
var (
	defaultCameraRotation = math3.Vec{
		X: defaultXDeg,
		Z: defaultZDeg,
	}
)

// Projection defines how to project the given 3d point.
type Projection interface {
	Project(math3.Vec) math3.Vec

	SetOffset(math3.Vec)

	GetScale() float64
	SetScale(float64)

	GetAngle() math3.Vec
	SetAngle(math3.Vec)
}

type direct struct {
	offset math3.Vec
	scale  float64
}

// NewDirect creates the default projection without direct changes.
func NewDirect() Projection { return &direct{scale: 1} }

func (d direct) Project(vec math3.Vec) math3.Vec {
	return math3.Vec{
		X: vec.X + d.offset.X,
		Y: vec.Y + d.offset.Y,
		Z: vec.Z + d.offset.Z,
	}
}

// SetOffset sets the internal offset.
func (d *direct) SetOffset(offset math3.Vec) {
	d.offset = offset
}

// GetScale accesses the value.
func (d direct) GetScale() float64 { return d.scale }

// SetScale sets the value.
func (d *direct) SetScale(s float64) { d.scale = s }

// GetAngle accesses the value.
func (direct) GetAngle() math3.Vec { return math3.Vec{} }

// SetAngle sets the value.
func (*direct) SetAngle(math3.Vec) {}

// isomorphic projection.
type isomorphic struct {
	scale int

	offset math3.Vec

	cameraRotation math3.Vec
}

// NewIsomorphic creates an isomorphic projection with the given scale.
func NewIsomorphic(scale int) Projection {
	return &isomorphic{
		scale:          scale,
		cameraRotation: defaultCameraRotation,
	}
}

// SetScale sets the value.
func (i *isomorphic) SetScale(s float64) {
	// Clamp.
	if s < 2 {
		s = 2
	} else if s > 100 {
		s = 100
	}
	i.scale = int(s)
}

// GetScale acceses the value.
func (i isomorphic) GetScale() float64 { return float64(i.scale) }

// GetAngle accesses the value.
func (i isomorphic) GetAngle() math3.Vec { return i.cameraRotation }

// SetAngle sets the value.
func (i *isomorphic) SetAngle(ang math3.Vec) { i.cameraRotation = ang }

// SetOffset sets the value.
func (i *isomorphic) SetOffset(offset math3.Vec) {
	i.offset = offset
}

// Project applies the proejection to the given vector.
func (i isomorphic) Project(vec math3.Vec) math3.Vec {
	// First scale the vector.
	vec = vec.ScaleAll(float64(i.scale))

	// Then rotate.
	vec = vec.Rotate(i.cameraRotation)

	// Then translate.
	return vec.Translate(i.offset)
}

// GetScale returns the scale factor to fit bounds in the given screen size.
func GetScale(screenWidth, screenHeight int, bounds image.Rectangle) int {
	width := bounds.Dx()
	height := bounds.Dy()

	return min((screenWidth-screenWidth/8.0)/width, (screenHeight-screenHeight/8.0)/height)
}

// GetOffsetCenter returns the offset to move `bounds` to the center of the screen.
func GetOffsetCenter(screenWidth, screenHeight int, bounds image.Rectangle) image.Point {
	width := bounds.Dx()
	height := bounds.Dy()

	offsetX := int(math.Round((float64(screenWidth - width)) / 2.0))
	offsetY := int(math.Round((float64(screenHeight - height)) / 2.0))
	if bounds.Min.X < 0 {
		offsetX -= bounds.Min.X
	}
	if bounds.Min.Y < 0 {
		offsetY -= bounds.Min.Y
	}
	return image.Pt(offsetX, offsetY)
}
