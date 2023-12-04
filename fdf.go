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

	heightFactor float64
}

func NewFdf() (*Fdf, error) {
	buf, err := mapData.ReadFile("maps/42.fdf")
	if err != nil {
		return nil, fmt.Errorf("fs readfile: %w", err)
	}

	g := &Fdf{
		projection:   projection.NewDirect(),
		heightFactor: 1,
	}

	m, err := parseMap(buf)
	if err != nil {
		return nil, fmt.Errorf("parseMap: %w", err)
	}
	g.Points = m

	return g, nil
}

func (m *Fdf) GetProjection() projection.Projection { return m.projection }

func (m *Fdf) SetProjection(p projection.Projection) image.Rectangle {
	m.projection = p
	// Process the new bounds and return them.
	m.bounds = m.getProjectedBounds()

	var offset math3.Vec
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

func (m *Fdf) GetHeightFactor() float64  { return m.heightFactor }
func (m *Fdf) SetHeightFactor(f float64) { m.heightFactor = f }

func (m *Fdf) Draw() image.Image {
	if m.bounds == (image.Rectangle{}) {
		m.bounds = m.getProjectedBounds()
	}
	img := image.NewRGBA(m.bounds)

	for j, line := range m.Points {
		for i, elem := range line {
			elem.Z *= m.heightFactor // Note: elem is a copy. Safe to mutate.

			v := m.projection.Project(elem.Vec)
			pv := image.Point{X: int(v.X), Y: int(v.Y)}

			if i+1 < len(line) {
				elem1 := m.Points[j][i+1]
				elem1.Z *= m.heightFactor
				v1 := m.projection.Project(elem1.Vec)
				pv1 := image.Point{X: int(v1.X), Y: int(v1.Y)}
				drawLine(img, pv, pv1, elem.color, elem1.color)
			}
			if j+1 < len(m.Points) && i < len(m.Points[j+1]) {
				elem1 := m.Points[j+1][i]
				elem1.Z *= m.heightFactor
				v1 := m.projection.Project(elem1.Vec)
				pv1 := image.Point{X: int(v1.X), Y: int(v1.Y)}
				drawLine(img, pv, pv1, elem.color, elem1.color)
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
			point := m.projection.Project(elem.Vec)

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
