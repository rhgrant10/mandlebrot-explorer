package explorer

import (
	"fmt"
	"image/color"

	"github.com/rhgrant10/mandlebrot-explorer/pkg/geometry"

	"golang.org/x/image/font/basicfont"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

func NewExplorer(screen Screen, window Window[float64], image Graph) Explorer {
	originalWindowRect := window.Rect
	return Explorer{
		Screen:             screen,
		Window:             window,
		Graph:              image,
		selection:          NewSelection(10),
		originalWindowRect: originalWindowRect,
	}
}

// Explorer implements the ebiten.Game interface.
type Explorer struct {
	Screen             Screen
	Window             Window[float64]
	Graph              Graph
	selection          Selection
	isFresh            bool
	originalWindowRect geometry.Rect[float64]
}

func (g *Explorer) ResetWindow() {
	newWindowRect := g.originalWindowRect
	g.Window.Rect = newWindowRect
	g.isFresh = false
}

// Update updates the game logic.
func (g *Explorer) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEqual) {
		g.AdjustIterLimit(20)
	}
	if ebiten.IsKeyPressed(ebiten.KeyMinus) {
		g.AdjustIterLimit(-20)
	}
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		g.ResetWindow()
	}
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		g.isFresh = false
	}
	if !g.isFresh && !g.Graph.IsRendering() {
		go func() {
			ebiten.SetCursorShape(ebiten.CursorShapeNotAllowed)
			defer ebiten.SetCursorShape(ebiten.CursorShapeCrosshair)
			g.Graph.Render(g.Window)
			g.isFresh = true
		}()
	}
	if !g.Graph.IsRendering() && (inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) ||
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButton2)) {
		g.selection.Start()
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) {
		g.SetWindow(g.TransformRect(g.selection.End()))
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButton2) {
		sel := g.selection.End()
		selWindow := g.TransformRect(sel)
		scaleFactor := float64(sel.Width()) / float64(g.Screen.Width())
		from := g.Window.CenterPoint()
		to := selWindow.CenterPoint()
		rect := g.Window.Rect
		(rect.
			Translate(*from.Negate()).
			Scale(1.0 / scaleFactor).
			Translate(to))
		g.SetWindow(rect)
	}
	if ebiten.IsKeyPressed(ebiten.KeyEscape) ||
		ebiten.IsKeyPressed(ebiten.KeyControl) &&
			ebiten.IsKeyPressed(ebiten.KeyC) {
		return ebiten.Termination
	}
	return nil
}

func (g *Explorer) TransformRect(rect geometry.Rect[int]) geometry.Rect[float64] {
	return geometry.Rect[float64]{
		Min: g.Window.Transform(rect.Min.X, rect.Min.Y, g.Screen.Width(), g.Screen.Height()),
		Max: g.Window.Transform(rect.Max.X, rect.Max.Y, g.Screen.Width(), g.Screen.Height()),
	}
}

func (g *Explorer) AdjustIterLimit(change int) {
	g.Graph.Iterator.IterLimit = max(64, g.Graph.Iterator.IterLimit+change)
	g.Graph.Colorizer.SetResolution(g.Graph.Iterator.IterLimit)
}

func (g *Explorer) SetWindow(rect geometry.Rect[float64]) {
	g.Window.Rect = rect
	g.isFresh = false
}

// Draw renders the game screen.
func (g *Explorer) Draw(screen *ebiten.Image) {
	g.Graph.DrawOn(screen)
	g.selection.Render(screen, g.Window)
	debugText := g.getDebugText()
	ff := text.FaceWithLineHeight(basicfont.Face7x13, 16)
	text.Draw(screen, debugText, ff, 18, 18, color.Black)
	text.Draw(screen, debugText, ff, 17, 17, color.Black)
	text.Draw(screen, debugText, ff, 16, 16, color.White)
}

func (g *Explorer) getDebugText() string {
	mouse := g.selection.GetMouse()
	hover := g.Window.Transform(mouse.X, mouse.Y, g.Screen.Width(), g.Screen.Height())
	text := fmt.Sprintf(
		"iterLimit: %d\nmouse X: % 5d (% .16e)\nmouse Y: % 5d (% .16e)\n",
		g.Graph.Iterator.IterLimit, mouse.X, hover.X, mouse.Y, hover.Y,
	)
	if g.selection.IsBeingMade() {
		text = fmt.Sprintf("%s\n%v", text, g.selection.Get())
	}
	return text
}

// Layout defines the screen layout.
func (g *Explorer) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.Screen.Width(), g.Screen.Height()
}
