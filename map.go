package main

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"

	"go.creack.net/fdf/math3"
)

//nolint:gochecknoglobals // Expected "readonly" global.
var defaultColor = color.White

func parseMap(mapData []byte) ([][]MapPoint, error) {
	// Start by cleaning up the input, removing blank lines and dup spaces.
	//nolint:prealloc // False positive.
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
	//nolint:prealloc // False positive.
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
				Vec: math3.Vec{
					X: float64(x),
					Y: float64(y),
					Z: float64(h),
				},
				color: defaultColor,
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

// MapPoint represents an individual point for the wireframe.
// 3d vector with color.
type MapPoint struct {
	math3.Vec
	color color.Color
}
