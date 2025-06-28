package ui

import "image/color"

type SetSprite struct {
	SpriteName string
	X, Y       float64
}

const SetSpriteWidth = 70
const SetSpriteHeight = 36

func (s *SetSprite) Update() {
	// This is a placeholder for any update logic needed for the sprite.
	// Currently, it does nothing.
}
func (s *SetSprite) Draw(screen *ScaledScreen) {
	screen.DrawRect(float64(s.X), float64(s.Y), SetSpriteWidth, SetSpriteHeight, color.RGBA{255, 255, 255, 255})
	screen.DrawTextCenteredAt(s.SpriteName, 32, int(s.X+SetSpriteWidth/2), int(s.Y+SetSpriteHeight/2), color.Black)
	screen.DrawUnfilledRect(float64(s.X), float64(s.Y), SetSpriteWidth, SetSpriteHeight, 4, color.RGBA{50, 60, 55, 255})
}
func (s *SetSprite) DrawWithColors(screen *ScaledScreen, colors []color.Color) {
	screen.DrawRect(float64(s.X), float64(s.Y), SetSpriteWidth, SetSpriteHeight, color.RGBA{255, 255, 255, 255})
	screen.DrawTextCenteredAtWithColors(s.SpriteName, 32, int(s.X+SetSpriteWidth/2), int(s.Y+SetSpriteHeight/2), colors)
	screen.DrawUnfilledRect(float64(s.X), float64(s.Y), SetSpriteWidth, SetSpriteHeight, 4, color.RGBA{50, 60, 55, 255})
}

func (s *SetSprite) MoveTo(x, y float64) {
	s.X = x
	s.Y = y
}

func (s *SetSprite) MoveBy(dx, dy float64) {
	s.X += dx
	s.Y += dy
}

func (s *SetSprite) Contains(x, y float64) bool {
	return x >= s.X && x < s.X+SetSpriteWidth && y >= s.Y && y < s.Y+SetSpriteHeight
}
