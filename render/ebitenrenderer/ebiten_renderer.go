package ebitenrenderer

import (
	"fmt"
	_ "image/png"

	"fdf/render"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 1024
	screenHeight = 1024
)

type Game struct {
	keys []ebiten.Key

	Fdf render.Engine
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
			g.Fdf.IncScale(1)
		case ebiten.Key4:
			g.Fdf.IncScale(-1)

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
	fdfImg := ebiten.NewImageFromImage(g.Fdf.Draw())

	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(fdfImg, op)

	msg := fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS())
	ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func Render(fdf render.Engine) error {
	g := &Game{Fdf: fdf}
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("FDF")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(g); err != nil {
		return fmt.Errorf("runGame: %w", err)
	}
	return nil
}
