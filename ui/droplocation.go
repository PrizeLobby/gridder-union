package ui

import "image/color"

type DropLocation struct {
	X, Y, W, H float64
	Index      int
	SetSprite  *SetSprite
}

func (d *DropLocation) Contains(x, y float64) bool {
	return x >= d.X && x < d.X+d.W && y >= d.Y && y < d.Y+d.H
}
func (d *DropLocation) Draw(screen *ScaledScreen, extras []color.Color) {
	screen.DrawRect(float64(d.X), float64(d.Y), float64(d.W), float64(d.H), color.RGBA{124, 194, 154, 255})
	screen.DrawUnfilledRect(float64(d.X), float64(d.Y), float64(d.W), float64(d.H), 10, color.RGBA{101, 153, 145, 255})
	if d.SetSprite != nil {
		d.SetSprite.DrawWithColors(screen, extras)
	}
}

func (d *DropLocation) Update() {
	// This is a placeholder for any update logic needed for the drop location.
	// Currently, it does nothing.
}
