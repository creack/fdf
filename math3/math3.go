package math3

import (
	"math"
)

type Vec struct {
	X, Y, Z float64
}

// Scale the vector by the given factor.
func (v Vec) Scale(scale float64) Vec {
	return Vec{
		X: v.X * scale,
		Y: v.Y * scale,
		Z: v.Z * scale,
	}
}

// Translate the vector.
func (v Vec) Translate(offset Vec) Vec {
	return Vec{
		X: v.X + offset.X,
		Y: v.Y + offset.Y,
		Z: v.Z + offset.Z,
	}
}

// Rotate the vector.
func (v Vec) Rotate(angle Vec) Vec {
	v = v.MultiplyMatrix(GetRotationMatrix(angle.Z, AxisZ))
	v = v.MultiplyMatrix(GetRotationMatrix(angle.X, AxisX))
	v = v.MultiplyMatrix(GetRotationMatrix(angle.Y, AxisY))
	return v
}

func (v Vec) MultiplyMatrix(m Matrix) Vec {
	return Vec{
		X: v.X*m.i.X + v.Y*m.i.Y + v.Z*m.i.Z,
		Y: v.X*m.j.X + v.Y*m.j.Y + v.Z*m.j.Z,
		Z: v.X*m.k.X + v.Y*m.k.Y + v.Z*m.k.Z,
	}
}

// Matrix is a 3D matrix.
type Matrix struct {
	i, j, k Vec
}

func (m Matrix) Multiply(m2 Matrix) Matrix {
	var result [3][3]float64

	mA := [3][3]float64{
		{m.i.X, m.j.X, m.k.X},
		{m.i.Y, m.j.Y, m.k.Y},
		{m.i.Z, m.j.Z, m.k.Z},
	}

	mB := [3][3]float64{
		{m2.i.X, m2.j.X, m2.k.X},
		{m2.i.Y, m2.j.Y, m2.k.Y},
		{m2.i.Z, m2.j.Z, m2.k.Z},
	}

	// Multiplying matrices and storing result.
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			var total float64
			for k := 0; k < 3; k++ {
				// fmt.Printf("mA[%v][%v] * mB[%v][%v] = %v * %v = %v\n", i, k, k, j, mA[i][k], mB[k][j], mA[i][k]*mB[k][j])
				total = total + mA[i][k]*mB[k][j]
			}
			result[i][j] = total
		}
	}
	return Matrix{
		i: Vec{X: result[0][0], Y: result[0][1], Z: result[0][2]},
		j: Vec{X: result[1][0], Y: result[1][1], Z: result[1][2]},
		k: Vec{X: result[2][0], Y: result[2][1], Z: result[2][2]},
	}
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

func GetRotationMatrix(deg float64, axis Axis) Matrix {
	switch axis {
	case AxisX:
		return Matrix{
			i: Vec{
				X: 1,
				Y: 0,
				Z: 0,
			},
			j: Vec{
				X: 0,
				Y: math.Cos(deg),
				Z: math.Sin(deg),
			},
			k: Vec{
				X: 0,
				Y: -math.Sin(deg),
				Z: math.Cos(deg),
			},
		}
	case AxisY:
		return Matrix{
			i: Vec{
				X: math.Cos(deg),
				Y: 0,
				Z: -math.Sin(deg),
			},
			j: Vec{
				X: 0,
				Y: 1,
				Z: 0,
			},
			k: Vec{
				X: math.Sin(deg),
				Y: 0,
				Z: math.Cos(deg),
			},
		}
	case AxisZ:
		return Matrix{
			i: Vec{
				X: math.Cos(deg),
				Y: math.Sin(deg),
				Z: 0,
			},
			j: Vec{
				X: -math.Sin(deg),
				Y: math.Cos(deg),
				Z: 0,
			},
			k: Vec{
				X: 0,
				Y: 0,
				Z: 1,
			},
		}
	default:
		panic("invalid axis")
	}
}
