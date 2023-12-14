package main

import (
	"embed"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io/fs"
	"math"
	"strings"

	"go.creack.net/fdf/math3"
	"go.creack.net/fdf/projection"
)

// Fdf reprensents the main engine to draw wireframes.
type Fdf struct {
	Points [][]MapPoint

	projection projection.Projection

	heightFactor float64

	mapData embed.FS
	mapName string
}

// NewFdf loads/parses the map and creates a fdf engine.
func NewFdf(mapData embed.FS, mapName string) (*Fdf, error) {
	g := &Fdf{
		projection:   projection.NewDirect(),
		heightFactor: 1,

		mapData: mapData,
	}
	if err := g.LoadMap(mapName); err != nil {
		return nil, fmt.Errorf("loadMap %q: %w", mapName, err)
	}
	return g, nil
}

// CurrentMapName returns the current map name.
func (m *Fdf) CurrentMapName() string { return m.mapName }

// ListMaps returns the list of available maps.
func (m *Fdf) ListMaps() []fs.DirEntry {
	out, err := m.mapData.ReadDir("maps")
	if err != nil { // Should never happen, but check just in case.
		panic(fmt.Errorf("open maps dir: %w", err))
	}
	return out
}

// LoadMap loads the given map name.
func (m *Fdf) LoadMap(mapName string) error {
	fmt.Println("Loading new map", mapName)
	buf, err := mapData.ReadFile(mapName)
	if err != nil {
		return fmt.Errorf("fs readfile: %w", err)
	}
	newMap, err := parseMap(buf)
	if err != nil {
		return fmt.Errorf("parseMap: %w", err)
	}
	m.Points = newMap
	m.mapName = strings.TrimPrefix(mapName, "maps/")
	return nil
}

// GetProjection accesses the value.
func (m *Fdf) GetProjection() projection.Projection { return m.projection }

// SetProjection applies the given projection and re-process the projected bounds.
func (m *Fdf) SetProjection(p projection.Projection) image.Rectangle {
	m.projection = p

	bounds := m.getProjectedBounds()

	// Adjust the offset to make sure we don't go in negative.
	var offset math3.Vec
	if bounds.Min.X <= 0 {
		offset.X = float64(-bounds.Min.X)
		bounds.Max.X += -bounds.Min.X
		bounds.Min.X = 0
	}
	if bounds.Min.Y <= 0 {
		offset.Y = float64(-bounds.Min.Y)
		bounds.Max.Y += -bounds.Min.Y
		bounds.Min.Y = 0
	}
	m.projection.SetOffset(offset)

	return bounds
}

// GetHeightFactor accesses the value.
func (m *Fdf) GetHeightFactor() float64 { return m.heightFactor }

// SetHeightFactor sets the value.
func (m *Fdf) SetHeightFactor(f float64) { m.heightFactor = f }

// Draw renders the image.
func (m *Fdf) Draw() image.Image {
	bounds := m.getProjectedBounds()
	img := image.NewRGBA(bounds)

	// Add black background.
	// Remove this to allow for transparency.
	draw.Draw(img, img.Bounds(), image.NewUniform(color.Black), image.Point{}, draw.Over)

	for y, line := range m.Points {
		for x, elem := range line {
			v := m.projection.Project(elem.Vec.ScaleZ(m.heightFactor))
			pv := image.Point{X: int(v.X), Y: int(v.Y)}

			if x+1 < len(line) {
				elem1 := m.Points[y][x+1]
				v1 := m.projection.Project(elem1.Vec.ScaleZ(m.heightFactor))
				pv1 := image.Point{X: int(v1.X), Y: int(v1.Y)}
				drawLine(img, pv, pv1, elem.color, elem1.color)
			}
			if y+1 < len(m.Points) && x < len(m.Points[y+1]) {
				elem1 := m.Points[y+1][x]
				v1 := m.projection.Project(elem1.Vec.ScaleZ(m.heightFactor))
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
	var bounds image.Rectangle

	for _, line := range m.Points {
		for _, elem := range line {
			point := m.projection.Project(elem.Vec)

			if math.Floor(point.X) < float64(bounds.Min.X) {
				bounds.Min.X = int(math.Floor(point.X))
			} else if math.Ceil(point.X) > float64(bounds.Max.X) {
				bounds.Max.X = int(math.Ceil(point.X))
			}
			if math.Floor(point.Y) < float64(bounds.Min.Y) {
				bounds.Min.Y = int(math.Floor(point.Y))
			} else if math.Ceil(point.Y) > float64(bounds.Max.Y) {
				bounds.Max.Y = int(math.Ceil(point.Y))
			}
		}
	}

	return bounds
}
