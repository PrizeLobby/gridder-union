package scene

import "github.com/prizelobby/union-gridder/ui"

type BaseScene struct {
	SceneManager *SceneManager
}

func (b *BaseScene) Update() {

}

func (b *BaseScene) Draw(screen *ui.ScaledScreen) {

}

func (b *BaseScene) OnSwitch() {

}

func (b *BaseScene) SetSceneManager(sm *SceneManager) {
	b.SceneManager = sm
}
