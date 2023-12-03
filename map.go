package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"strconv"
	"strings"
)

// TODO: Consider moving this to the renderer.
//
//nolint:gochecknoglobals // Expected "readonly" global.
var defaultColor = color.White

func parseMap(mapData []byte) ([][]MapPoint, error) {
	// Start by cleaning up the input, removing blank lines and dup spaces.
	var grid [][]string
	for _, line := range strings.Split(string(mapData), "\n") {
		if line == "" {
			continue
		}
		var gridLine []string
		for _, elem := range strings.Split(line, " ") {
			if elem == "" {
				continue
			}
			gridLine = append(gridLine, elem)
		}
		grid = append(grid, gridLine)
	}

	// Then for each point, parse the height and optional color.
	var m [][]MapPoint
	for y, line := range grid {
		var points []MapPoint
		for x, elem := range line {
			parts := strings.Split(elem, ",")

			h, err := strconv.Atoi(parts[0])
			if err != nil {
				return nil, fmt.Errorf("invalid height %q for %d/%d: %w", parts[0], y, x, err)
			}

			p := MapPoint{
				x:      x,
				y:      y,
				height: h,
				color:  defaultColor,
			}

			if len(parts) > 1 {
				col, err := rgbaFromHexString(parts[1])
				if err != nil {
					return nil, fmt.Errorf("invalid color %q for %d/%d: %w", parts[1], y, x, err)
				}
				col.A = 255
				p.color = col
			}

			points = append(points, p)
		}
		m = append(m, points)
	}

	if len(m) == 0 {
		return nil, fmt.Errorf("no points")
	}

	return m, nil
}

type MapPoint struct {
	x      int
	y      int
	height int
	color  color.Color
}

// Vector returns the map point as a vec3.
func (mp MapPoint) Vector() vec3 {
	return vec3{
		x:     float64(mp.x),
		y:     float64(mp.y),
		z:     float64(mp.height),
		color: mp.color,
	}
}

// getProjectedBounds projects and scales each points of the map
// and returns the smallest boundaries fitting everything.
//
// Going over each point as the projection of any given point can result
// in a bigger viewport.
func (m *Fdf) getProjectedBounds(scale float64) image.Rectangle {
	var border image.Rectangle

	for _, line := range m.Points {
		for _, elem := range line {
			point := m.cartesianToIsometric(elem.Vector().Scale(scale))

			if math.Floor(point.x) < float64(border.Min.X) {
				border.Min.X = int(math.Floor(point.x))
			} else if math.Ceil(point.x) > float64(border.Max.X) {
				border.Max.X = int(math.Ceil(point.x))
			}
			if math.Floor(point.y) < float64(border.Min.Y) {
				border.Min.Y = int(math.Floor(point.y))
			} else if math.Ceil(point.y) > float64(border.Max.Y) {
				border.Max.Y = int(math.Ceil(point.y))
			}

		}
	}

	return border
}
