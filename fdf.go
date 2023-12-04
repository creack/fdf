package main

import (
	"fmt"
	"image"
	"math"

	"fdf/math3"
	"fdf/projection"
)

type Fdf struct {
	Points [][]MapPoint

	bounds     image.Rectangle
	projection projection.Projection
}

func newFdf(mapData []byte) (*Fdf, error) {
	g := &Fdf{
		projection: projection.NewDirect(),
	}

	m, err := parseMap(mapData)
	if err != nil {
		return nil, fmt.Errorf("parseMap: %w", err)
	}
	g.Points = m

	// g.depthScale = 1
	// g.cameraRotation.X = defaultXDeg
	// g.cameraRotation.Z = defaultZDeg

	// bounds := g.getProjectedBounds(1)
	// g.scale = getScale(bounds)
	// bounds = g.getProjectedBounds(g.scale)
	// g.offset = getOffset(WIDTH, HEIGHT, bounds)
	// fmt.Printf("%f - %d/%d\n", g.scale, g.offset.X, g.offset.Y)
	// g.scale = 42
	// g.offset.X = 100
	// g.offset.Y = 100
	// g.projection = projection.NewIsomorphic(int(g.scale), g.offset, g.cameraRotation)
	return g, nil
}

func (m *Fdf) GetProjection() projection.Projection { return m.projection }

func (m *Fdf) SetProjection(p projection.Projection) image.Rectangle {
	m.projection = p
	// Process the new bounds and return them.
	m.bounds = m.getProjectedBounds()

	var offset math3.Vec3
	if m.bounds.Min.X <= 0 {
		offset.X = float64(-m.bounds.Min.X)
		m.bounds.Max.X += -m.bounds.Min.X
		m.bounds.Min.X = 0
	}
	if m.bounds.Min.Y <= 0 {
		offset.Y = float64(-m.bounds.Min.Y)
		m.bounds.Max.Y += -m.bounds.Min.Y
		m.bounds.Min.Y = 0
	}
	m.projection.SetOffset(offset)

	return m.bounds
}

func (m *Fdf) Draw() image.Image {
	m.bounds = m.getProjectedBounds()
	img := image.NewRGBA(m.bounds)

	// Draw a transparent blue overlay on the whole fdf image.
	// draw.Draw(img, img.Bounds(), image.NewUniform(color.RGBA{A: 100, B: 200, G: 50}), image.Point{}, draw.Over)

	for j, line := range m.Points {
		for i, elem := range line {
			v := m.projection.Project(elem.Vector())
			pv := image.Point{X: int(v.X), Y: int(v.Y)}

			if i+1 < len(line) {
				v1 := m.projection.Project(m.Points[j][i+1].Vector())
				pv1 := image.Point{X: int(v1.X), Y: int(v1.Y)}
				drawLine(img, pv, pv1, v1.Color)
			}
			if j+1 < len(m.Points) && i < len(m.Points[j+1]) {
				v1 := m.projection.Project(m.Points[j+1][i].Vector())
				pv1 := image.Point{X: int(v1.X), Y: int(v1.Y)}
				drawLine(img, pv, pv1, v1.Color)
			}

		}
	}

	return img
}

// getProjectedBounds projects and scales each points of the map
// and returns the smallest boundaries fitting everything.
//
// Going over each point as the projection of any given point can result
// in a bigger viewport.
func (m *Fdf) getProjectedBounds() image.Rectangle {
	var border image.Rectangle

	for _, line := range m.Points {
		for _, elem := range line {
			point := m.projection.Project(elem.Vector())

			if math.Floor(point.X) < float64(border.Min.X) {
				border.Min.X = int(math.Floor(point.X))
			} else if math.Ceil(point.X) > float64(border.Max.X) {
				border.Max.X = int(math.Ceil(point.X))
			}
			if math.Floor(point.Y) < float64(border.Min.Y) {
				border.Min.Y = int(math.Floor(point.Y))
			} else if math.Ceil(point.Y) > float64(border.Max.Y) {
				border.Max.Y = int(math.Ceil(point.Y))
			}

		}
	}

	return border
}
