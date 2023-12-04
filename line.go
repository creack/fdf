package main

import (
	"image"
	"image/color"
	"math"
)

// drawLine from p0 to p1.
func drawLine(dst *image.RGBA, p0, p1 image.Point, col1, col2 color.Color) {
	rect := image.Rectangle{Min: p0, Max: p1}.Canon()
	if rect.Dx() > rect.Dy() {
		drawLineHoriz(dst, p0, p1, col1, col2)
	} else {
		drawLineVert(dst, p0, p1, col1, col2)
	}
}

// drawLineHoriz along the X axis.
func drawLineHoriz(dst *image.RGBA, p0, p1 image.Point, col1, col2 color.Color) {
	// Start with the lowest point, swap if needed.
	if p1.X-p0.X < 0 {
		p0, p1 = p1, p0
	}

	// Dims.
	dx, dy := p1.X-p0.X, p1.Y-p0.Y

	// Other axis direction.
	yDir := 1
	if dy < 0 {
		yDir = -1
		dy *= -1
	}

	d := 2*dy - dx
	for x, y := p0.X, p0.Y; x <= p1.X; x++ {
		dst.Set(x, y, lookupGradient(col1, col2, p0, p1, x, y))
		if d > 0 {
			y += yDir
			d += -2 * dx
		}
		d += 2 * dy
	}
}

// drawLineVert along the Y axis.
func drawLineVert(dst *image.RGBA, p0, p1 image.Point, col1, col2 color.Color) {
	// Start with the lowest point, swap if needed.
	if p1.Y-p0.Y < 0 {
		p0, p1 = p1, p0
	}

	// Dims.
	dx, dy := p1.X-p0.X, p1.Y-p0.Y

	// Other axis direction.
	xDir := 1
	if dx < 0 {
		xDir = -1
		dx *= -1
	}

	d := 2*dx - dy
	for x, y := p0.X, p0.Y; y <= p1.Y; y++ {
		dst.Set(x, y, lookupGradient(col1, col2, p0, p1, x, y))
		if d > 0 {
			x += xDir
			d += -2 * dy
		}
		d += 2 * dx
	}
}

// getLinePosition returns the % of the current position from start to end.
func getLinePosition(start, end image.Point, curX, curY int) float64 {
	// a^2 + b^2 = c^2. We have a and b.
	fullLength := math.Sqrt(math.Pow(float64(end.X-start.X), 2) + math.Pow(float64(end.Y-start.Y), 2))
	curLength := math.Sqrt(math.Pow(float64(curX-start.X), 2) + math.Pow(float64(curY-start.Y), 2))
	if fullLength != 0 {
		return curLength / fullLength
	}
	return 1
}

// lookupGradient gets the relative current position in the
// line and returns the color gradient.
func lookupGradient(col1, col2 color.Color, start, end image.Point, curX, curY int) color.Color {
	return getGradientColor(col1, col2, getLinePosition(start, end, curX, curY))
}
