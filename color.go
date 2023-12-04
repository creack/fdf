package main

import (
	"fmt"
	"image/color"
	"math"
	"strconv"
	"strings"
)

// rgbaFromHexString parses the given string into RGBA.
func rgbaFromHexString(str string) (color.RGBA, error) {
	colStr := strings.TrimPrefix(str, "0x")
	c, err := strconv.ParseUint(colStr, 16, 32)
	if err != nil {
		return color.RGBA{}, fmt.Errorf("parse uint: %w", err)
	}
	return rgbaFromUint(uint32(c)), nil
}

// rgbaFromUint handles the bitshift magic to parse RGBA.
// The magic is needed to support bot 0xFFFFFF and 0xFFFFFFFF.
func rgbaFromUint(rgba uint32) color.RGBA {
	var cc color.RGBA

	if rgba > 0xFFFFFF {
		cc.A = 8
	}
	cc.R = uint8((rgba >> (16 + cc.A)) & 0xFF)
	cc.G = uint8((rgba >> (8 + cc.A)) & 0xFF)
	cc.B = uint8((rgba >> cc.A) & 0xFF)
	if cc.A != 0 {
		cc.A = uint8(rgba & 0xFF)
	}

	return cc
}

// getGradientColor returns the color in between col1 and col2 with position as strength.
// position is a %, between 0 and 1.
func getGradientColor(col1, col2 color.Color, position float64) color.Color {
	r1, g1, b1, a1 := col1.RGBA()
	r2, g2, b2, a2 := col2.RGBA()
	return color.RGBA{
		R: uint8(math.Round(float64(uint8(r1))*(1-position) + float64(uint8(r2))*position)),
		G: uint8(math.Round(float64(uint8(g1))*(1-position) + float64(uint8(g2))*position)),
		B: uint8(math.Round(float64(uint8(b1))*(1-position) + float64(uint8(b2))*position)),
		A: uint8(math.Round(float64(uint8(a1))*(1-position) + float64(uint8(a2))*position)),
	}
}
