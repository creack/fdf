package math3

import (
	"image/color"
	"math"
)

type Vec3 struct {
	X, Y, Z float64
	Color   color.Color
}

// Scale the vector by the given factor.
func (v Vec3) Scale(scale float64) Vec3 {
	return Vec3{
		X:     v.X * scale,
		Y:     v.Y * scale,
		Z:     v.Z * scale,
		Color: v.Color,
	}
}

// matrix3 is a 3D matrix.
type matrix3 struct {
	i, j, k Vec3
}

// MultiplyVectorMatrix performs v * m.
func MultiplyVectorMatrix(v Vec3, m matrix3) Vec3 {
	lv := [3]float64{v.X, v.Y, v.Z}
	v.X = lv[0]*m.i.X + lv[1]*m.i.Y + lv[2]*m.i.Z
	v.Y = lv[0]*m.j.X + lv[1]*m.j.Y + lv[2]*m.j.Z
	v.Z = lv[0]*m.k.X + lv[1]*m.k.Y + lv[2]*m.k.Z
	return v
}

func rad(deg float64) float64 {
	const radFactor = math.Pi / 180
	return deg * radFactor
}

// Axis enum type.
type Axis byte

// Axis enum values.
const (
	AxisNone Axis = iota
	AxisX
	AxisY
	AxisZ
)

// GetRotationMatrix returns a matrix that can be used to rotate along the given axis by the given degs.
func GetRotationMatrix(deg float64, axis Axis) matrix3 {
	switch axis {
	case AxisX:
		return matrix3{
			i: Vec3{
				X: 1,
				Y: 0,
				Z: 0,
			},
			j: Vec3{
				X: 0,
				Y: math.Cos(deg),
				Z: -math.Sin(deg),
			},
			k: Vec3{
				X: 0,
				Y: math.Sin(deg),
				Z: math.Cos(deg),
			},
		}
	case AxisY:
		return matrix3{
			i: Vec3{
				X: math.Cos(deg),
				Y: 0,
				Z: -math.Sin(deg),
			},
			j: Vec3{
				X: 0,
				Y: 1,
				Z: 0,
			},
			k: Vec3{
				X: math.Sin(deg),
				Y: 0,
				Z: math.Cos(deg),
			},
		}
	case AxisZ:
		return matrix3{
			i: Vec3{
				X: math.Cos(rad(deg)),
				Y: -math.Sin(rad(deg)),
				Z: 0,
			},
			j: Vec3{
				X: math.Sin(rad(deg)),
				Y: math.Cos(rad(deg)),
				Z: 0,
			},
			k: Vec3{
				X: 0,
				Y: 0,
				Z: 1,
			},
		}
	default:
		return matrix3{
			i: Vec3{
				X: 1,
				Y: 0,
				Z: 0,
			},
			j: Vec3{
				X: 0,
				Y: 1,
				Z: 0,
			},
			k: Vec3{
				X: 0,
				Y: 0,
				Z: 1,
			},
		}

	}
}
