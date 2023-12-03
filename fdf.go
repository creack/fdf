package main

import (
	_ "embed"
	"fmt"
	"image"
	"image/color"
	"math"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	WIDTH  = screenWidth
	HEIGHT = screenHeight
)

//go:embed maps/42.fdf
var mapData []byte

type MapPoint struct {
	height int
	color  color.Color
}

type Map struct {
	scale  float64
	offset image.Point
	deg    vec3

	depthChange float64

	Points [][]MapPoint
}

var (
	zDeg float64 = 45
	xDeg float64 = math.Atan(math.Sqrt2)
)

func (m *Map) cartesianToIso(vec vec3) vec3 {
	vec.z *= m.depthChange
	vec = multiplyVectorMatrix(vec, getRotationMatrix(m.deg.z, 'z'))
	vec = multiplyVectorMatrix(vec, getRotationMatrix(m.deg.x, 'x'))
	vec = multiplyVectorMatrix(vec, getRotationMatrix(m.deg.y, 'y'))
	return vec
}

func (m *Map) Update(keys []ebiten.Key) {
	for _, k := range keys {
		switch k {
		case ebiten.Key1:
			m.depthChange -= 0.1
		case ebiten.Key2:
			m.depthChange += 0.1
		case ebiten.Key3:
			m.scale += 1
		case ebiten.Key4:
			m.scale -= 1

		case ebiten.KeyUp:
			m.deg.x += 0.01
		case ebiten.KeyDown:
			m.deg.x -= 0.01
		case ebiten.KeyRight:
			m.deg.y += 0.01
		case ebiten.KeyLeft:
			m.deg.y -= 0.01

		case ebiten.KeyW:
			m.offset.Y -= int(m.scale)
		case ebiten.KeyS:
			m.offset.Y += int(m.scale)
		case ebiten.KeyA:
			m.offset.X -= int(m.scale)
		case ebiten.KeyD:
			m.offset.X += int(m.scale)
		}
	}
}

func (m *Map) Draw(screen *ebiten.Image) {
	img := ebiten.NewImage(WIDTH, HEIGHT)
	for j, line := range m.Points {
		for i, elem := range line {
			v := m.scaleOffset(vec3{
				x:     float64(i),
				y:     float64(j),
				z:     float64(elem.height),
				color: elem.color,
			})

			if i+1 < len(line) {
				v1 := m.scaleOffset(vec3{
					x:     float64(i + 1),
					y:     float64(j),
					z:     float64(m.Points[j][i+1].height),
					color: m.Points[j][i+1].color,
				})
				vector.StrokeLine(img, float32(v.x), float32(v.y), float32(v1.x), float32(v1.y), 1, v1.color, true)
			}
			if j+1 < len(m.Points) && i < len(m.Points[j+1]) {
				v1 := m.scaleOffset(vec3{
					x:     float64(i),
					y:     float64(j + 1),
					z:     float64(m.Points[j+1][i].height),
					color: m.Points[j+1][i].color,
				})
				vector.StrokeLine(img, float32(v.x), float32(v.y), float32(v1.x), float32(v1.y), 1, v1.color, true)
			}

		}
	}
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(img, op)
}

type vec3 struct {
	x, y, z float64
	color   color.Color
}

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
	return deg * (math.Pi / 180)
}

func getRotationMatrix(deg float64, axis byte) matrix3 {
	switch axis {
	case 'x':
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
	case 'y':
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
	case 'z':
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

func rgbaFromHexString(str string) (color.RGBA, error) {
	colStr := strings.TrimPrefix(str, "0x")
	c, err := strconv.ParseUint(colStr, 16, 32)
	if err != nil {
		return color.RGBA{}, fmt.Errorf("parse uint: %w", err)
	}
	return rgbaFromUint(uint32(c)), nil
}

func rgbaFromUint(rgba uint32) color.RGBA {
	cc := color.RGBA{}
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

func loadMap() (*Map, error) {
	m := &Map{}
	for i, line := range strings.Split(string(mapData), "\n") {
		if line == "" {
			continue
		}
		var points []MapPoint
		for j, elem := range strings.Split(line, " ") {
			if elem == "" {
				continue
			}
			parts := strings.Split(elem, ",")

			h, err := strconv.Atoi(parts[0])
			if err != nil {
				return nil, fmt.Errorf("invalid height %q for %d/%d: %w", parts[0], i, j, err)
			}
			p := MapPoint{height: h}

			if len(parts) > 1 {
				col, err := rgbaFromHexString(parts[1])
				if err != nil {
					return nil, fmt.Errorf("invalid color %q for %d/%d: %w", parts[1], i, j, err)
				}
				col.A = 255
				p.color = col
			} else {
				p.color = color.White
			}

			points = append(points, p)
		}
		m.Points = append(m.Points, points)
	}

	if len(m.Points) == 0 {
		return nil, fmt.Errorf("no points")
	}

	m.depthChange = 1
	m.deg.x = xDeg
	m.deg.z = zDeg

	border := m.getMapBorder(1)
	m.scale = getScale(border)
	border = m.getMapBorder(m.scale)
	m.offset = getOffset(border)
	fmt.Printf("%f - %d/%d\n", m.scale, m.offset.X, m.offset.Y)

	return m, nil
}

func (m *Map) toOffset(v vec3, offset image.Point) vec3 {
	nv := m.cartesianToIso(v)
	return vec3{
		x:     nv.x + float64(offset.X),
		y:     nv.y + float64(offset.Y),
		z:     nv.z,
		color: nv.color,
	}
}

func toScale(v vec3, scale float64) vec3 {
	return vec3{
		x:     v.x * scale,
		y:     v.y * scale,
		z:     v.z * scale,
		color: v.color,
	}
}

func (m *Map) scaleOffset(v vec3) vec3 {
	return m.toOffset(toScale(v, m.scale), m.offset)
}

// returns a struct containing min and max of x,y
func (m *Map) getMapBorder(scale float64) image.Rectangle {
	border := image.Rectangle{}

	for j, line := range m.Points {
		for i, elem := range line {
			point := vec3{
				x:     float64(i),
				y:     float64(j),
				z:     float64(elem.height),
				color: elem.color,
			}
			point = m.cartesianToIso(toScale(point, scale))
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
	fmt.Printf("%v\n", border)

	return border
}

func getScale(border image.Rectangle) float64 {
	width := math.Abs(float64(border.Max.X - border.Min.X))
	height := math.Abs(float64(border.Max.Y - border.Min.Y))
	return math.Floor(math.Min((WIDTH-WIDTH/8.0)/width, (HEIGHT-HEIGHT/8.0)/height))
}

func getOffset(border image.Rectangle) image.Point {
	width := math.Abs(float64(border.Max.X - border.Min.X))
	height := math.Abs(float64(border.Max.Y - border.Min.Y))
	offsetX := math.Round((WIDTH - width) / 2.0)
	offsetY := math.Round((HEIGHT - height) / 2.0)
	if border.Min.X < 0 {
		offsetX += math.Abs(math.Round(float64(border.Min.X)))
	}
	if border.Min.Y < 0 {
		offsetY += math.Abs(math.Round(float64(border.Min.Y)))
	}
	return image.Point{X: int(offsetX), Y: int(offsetY)}
}
