package projection

import (
	"image"
	"math"

	"fdf/math3"
)

//nolint:gochecknoglobals // Expected "readonly" globals.
var (
	defaultXDeg float64 = math.Atan(math.Sqrt2)
	defaultZDeg float64 = 45

	defaultCameraRotation = math3.Vec3{
		X: defaultXDeg,
		Z: defaultZDeg,
	}
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
	scale int

	cameraRotation math3.Vec3
}

func NewIsomorphic2(scale int) Projection {
	return &isomorphic{
		scale:          scale,
		cameraRotation: defaultCameraRotation,
	}
}

func (i isomorphic) Project(vec math3.Vec3) math3.Vec3 {
	// First scale the vector.
	vec = vec.Scale(float64(i.scale))

	// Then rotate.
	vec = math3.MultiplyVectorMatrix(vec, math3.GetRotationMatrix(i.cameraRotation.Z, math3.AxisZ))
	vec = math3.MultiplyVectorMatrix(vec, math3.GetRotationMatrix(i.cameraRotation.X, math3.AxisX))
	vec = math3.MultiplyVectorMatrix(vec, math3.GetRotationMatrix(i.cameraRotation.Y, math3.AxisY))

	return vec
}

func GetScale(screenWidth, screenHeight int, bounds image.Rectangle) int {
	width := bounds.Dx()
	height := bounds.Dy()

	return min((screenWidth-screenWidth/8.0)/width, (screenHeight-screenHeight/8.0)/height)
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
