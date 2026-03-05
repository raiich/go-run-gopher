// Package game implements the main game scene with gameplay logic and state management.
package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/raiich/go-run-gopher/internal/credits"
	"github.com/raiich/go-run-gopher/internal/game/audio"
	"github.com/raiich/go-run-gopher/internal/game/effect"
	"github.com/raiich/go-run-gopher/internal/game/gopher"
	"github.com/raiich/go-run-gopher/internal/game/pregame"
	"github.com/raiich/go-run-gopher/internal/game/ui"
	"github.com/raiich/go-run-gopher/lib/scene"
	"github.com/raiich/kazura/must"
	"github.com/raiich/kazura/state"
)

// Scene represents the main game scene with game state and rendering
type Scene struct {
	*sceneData
	machine *state.Machine[sceneState, *sceneData]
}

// Update updates the scene state, handling input and delegating to the current game state
func (s *Scene) Update() error {
	s.effect.Update()

	button, justPressed := s.ui.GetButtonState()
	must.Must(s.machine.CurrentState()).handleInput(s, button, justPressed)

	// Delegate state-specific update to current state
	return must.Must(s.machine.CurrentState()).update(s)
}

// Draw renders the scene including field, gopher, UI, and state-specific elements
func (s *Scene) Draw(screen *ebiten.Image) {
	// draw field background (programmatic gymnasium floor)
	floorColor := color.RGBA{R: 140, G: 100, B: 70, A: 255}
	goalZoneColor := color.RGBA{R: 110, G: 80, B: 55, A: 255}
	screen.Fill(floorColor)
	// Goal zones: areas outside the boundary lines
	vector.DrawFilledRect(screen, 0, 0, float32(boundaryMargin), float32(ScreenHeight), goalZoneColor, false)
	vector.DrawFilledRect(screen, float32(ScreenWidth-boundaryMargin), 0, float32(boundaryMargin), float32(ScreenHeight), goalZoneColor, false)

	// draw boundary margin lines (part of the field, like gymnasium floor markings)
	s.drawBoundaryLines(screen)

	// draw Gopher (always visible, handles its own bounce offsets)
	s.gopher.Draw(screen)

	// draw UI and state-specific UI
	must.Must(s.machine.CurrentState()).draw(s, screen)

	s.effect.Draw(screen)
	s.overlay.Draw(screen)
}

// drawBoundaryLines draws vertical boundary lines with top/bottom margins.
func (s *Scene) drawBoundaryLines(screen *ebiten.Image) {
	lineColor := color.RGBA{R: 255, G: 255, B: 255, A: 200}
	lineWidth := float32(4)
	margin := float32(ScreenHeight) / 4

	leftX := float32(boundaryMargin) - lineWidth/2
	rightX := float32(ScreenWidth-boundaryMargin) - lineWidth/2
	h := float32(ScreenHeight) - margin*2

	vector.DrawFilledRect(screen, leftX, margin, lineWidth, h, lineColor, false)
	vector.DrawFilledRect(screen, rightX, margin, lineWidth, h, lineColor, false)
}

// Layout returns the logical screen size for the game
func (s *Scene) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

// showCredits pushes the credits scene onto the scene stack
func (s *Scene) showCredits() {
	creditsScene := credits.NewScene(s.scenes)
	must.NoError(s.scenes.PushScene(creditsScene))
}

// NewScene creates a new game scene with all necessary components initialized
func NewScene(scenes *scene.Manager) *Scene {
	dispatcher := scenes.Dispatcher()

	// Create screen data
	g := gopher.New(dispatcher, ScreenWidth, ScreenHeight)
	eff := effect.New(ui.TextFace36, ScreenWidth, ScreenHeight)

	data := &sceneData{
		scenes:      scenes,
		dispatcher:  dispatcher,
		ui:          ui.NewController(ScreenWidth, ScreenHeight),
		gopher:      g,
		effect:      eff,
		overlay:     pregame.NewOverlay(dispatcher, ui.TextFace80, ScreenWidth, ScreenHeight),
		scalePlayer: audio.NewScalePlayer(),
		uiConfig:    defaultUIConfig(),
	}

	g.OnRecover = func() { eff.Show(effect.KindRecover) }

	// Create state machine
	machine := state.NewMachine[sceneState, *sceneData](stateGraph, data)
	must.NoError(machine.Launch())

	return &Scene{
		sceneData: data,
		machine:   machine,
	}
}
