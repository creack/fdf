package main

import (
	"fmt"
	"image"
	"math"

	"fdf/render"
)

const (
	WIDTH  = screenWidth
	HEIGHT = screenHeight
)

type Fdf struct {
	render.StubEngine

	scale          float64
	offset         image.Point
	cameraRotation vec3

	depthScale float64

	Points [][]MapPoint
}

func (m *Fdf) IncScale(n float64) { m.scale += n }

func (m *Fdf) SetScale(n float64) { m.scale = n }

var (
	defaultZDeg float64 = 45
	defaultXDeg float64 = math.Atan(math.Sqrt2)
)

func (m *Fdf) cartesianToIsometric(vec vec3) vec3 {
	vec.z *= m.depthScale
	vec = multiplyVectorMatrix(vec, getRotationMatrix(m.cameraRotation.z, AxisZ))
	vec = multiplyVectorMatrix(vec, getRotationMatrix(m.cameraRotation.x, AxisX))
	vec = multiplyVectorMatrix(vec, getRotationMatrix(m.cameraRotation.y, AxisY))
	return vec
}

func (m *Fdf) Draw() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight))

	for j, line := range m.Points {
		for i, elem := range line {
			v := m.scaleOffset(elem.Vector())
			pv := image.Point{X: int(v.x), Y: int(v.y)}

			if i+1 < len(line) {
				v1 := m.scaleOffset(m.Points[j][i+1].Vector())
				pv1 := image.Point{X: int(v1.x), Y: int(v1.y)}
				drawLine(img, pv, pv1, v1.color)
			}
			if j+1 < len(m.Points) && i < len(m.Points[j+1]) {
				v1 := m.scaleOffset(m.Points[j+1][i].Vector())
				pv1 := image.Point{X: int(v1.x), Y: int(v1.y)}
				drawLine(img, pv, pv1, v1.color)
			}

		}
	}
	return img
}

func newFdf(mapData []byte) (*Fdf, error) {
	g := &Fdf{}

	m, err := parseMap(mapData)
	if err != nil {
		return nil, fmt.Errorf("parseMap: %w", err)
	}
	g.Points = m

	g.depthScale = 1
	g.cameraRotation.x = defaultXDeg
	g.cameraRotation.z = defaultZDeg

	// bounds := g.getProjectedBounds(1)
	// g.scale = getScale(bounds)
	// bounds = g.getProjectedBounds(g.scale)
	// g.offset = getOffset(WIDTH, HEIGHT, bounds)
	// fmt.Printf("%f - %d/%d\n", g.scale, g.offset.X, g.offset.Y)
	g.scale = 42
	g.offset.X = 100
	g.offset.Y = 100
	return g, nil
}

func (m *Fdf) toOffset(v vec3, offset image.Point) vec3 {
	nv := m.cartesianToIsometric(v)
	return vec3{
		x:     nv.x + float64(offset.X),
		y:     nv.y + float64(offset.Y),
		z:     nv.z,
		color: nv.color,
	}
}

func (m *Fdf) scaleOffset(v vec3) vec3 {
	return m.toOffset(v.Scale(m.scale), m.offset)
}

func getScale(border image.Rectangle) float64 {
	width := math.Abs(float64(border.Max.X - border.Min.X))
	height := math.Abs(float64(border.Max.Y - border.Min.Y))
	return math.Floor(math.Min((WIDTH-WIDTH/8.0)/width, (HEIGHT-HEIGHT/8.0)/height))
}

func getOffset(screenWidth, screenHeight int, bounds image.Rectangle) image.Point {
	width := math.Abs(float64(bounds.Max.X - bounds.Min.X))
	height := math.Abs(float64(bounds.Max.Y - bounds.Min.Y))
	offsetX := math.Round((float64(screenWidth) - width) / 2.0)
	offsetY := math.Round((float64(screenHeight) - height) / 2.0)
	if bounds.Min.X < 0 {
		offsetX += math.Abs(math.Round(float64(bounds.Min.X)))
	}
	if bounds.Min.Y < 0 {
		offsetY += math.Abs(math.Round(float64(bounds.Min.Y)))
	}
	return image.Point{X: int(offsetX), Y: int(offsetY)}
}
