package ebitenrenderer

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"fdf/math3"
	"fdf/projection"
	"fdf/render"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	keys []ebiten.Key

	fdf render.Engine
}

func (g *Game) Update() error {
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])

	g.handleFdfKeys(g.keys)
	return nil
}

func (g *Game) handleFdfKeys(keys []ebiten.Key) {
	for _, k := range keys {
		switch k {
		// case ebiten.Key1:
		// 	m.depthChange -= 0.1
		// case ebiten.Key2:
		// 	m.depthChange += 0.1
		case ebiten.Key3:
			g.fdf.IncScale(1)
		case ebiten.Key4:
			g.fdf.IncScale(-1)

			// case ebiten.KeyUp:
			// 	m.deg.x += 0.01
			// case ebiten.KeyDown:
			// 	m.deg.x -= 0.01
			// case ebiten.KeyRight:
			// 	m.deg.y += 0.01
			// case ebiten.KeyLeft:
			// 	m.deg.y -= 0.01

			// case ebiten.KeyW:
			// 	m.offset.Y -= int(m.scale)
			// case ebiten.KeyS:
			// 	m.offset.Y += int(m.scale)
			// case ebiten.KeyA:
			// 	m.offset.X -= int(m.scale)
			// case ebiten.KeyD:
			// 	m.offset.X += int(m.scale)
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Get the screen size.
	screenWidth := screen.Bounds().Dx()
	screenHeight := screen.Bounds().Dy()

	vector.StrokeRect(screen, 1, 1, float32(screenWidth-2), float32(screenHeight-2), 2, color.White, false)

	// Set the inital projection with scale 1, no offset, no camera rotation.
	g.fdf.SetProjection(projection.NewIsomorphic(1, image.Point{}, math3.Vec3{}))

	// Get the inital bounds.
	bounds := g.fdf.Draw().Bounds()

	vector.StrokeRect(screen, float32(bounds.Min.X), float32(bounds.Min.Y), float32(bounds.Max.X), float32(bounds.Max.Y), 2, color.RGBA{A: 255, B: 255}, false)

	// Lookup the scale based on the screen size and the scaled boundaries.
	s := float64(projection.GetScale(screenWidth, screenHeight, bounds))
	// Update the projection with the new scale and set default camera.
	g.fdf.SetProjection(projection.NewIsomorphic(int(s), image.Point{}, math3.Vec3{
		X: math.Atan(math.Sqrt2),
		Z: 45,
	}))
	bounds = g.fdf.Draw().Bounds()

	vector.StrokeRect(screen, float32(bounds.Min.X), float32(bounds.Min.Y), float32(bounds.Max.X), float32(bounds.Max.Y), 2, color.RGBA{A: 255, G: 255}, false)

	offset := projection.GetOffsetCenter(screenWidth, screenHeight, bounds)
	g.fdf.SetProjection(projection.NewIsomorphic(int(s), offset, math3.Vec3{
		X: math.Atan(math.Sqrt2),
		Z: 45,
	}))

	vector.DrawFilledCircle(screen, float32(offset.X), float32(offset.Y), 5, color.White, false)
	vector.DrawFilledCircle(screen, float32(screenWidth)/2, float32(screenHeight)/2, 15, color.RGBA{A: 255, G: 255, B: 100}, false)

	bounds = g.fdf.Draw().Bounds()
	offset2 := projection.GetOffsetCenter(screenWidth, screenHeight, bounds)
	vector.StrokeRect(screen, float32(bounds.Min.X), float32(bounds.Min.Y), float32(bounds.Max.X), float32(bounds.Max.Y), 2, color.RGBA{A: 255, R: 255}, false)

	col11 := g.fdf.Draw().At(21, 12)

	fdfImg := ebiten.NewImageFromImage(g.fdf.Draw())

	op := &ebiten.DrawImageOptions{}
	// op.GeoM.Translate(float64(offset.X), float64(offset.Y))
	screen.DrawImage(fdfImg, op)

	msg := fmt.Sprintf(`TPS: %0.2f
FPS: %0.2f
Scale: %0.2f
Sizes:
  - Screen: %d/%d
  - Fdf: %d/%d
Offset0: %d/%d
Offset1: %d/%d
Col11: %v
`, ebiten.ActualTPS(),
		ebiten.ActualFPS(),
		s,
		screenWidth, screenHeight,
		bounds.Dx(), bounds.Dy(),
		offset.X, offset.Y,
		offset2.X, offset2.Y,
		col11,
	)
	ebitenutil.DebugPrintAt(screen, msg, screenWidth-150, 1)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

type renderer struct {
	width, height int
}

func New(width, height int) render.Renderer {
	ebiten.SetWindowSize(width*2, height*2)
	ebiten.SetWindowTitle("FDF")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	return &renderer{
		width:  width,
		height: height,
	}
}

func (r *renderer) Run(fdf render.Engine) error {
	g := &Game{fdf: fdf}

	if err := ebiten.RunGame(g); err != nil {
		return fmt.Errorf("runGame: %w", err)
	}
	return nil
}
