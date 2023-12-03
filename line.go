package main

import (
	"image"
	"image/color"
)

// drawLine from p0 to p1.
func drawLine(dst *image.RGBA, p0, p1 image.Point, clr color.Color) {
	rect := image.Rectangle{Min: p0, Max: p1}.Canon()
	if rect.Dx() > rect.Dy() {
		drawLineHoriz(dst, p0, p1, clr)
	} else {
		drawLineVert(dst, p0, p1, clr)
	}
}

// drawLineHoriz along the X axis.
func drawLineHoriz(dst *image.RGBA, p0, p1 image.Point, clr color.Color) {
	// Start with the lowest point, swap if needed.
	if p1.X-p0.X < 0 {
		p0, p1 = p1, p0
	}

	// Dims.
	dx, dy := int(p1.X-p0.X), int(p1.Y-p0.Y)

	// Other axis direction.
	yDir := 1
	if dy < 0 {
		yDir = -1
		dy *= -1
	}

	d := (2 * dy) - dx
	for x, y := p0.X, p0.Y; x <= p1.X; x++ {
		dst.Set(int(x), int(y), clr)
		if d > 0 {
			y += yDir
			d += -2 * dx
		}
		d += 2 * dy
	}
}

// drawLineVert along the Y axis.
func drawLineVert(dst *image.RGBA, p0, p1 image.Point, clr color.Color) {
	// Start with the lowest point, swap if needed.
	if p1.Y-p0.Y < 0 {
		p0, p1 = p1, p0
	}

	// Dims.
	dx, dy := int(p1.X-p0.X), int(p1.Y-p0.Y)

	// Other axis direction.
	xDir := 1
	if dx < 0 {
		xDir = -1
		dx *= -1
	}

	d := 2*dx - dy
	for x, y := p0.X, p0.Y; y <= p1.Y; y++ {
		dst.Set(int(x), int(y), clr)
		if d > 0 {
			x += xDir
			d += -2 * dy
		}
		d += 2 * dx
	}
}
