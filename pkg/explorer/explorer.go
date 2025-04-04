package explorer

import (
	"bytes"
	_ "embed"
	"fmt"
	"image/color"

	"github.com/rhgrant10/mandlebrot-explorer/pkg/geometry"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

//go:embed InterVariable.ttf
var ttfBytes []byte

func loadFont() text.Face {
	fontFile := bytes.NewReader(ttfBytes)
	faceSource, err := text.NewGoTextFaceSource(fontFile)
	if err != nil {
		panic(err)
	}
	return &text.GoTextFace{
		Source: faceSource,
		Size:   18,
	}
}

func NewGame(window Window[float64], image Graph) Game {
	originalWindowRect := window.Rect
	return Game{
		window:             window,
		graph:              image,
		selection:          NewSelection(10),
		font:               loadFont(),
		originalWindowRect: originalWindowRect,
	}
}

// Game implements the ebiten.Game interface.
type Game struct {
	window             Window[float64]
	graph              Graph
	selection          Selection
	font               text.Face
	isFresh            bool
	originalWindowRect geometry.Rect[float64]
}

// Update updates the game logic.
func (g *Game) Update() error {
	if !g.isFresh {
		if !g.graph.IsRendering() {
			go func() {
				ebiten.SetCursorShape(ebiten.CursorShapeNotAllowed)
				defer ebiten.SetCursorShape(ebiten.CursorShapeCrosshair)
				g.graph.Render(g.window)
				g.isFresh = true
			}()
		}
	} else {
		if ebiten.IsKeyPressed(ebiten.KeyEqual) {
			g.adjustIterLimit(20)
		}
		if ebiten.IsKeyPressed(ebiten.KeyMinus) {
			g.adjustIterLimit(-20)
		}
		if ebiten.IsKeyPressed(ebiten.KeyR) {
			g.resetWindow()
		}
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			g.isFresh = false
		}
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) ||
			inpututil.IsMouseButtonJustPressed(ebiten.MouseButton2) {
			g.selection.Start()
		}
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) {
			g.setWindow(g.transformRect(g.selection.End()))
		}
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButton2) {
			sel := g.selection.End()
			selWindow := g.transformRect(sel)
			scaleFactor := float64(sel.Width()) / float64(g.graph.Resolution().X)
			from := g.window.CenterPoint()
			to := selWindow.CenterPoint()
			rect := g.window.Rect
			(rect.
				Translate(*from.Negate()).
				Scale(1.0 / scaleFactor).
				Translate(to))
			g.setWindow(rect)
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyEscape) ||
		ebiten.IsKeyPressed(ebiten.KeyControl) &&
			ebiten.IsKeyPressed(ebiten.KeyC) {
		return ebiten.Termination
	}
	return nil
}

// Draw renders the game screen.
func (g *Game) Draw(screen *ebiten.Image) {
	g.graph.DrawOn(screen)
	g.selection.Render(screen, g.window)
	debugText := g.getDebugText()
	ops := text.DrawOptions{}
	ops.LineSpacing = 32
	ops.GeoM.Translate(16, 16)
	ops.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, debugText, g.font, &ops)
	ops.GeoM.Translate(1, 1)
	ops.ColorScale.ScaleWithColor(color.Black)
	text.Draw(screen, debugText, g.font, &ops)
	ops.GeoM.Translate(1, 1)
	text.Draw(screen, debugText, g.font, &ops)
}

// Layout defines the screen layout.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	size := g.graph.Resolution()
	return size.X, size.Y
}

func (g *Game) transformRect(rect geometry.Rect[int]) geometry.Rect[float64] {
	size := g.graph.Resolution()
	return geometry.Rect[float64]{
		Min: g.window.Transform(rect.Min.X, rect.Min.Y, size.X, size.Y),
		Max: g.window.Transform(rect.Max.X, rect.Max.Y, size.X, size.Y),
	}
}

func (g *Game) adjustIterLimit(change int) {
	g.graph.Iterator.IterLimit = max(64, g.graph.Iterator.IterLimit+change)
	g.graph.Colorizer.SetResolution(g.graph.Iterator.IterLimit)
}

func (g *Game) setWindow(rect geometry.Rect[float64]) {
	g.window.Rect = rect
	g.isFresh = false
}

func (g *Game) resetWindow() {
	newWindowRect := g.originalWindowRect
	g.window.Rect = newWindowRect
	g.isFresh = false
}

func (g *Game) getDebugText() string {
	size := g.graph.Resolution()
	mouse := g.selection.GetMouse()
	hover := g.window.Transform(mouse.X, mouse.Y, size.Y, size.Y)
	text := fmt.Sprintf(
		"iterLimit: %d\nmouse X: % 5d (% .16e)\nmouse Y: % 5d (% .16e)\n",
		g.graph.Iterator.IterLimit, mouse.X, hover.X, mouse.Y, hover.Y,
	)
	if g.selection.IsBeingMade() {
		text = fmt.Sprintf("%s\n%v", text, g.selection.Get())
	}
	return text
}
