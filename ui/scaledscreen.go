package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/tinne26/etxt"
	"github.com/tinne26/etxt/fract"
)

type ScaledScreen struct {
	Screen         *ebiten.Image
	scaleFactor    float64
	Etxt           *etxt.Renderer
	debugPrintLoc  fract.Point
	debugPrintSize float64
}

func NewScaledScreen(renderer *etxt.Renderer) *ScaledScreen {
	return &ScaledScreen{
		Etxt:           renderer,
		scaleFactor:    ebiten.Monitor().DeviceScaleFactor(),
		debugPrintSize: 16 * ebiten.Monitor().DeviceScaleFactor(),
	}
}

func (s *ScaledScreen) SetTarget(t *ebiten.Image) {
	s.Screen = t
	s.debugPrintLoc = fract.IntsToPoint(0, 0)
}

func (s *ScaledScreen) DrawImage(image *ebiten.Image, options *ebiten.DrawImageOptions) {
	options.GeoM.Scale(s.scaleFactor, s.scaleFactor)
	s.Screen.DrawImage(image, options)
}

func (s *ScaledScreen) DrawRect(x, y, w, h float64, color color.Color) {
	xx := float32(x * s.scaleFactor)
	yy := float32(y * s.scaleFactor)
	hh := float32(h * s.scaleFactor)
	ww := float32(w * s.scaleFactor)

	vector.DrawFilledRect(s.Screen, xx, yy, ww, hh, color, false)
}

func (s *ScaledScreen) DrawUnfilledRect(x, y, w, h, strokeWidth float64, color color.Color) {
	xx := float32(x * s.scaleFactor)
	yy := float32(y * s.scaleFactor)
	hh := float32(h * s.scaleFactor)
	ww := float32(w * s.scaleFactor)
	sw := float32(strokeWidth * s.scaleFactor)

	vector.StrokeRect(s.Screen, xx, yy, ww, hh, sw, color, false)
}

func (s *ScaledScreen) DrawCircle(cx, cy, r float64, color color.Color) {
	xx := float32(cx * s.scaleFactor)
	yy := float32(cy * s.scaleFactor)
	rr := float32(r * s.scaleFactor)

	vector.DrawFilledCircle(s.Screen, xx, yy, rr, color, false)
}

func (s *ScaledScreen) DrawRectShader(w, h int, shader *ebiten.Shader, opts *ebiten.DrawRectShaderOptions) {
	ww := int(float64(w) * s.scaleFactor)
	hh := int(float64(h) * s.scaleFactor)

	opts.GeoM.Scale(s.scaleFactor, s.scaleFactor)
	s.Screen.DrawRectShader(ww, hh, shader, opts)
}

func (s *ScaledScreen) scaledTextSize(size float64) float64 {
	return size * s.scaleFactor
}

func (s *ScaledScreen) TextSelectionRectSize(t string, size float64) (float64, float64) {
	s.Etxt.SetSize(s.scaledTextSize(size))
	r := s.Etxt.Measure(t)
	return r.Width().ToFloat64(), r.Height().ToFloat64()
}

func (s *ScaledScreen) DrawText(t string, size float64, x, y int, color color.Color) {
	xx := int(float64(x) * s.scaleFactor)
	yy := int(float64(y) * s.scaleFactor)

	s.Etxt.SetColor(color)
	s.Etxt.SetSize(s.scaledTextSize(size))
	s.Etxt.SetAlign(etxt.Top | etxt.Left)
	s.Etxt.Draw(s.Screen, t, xx, yy)
}

func (s *ScaledScreen) DrawTextCenteredAt(t string, size float64, x, y int, color color.Color) {
	xx := int(float64(x) * s.scaleFactor)
	yy := int(float64(y) * s.scaleFactor)

	s.Etxt.SetColor(color)
	s.Etxt.SetSize(s.scaledTextSize(size))
	s.Etxt.SetAlign(etxt.HorzCenter | etxt.VertCenter)
	s.Etxt.Draw(s.Screen, t, xx, yy)
}

func (s *ScaledScreen) DrawTextWithAlign(t string, size float64, x, y int, color color.Color, vAlign etxt.Align, hAlign etxt.Align) {
	xx := int(float64(x) * s.scaleFactor)
	yy := int(float64(y) * s.scaleFactor)

	s.Etxt.SetColor(color)
	s.Etxt.SetSize(s.scaledTextSize(size))
	s.Etxt.SetAlign(vAlign | hAlign)
	s.Etxt.Draw(s.Screen, t, xx, yy)
}

func (s *ScaledScreen) DebugPrint(str string) {
	s.Etxt.SetSize(s.debugPrintSize)
	s.Etxt.SetAlign(etxt.Top | etxt.Left)
	s.Etxt.SetColor(color.White)
	r := s.Etxt.Measure(str)
	s.Etxt.Draw(s.Screen, str+"\n", s.debugPrintLoc.X.ToInt(), s.debugPrintLoc.Y.ToInt())
	s.debugPrintLoc = s.debugPrintLoc.AddUnits(fract.FromInt(0), r.Height())
}

func AdjustedCursorPosition() (float64, float64) {
	cx, cy := ebiten.CursorPosition()
	return float64(cx) / ebiten.Monitor().DeviceScaleFactor(), float64(cy) / ebiten.Monitor().DeviceScaleFactor()
}
