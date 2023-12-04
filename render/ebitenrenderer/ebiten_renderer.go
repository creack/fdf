package ebitenrenderer

import (
	"fmt"
	"image"

	"fdf/math3"
	"fdf/projection"
	"fdf/render"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	keys []ebiten.Key

	fdf render.Engine

	offset math3.Vec3
	img    image.Image
}

func (g *Game) Update() error {
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])

	g.handleFdfKeys(g.keys)
	return nil
}

func (g *Game) handleFdfKeys(keys []ebiten.Key) {
	p := g.fdf.GetProjection()
	scale := p.GetScale()
	angle := p.GetAngle()

	for _, k := range keys {
		switch k {
		// case ebiten.Key1:
		// 	m.depthChange -= 0.1
		// case ebiten.Key2:
		// 	m.depthChange += 0.1
		case ebiten.Key3:
			scale++
			p.SetScale(scale)
			// g.fdf.SetProjection(p)
			// g.img = g.fdf.Draw()
		case ebiten.Key4:
			scale--
			p.SetScale(scale)
			// g.fdf.SetProjection(p)
			// g.img = g.fdf.Draw()

		case ebiten.KeyUp:
			angle = angle.Translate(math3.Vec3{X: 0.01})
			p.SetAngle(angle)
		case ebiten.KeyDown:
			angle = angle.Translate(math3.Vec3{X: -0.01})
			p.SetAngle(angle)
		case ebiten.KeyRight:
			angle = angle.Translate(math3.Vec3{Y: 0.01})
			p.SetAngle(angle)
		case ebiten.KeyLeft:
			angle = angle.Translate(math3.Vec3{Y: -0.01})
			p.SetAngle(angle)
		case ebiten.KeyShiftRight:
			angle = angle.Translate(math3.Vec3{Z: 0.01})
			p.SetAngle(angle)
		case ebiten.KeyShiftLeft:
			angle = angle.Translate(math3.Vec3{Z: -0.01})
			p.SetAngle(angle)

		case ebiten.Key0:
			angle = math3.Vec3{}
			p.SetAngle(angle)

		case ebiten.KeyW:
			g.offset.Y -= scale
		case ebiten.KeyS:
			g.offset.Y += scale
		case ebiten.KeyA:
			g.offset.X -= scale
		case ebiten.KeyD:
			g.offset.X += scale
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Get the screen size.
	screenWidth, screenHeight := screen.Bounds().Dx(), screen.Bounds().Dy()

	if g.img == nil {
		// Set the isometric projection.
		render.Iso(g.fdf, screenWidth, screenHeight)
	}

	// Render the fdf.
	g.img = g.fdf.Draw()

	// Draw the rendered image on the screen.
	op := &ebiten.DrawImageOptions{}
	// Translate to center.
	centerOffset := projection.GetOffsetCenter(screenWidth, screenHeight, g.img.Bounds())
	op.GeoM.Translate(float64(centerOffset.X), float64(centerOffset.Y))
	// Translate to local offset (controlled by keyboard).
	op.GeoM.Translate(g.offset.X, g.offset.Y)

	screen.DrawImage(ebiten.NewImageFromImage(g.img), op)

	msg := fmt.Sprintf(`TPS: %0.2f
FPS: %0.2f
Sizes:
  - Screen: %d/%d
  - Bounds: %v
  - Scale: %0.2f
Camera:
 - %0.2f
 - %0.2f
 - %0.2f
`, ebiten.ActualTPS(),
		ebiten.ActualFPS(),
		screenWidth, screenHeight,
		g.img.Bounds(),
		g.fdf.GetProjection().GetScale(),
		g.fdf.GetProjection().GetAngle().X,
		g.fdf.GetProjection().GetAngle().Y,
		g.fdf.GetProjection().GetAngle().Z,
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
