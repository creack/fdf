package projection

import (
	"fmt"
	"image"
	"math"

	"fdf/math3"
)

// Projection defines how to project the given 3d point.
type Projection interface {
	Project(math3.Vec3) math3.Vec3
}

type direct struct{}

func NewDirect() Projection { return direct{} }

func (direct) Project(vec math3.Vec3) math3.Vec3 { return vec }

// isomorphic projection.
type isomorphic struct {
	scale  int
	offset image.Point

	cameraRotation math3.Vec3
}

// NewIsomorphic creates the projection.
func NewIsomorphic(scale int, offset image.Point, cameraRotation math3.Vec3) Projection {
	return &isomorphic{
		scale:          scale,
		offset:         offset,
		cameraRotation: cameraRotation,
	}
}

func (i isomorphic) Project(vec math3.Vec3) math3.Vec3 {
	// First scale the vector.
	vec = vec.Scale(float64(i.scale))

	// Then rotate.
	vec = math3.MultiplyVectorMatrix(vec, math3.GetRotationMatrix(i.cameraRotation.Z, math3.AxisZ))
	vec = math3.MultiplyVectorMatrix(vec, math3.GetRotationMatrix(i.cameraRotation.X, math3.AxisX))
	vec = math3.MultiplyVectorMatrix(vec, math3.GetRotationMatrix(i.cameraRotation.Y, math3.AxisY))

	// Then translate.
	vec = math3.Vec3{
		X:     vec.X + float64(i.offset.X),
		Y:     vec.Y + float64(i.offset.Y),
		Z:     vec.Z,
		Color: vec.Color,
	}

	return vec
}

func GetScale(screenWidth, screenHeight int, bounds image.Rectangle) int {
	width := bounds.Dx()
	height := bounds.Dy()

	out := min((screenWidth-screenWidth/8.0)/width, (screenHeight-screenHeight/8.0)/height)
	fmt.Printf("Screen H: %d, H: %d, Scaled H: %d\n", screenHeight, height, out*height)
	return out
}

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
	return image.Point{X: int(offsetX), Y: int(offsetY)}
}
