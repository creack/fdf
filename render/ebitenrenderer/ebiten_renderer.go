// Package ebitenrenderer provides a ebiten impementation of the renderer.
package ebitenrenderer

import (
	"fmt"
	"image"

	"go.creack.net/fdf/math3"
	"go.creack.net/fdf/projection"
	"go.creack.net/fdf/render"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Game holds the state.
type Game struct {
	keys []ebiten.Key

	fdf                       render.Engine
	screenWidth, screenHeight int

	img    image.Image
	offset math3.Vec // Offset of the rendered image.

	tainted bool // Flag to know when to redraw img.
}

// Update implements the ebiten interface.
//
// Handles key presses.
func (g *Game) Update() error {
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])

	g.handleFdfKeys(g.keys)

	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		entries := g.fdf.ListMaps()
		i := -1
		for ii, elem := range entries {
			if elem.Name() == g.fdf.CurrentMapName() {
				i = ii
				break
			}
		}
		if err := g.fdf.LoadMap("maps/" + entries[(i+1)%len(entries)].Name()); err != nil {
			return fmt.Errorf("loadMap: %w", err)
		}
		g.img = image.NewRGBA(image.Rect(0, 0, 0, 0))
		g.tainted = true
	}
	return nil
}

func (g *Game) handleFdfKeys(keys []ebiten.Key) {
	p := g.fdf.GetProjection()
	scale := p.GetScale()
	angle := p.GetAngle()
	heightFactor := g.fdf.GetHeightFactor()

loop:
	for _, k := range keys {
		//nolint:exhaustive // False positive.
		switch k {
		case ebiten.KeyW:
			g.offset.Y -= scale
		case ebiten.KeyS:
			g.offset.Y += scale
		case ebiten.KeyA:
			g.offset.X -= scale
		case ebiten.KeyD:
			g.offset.X += scale
		}

		//nolint:exhaustive // False positive.
		switch k {
		// Height factor.
		case ebiten.Key1:
			heightFactor -= 0.1
		case ebiten.Key2:
			heightFactor += 0.1

		// Scale.
		case ebiten.Key3:
			scale++
		case ebiten.Key4:
			scale--

		// Rotations.
		case ebiten.KeyUp:
			angle = angle.Translate(math3.Vec{X: 0.01})
		case ebiten.KeyDown:
			angle = angle.Translate(math3.Vec{X: -0.01})
		case ebiten.KeyRight:
			angle = angle.Translate(math3.Vec{Y: 0.01})
		case ebiten.KeyLeft:
			angle = angle.Translate(math3.Vec{Y: -0.01})
		case ebiten.KeyShiftRight:
			angle = angle.Translate(math3.Vec{Z: 0.01})
		case ebiten.KeyShiftLeft:
			angle = angle.Translate(math3.Vec{Z: -0.01})

		// Misc.
		case ebiten.Key0: // Set all angles to 0.
			angle = math3.Vec{}
		case ebiten.KeyI:
			// Reset to default Isometric.
			if g.screenWidth > 0 && g.screenHeight > 0 {
				render.Iso(g.fdf, g.screenWidth, g.screenHeight)
			}
		default:
			continue loop
		}
		p.SetScale(scale)
		p.SetAngle(angle)
		g.fdf.SetHeightFactor(heightFactor)
		g.tainted = true
	}
}

// Draw implements the ebiten.Game interface.
//
// Draws the rendered wireframe on screen.
func (g *Game) Draw(screen *ebiten.Image) {
	// Get the screen size.
	screenWidth, screenHeight := screen.Bounds().Dx(), screen.Bounds().Dy()
	g.screenWidth, g.screenHeight = screenWidth, screenHeight

	if g.img == nil {
		// Set the isometric projection.
		render.Iso(g.fdf, screenWidth, screenHeight)
		g.tainted = true
	}

	// Render the fdf.
	if g.tainted {
		g.img = g.fdf.Draw()
		g.tainted = false
	}

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
  - HFact: %0.2f
Camera:
 - %0.2f
 - %0.2f
 - %0.2f
`, ebiten.ActualTPS(),
		ebiten.ActualFPS(),
		screenWidth, screenHeight,
		g.img.Bounds(),
		g.fdf.GetProjection().GetScale(),
		g.fdf.GetHeightFactor(),
		g.fdf.GetProjection().GetAngle().X,
		g.fdf.GetProjection().GetAngle().Y,
		g.fdf.GetProjection().GetAngle().Z,
	)
	ebitenutil.DebugPrintAt(screen, msg, screenWidth-150, 1)
	ebitenutil.DebugPrint(screen, fmt.Sprintf(`TPS: %0.2f, FPS: %0.2f
Resolution: %dx%d
Map: %s

Controls:
  W/A/S/D: Move
  Up/Down/Left/Right/Shift Left/Shift Right: Rotate
  C: Cycle maps
  I: Reset view to Isometric
  0: Reset view to 0 angles.
  1/2: Change height scale factor
  3/4: Zoom in/out
`, ebiten.ActualTPS(), ebiten.ActualFPS(), g.screenWidth, g.screenHeight, g.fdf.CurrentMapName()))
}

// Layout implements the ebiten.Game interface.
func (g *Game) Layout(outsideWidth, outsideHeight int) (w, h int) { return outsideWidth, outsideHeight }

type renderer struct{}

// New creates the renderer.
func New(initialWidth, initialHeight int) render.Renderer {
	ebiten.SetWindowSize(initialWidth*2, initialHeight*2)
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
