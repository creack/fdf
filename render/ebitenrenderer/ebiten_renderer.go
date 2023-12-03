package ebitenrenderer

import (
	"fmt"
	"image/color"

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
			// g.fdf.IncScale(1)
		case ebiten.Key4:
			// g.fdf.IncScale(-1)

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
	screenWidth, screenHeight := screen.Bounds().Dx(), screen.Bounds().Dy()

	// Set the isometric projection.
	bounds := g.fdf.SetProjection(projection.NewIsomorphic(projection.GetScale(screenWidth, screenHeight, g.fdf.SetProjection(projection.NewIsomorphic(1)))))

	vector.StrokeRect(screen, float32(bounds.Min.X), float32(bounds.Min.Y), float32(bounds.Max.X), float32(bounds.Max.Y), 2, color.RGBA{A: 255, R: 255}, false)

	// NOTE: NewImageFromImage automatically fits the image back into 0,0 origin.
	fdfImg := ebiten.NewImageFromImage(g.fdf.Draw())
	bounds = fdfImg.Bounds()

	vector.StrokeRect(screen, 1, 1, float32(screenWidth-2), float32(screenHeight-2), 2, color.White, false)

	vector.StrokeRect(fdfImg, float32(bounds.Min.X), float32(bounds.Min.Y), float32(bounds.Max.X), float32(bounds.Max.Y), 2, color.RGBA{A: 255, B: 255}, false)

	op := &ebiten.DrawImageOptions{}
	offset := projection.GetOffsetCenter(screenWidth, screenHeight, bounds)
	op.GeoM.Translate(float64(offset.X), float64(offset.Y))
	screen.DrawImage(fdfImg, op)

	msg := fmt.Sprintf(`TPS: %0.2f
FPS: %0.2f
Sizes:
  - Screen: %d/%d
  - Fdf: %v
`, ebiten.ActualTPS(),
		ebiten.ActualFPS(),
		screenWidth, screenHeight,
		bounds,
	)
	ebitenutil.DebugPrintAt(screen, msg, screenWidth-150, 1)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

type renderer struct{}

func New(width, height int) render.Renderer {
	ebiten.SetWindowSize(width*2, height*2)
	ebiten.SetWindowTitle("FDF")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	return &renderer{}
}

func (r *renderer) Run(fdf render.Engine) error {
	g := &Game{fdf: fdf}

	if err := ebiten.RunGame(g); err != nil {
		return fmt.Errorf("runGame: %w", err)
	}
	return nil
}
