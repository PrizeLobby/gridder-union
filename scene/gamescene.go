package scene

import (
	"image/color"
	"slices"
	"strings"

	"github.com/prizelobby/union-gridder/core"

	"github.com/prizelobby/union-gridder/ui"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const SET_START_X = 75
const SET_START_Y = 145

type GameScene struct {
	BaseScene
	Game          *core.Game
	Setsprites    []*ui.SetSprite
	Droplocations []*ui.DropLocation
	Colors        [6][]color.Color
	Stroke        *ui.Stroke
}

func NewGameScene(game *core.Game) *GameScene {
	setSprites := make([]*ui.SetSprite, 0, core.NUM_SETS)

	for i, s := range game.Sets {
		setSprites = append(setSprites, &ui.SetSprite{
			SpriteName: s,
			X:          SET_START_X,
			Y:          SET_START_Y + float64(i)*50,
		})
	}

	dropLocations := make([]*ui.DropLocation, 0, 9)
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			dropLocations = append(dropLocations, &ui.DropLocation{
				X:     960/2 - 180 - 60 + float64(i)*180,
				Y:     120 + float64(j)*180,
				W:     120,
				H:     120,
				Index: i + 3*j,
			})
		}
	}

	colors := [6][]color.Color{}
	for i := range 6 {
		colors[i] = make([]color.Color, len(game.Targets[i]))
		for j := range len(colors[i]) {
			colors[i][j] = color.RGBA{0, 0, 0, 255}
		}
	}

	return &GameScene{
		Game:          game,
		Setsprites:    setSprites,
		Droplocations: dropLocations,
		Colors:        colors,
	}
}

func (g *GameScene) Draw(screen *ui.ScaledScreen) {
	screen.Screen.Fill(color.RGBA{230, 228, 213, 255})
	for _, loc := range g.Droplocations {
		loc.Draw(screen)
	}
	screen.DrawTextCenteredAt("Gridder Union", 64, 960/2, 60, color.Black)
	screen.DrawTextWithColors(string(g.Game.Targets[0]), 32, 740, 120+60-ui.SetSpriteHeight/2, g.Colors[0])
	screen.DrawTextWithColors(string(g.Game.Targets[1]), 32, 740, 120+180+60-ui.SetSpriteHeight/2, g.Colors[1])
	screen.DrawTextWithColors(string(g.Game.Targets[2]), 32, 740, 120+2*180+60-ui.SetSpriteHeight/2, g.Colors[2])
	screen.DrawTextCenteredAtWithColors(string(g.Game.Targets[3]), 32, 960/2-180, 630, g.Colors[3])
	screen.DrawTextCenteredAtWithColors(string(g.Game.Targets[4]), 32, 960/2, 630, g.Colors[4])
	screen.DrawTextCenteredAtWithColors(string(g.Game.Targets[5]), 32, 960/2+180, 630, g.Colors[5])

	if g.Game.Solved {
		screen.DrawTextCenteredAt("You solved the puzzle!", 40, 960/2, 680, color.RGBA{90, 190, 90, 255})
	}

	for _, sprite := range g.Setsprites {
		sprite.Draw(screen)
	}
	if g.Stroke != nil {
		g.Stroke.DraggingObject.(*ui.SetSprite).Draw(screen)
	}
}

func (g *GameScene) Update() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		cursorX, cursorY := ui.AdjustedCursorPosition()
		selectedIndex := -1
		for i, setSprite := range g.Setsprites {
			if setSprite.Contains(cursorX, cursorY) {
				g.Stroke = ui.NewStroke(cursorX, cursorY, setSprite)
				selectedIndex = i
				break
			}
		}
		if selectedIndex != -1 {
			g.Setsprites = append(g.Setsprites[:selectedIndex], g.Setsprites[selectedIndex+1:]...)
		}
		for _, loc := range g.Droplocations {
			if loc.SetSprite != nil && loc.SetSprite.Contains(cursorX, cursorY) {
				g.Stroke = ui.NewStroke(cursorX, cursorY, loc.SetSprite)
				loc.SetSprite = nil
				g.Game.SetSlot(loc.Index, "")
				g.RecalculateMatches()
				break
			}
		}
	}

	if g.Stroke != nil {
		cursorX, cursorY := ui.AdjustedCursorPosition()
		g.Stroke.Update(cursorX, cursorY)

		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
			dropTaken := false
			for _, loc := range g.Droplocations {
				if loc.Contains(cursorX, cursorY) {
					g.Stroke.DraggingObject.MoveTo(loc.X+float64(loc.W)/2-ui.SetSpriteWidth/2, loc.Y+float64(loc.H)/2-ui.SetSpriteHeight/2)

					if loc.SetSprite != nil {
						// If there's already a sprite in the drop location, we need to remove it first
						g.Setsprites = append(g.Setsprites, loc.SetSprite)
						slices.SortFunc(g.Setsprites, func(a, b *ui.SetSprite) int {
							return strings.Compare(a.SpriteName, b.SpriteName)
						})
					}

					loc.SetSprite = g.Stroke.DraggingObject.(*ui.SetSprite)
					g.Game.SetSlot(loc.Index, g.Stroke.DraggingObject.(*ui.SetSprite).SpriteName)
					dropTaken = true
					break
				}
			}
			if !dropTaken {
				g.Setsprites = append(g.Setsprites, g.Stroke.DraggingObject.(*ui.SetSprite))
				slices.SortFunc(g.Setsprites, func(a, b *ui.SetSprite) int {
					return strings.Compare(a.SpriteName, b.SpriteName)
				})
			}
			for i, loc := range g.Setsprites {
				loc.X = SET_START_X
				loc.Y = SET_START_Y + float64(i)*50
			}
			g.RecalculateMatches()

			g.Stroke.DraggingObject = nil
			g.Stroke.Release()
			g.Stroke = nil
		}
	}
}
func (g *GameScene) RecalculateMatches() {
	for i := range 6 {
		for j, m := range g.Game.Matches[i] {
			if m {
				g.Colors[i][j] = color.RGBA{90, 190, 90, 255} // Green for matches
			} else {
				g.Colors[i][j] = color.RGBA{0, 0, 0, 255}
			}
		}
	}
}

func (g *GameScene) OnSwitch() {
}

func (g *GameScene) SetSceneManager(sm *SceneManager) {
	g.SceneManager = sm
}
