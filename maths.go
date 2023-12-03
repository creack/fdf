package main

import (
	"image/color"
	"math"
)

type vec3 struct {
	x, y, z float64
	color   color.Color
}

// Scale the vector by the given factor.
func (v vec3) Scale(scale float64) vec3 {
	return vec3{
		x:     v.x * scale,
		y:     v.y * scale,
		z:     v.z * scale,
		color: v.color,
	}
}

// matrix3 is a 3D matrix.
type matrix3 struct {
	i, j, k vec3
}

func multiplyVectorMatrix(v vec3, m matrix3) vec3 {
	lv := [3]float64{v.x, v.y, v.z}
	v.x = lv[0]*m.i.x + lv[1]*m.i.y + lv[2]*m.i.z
	v.y = lv[0]*m.j.x + lv[1]*m.j.y + lv[2]*m.j.z
	v.z = lv[0]*m.k.x + lv[1]*m.k.y + lv[2]*m.k.z
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

func getRotationMatrix(deg float64, axis Axis) matrix3 {
	switch axis {
	case AxisX:
		return matrix3{
			i: vec3{
				x: 1,
				y: 0,
				z: 0,
			},
			j: vec3{
				x: 0,
				y: math.Cos(deg),
				z: -math.Sin(deg),
			},
			k: vec3{
				x: 0,
				y: math.Sin(deg),
				z: math.Cos(deg),
			},
		}
	case AxisY:
		return matrix3{
			i: vec3{
				x: math.Cos(deg),
				y: 0,
				z: -math.Sin(deg),
			},
			j: vec3{
				x: 0,
				y: 1,
				z: 0,
			},
			k: vec3{
				x: math.Sin(deg),
				y: 0,
				z: math.Cos(deg),
			},
		}
	case AxisZ:
		return matrix3{
			i: vec3{
				x: math.Cos(rad(deg)),
				y: -math.Sin(rad(deg)),
				z: 0,
			},
			j: vec3{
				x: math.Sin(rad(deg)),
				y: math.Cos(rad(deg)),
				z: 0,
			},
			k: vec3{
				x: 0,
				y: 0,
				z: 1,
			},
		}
	default:
		return matrix3{
			i: vec3{
				x: 1,
				y: 0,
				z: 0,
			},
			j: vec3{
				x: 0,
				y: 1,
				z: 0,
			},
			k: vec3{
				x: 0,
				y: 0,
				z: 1,
			},
		}

	}
}
