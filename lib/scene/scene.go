// Package scene provides scene management functionality for Ebiten games.
package scene

import "github.com/hajimehoshi/ebiten/v2"

// Scene represents a game scene (screen) that can be updated, drawn, and laid out.
// This interface matches ebiten.Game to allow scenes to be managed uniformly.
type Scene interface {
	// Update updates the scene logic. Called every tick.
	Update() error
	// Draw draws the scene to the screen. Called every frame.
	Draw(screen *ebiten.Image)
	// Layout determines the game's logical screen size.
	Layout(outsideWidth, outsideHeight int) (int, int)
}
