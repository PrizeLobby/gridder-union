package main

import (
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/prizelobby/union-gridder/config"
	"github.com/prizelobby/union-gridder/core"
	"github.com/prizelobby/union-gridder/res"
	"github.com/prizelobby/union-gridder/scene"
	"github.com/prizelobby/union-gridder/ui"
	"github.com/tinne26/etxt"
)

const GAME_WIDTH = 960
const GAME_HEIGHT = 720

const SAMPLE_RATE = 48000

type GameState int

const (
	MENU GameState = iota
	PLAYING
	CREDITS
)

type EbitenGame struct {
	ScaledScreen *ui.ScaledScreen
	gameState    GameState
	SceneManager *scene.SceneManager
}

func (g *EbitenGame) Update() error {
	if config.DEBUG {
		if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
			os.Exit(0)
		}
	}

	g.SceneManager.Update()
	return nil
}

func (g *EbitenGame) Draw(screen *ebiten.Image) {
	g.ScaledScreen.SetTarget(screen)

	//msg := fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS())
	//g.ScaledScreen.DebugPrint(msg)

	g.SceneManager.Draw(g.ScaledScreen)
}

func (g *EbitenGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	panic("use Ebitengine >=v2.5.0")
}

func (g *EbitenGame) LayoutF(outsideWidth, outsideHeight float64) (screenWidth, screenHeight float64) {
	scale := ebiten.Monitor().DeviceScaleFactor()
	canvasWidth := GAME_WIDTH * scale
	canvasHeight := GAME_HEIGHT * scale
	return canvasWidth, canvasHeight
}

func main() {
	game := core.NewGame()

	// create a new text renderer and configure it
	txtRenderer := etxt.NewRenderer()
	txtRenderer.Utils().SetCache8MiB()
	txtRenderer.SetFont(res.GetFont("Roboto-Medium"))
	txtRenderer.SetAlign(etxt.HorzCenter | etxt.VertCenter)
	txtRenderer.SetSize(64)

	scaledScreen := ui.NewScaledScreen(txtRenderer)

	g := &EbitenGame{
		ScaledScreen: scaledScreen,
		gameState:    MENU,
	}
	sm := scene.NewSceneManager()
	gameScene := scene.NewGameScene(game)
	sm.AddScene("game", gameScene)
	g.SceneManager = sm
	g.SceneManager.SwitchToScene("game")

	ebiten.SetWindowSize(GAME_WIDTH, GAME_HEIGHT)
	ebiten.SetWindowTitle("Gridder Union")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
