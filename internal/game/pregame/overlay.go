// Package pregame provides a pregame overlay system for displaying the 3, 2, 1, Start! sequence.
package pregame

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/raiich/kazura/must"
	"github.com/raiich/kazura/state"
)

// Overlay manages the pregame display and state machine
type Overlay struct {
	machine *state.Machine[pregameState, *pregameData]
}

// Start starts the pregame sequence for the specified duration and calls the callback when complete
func (o *Overlay) Start(duration time.Duration, callback func()) error {
	return o.machine.Trigger(tickEvent{
		nextCount: duration,
		callback:  callback,
	})
}

// Draw renders the pregame overlay to the screen
func (o *Overlay) Draw(screen *ebiten.Image) {
	must.Must(o.machine.CurrentState()).Draw(o.machine.Value(), screen)
}

// NewOverlay creates a new pregame overlay with the specified font and screen dimensions
func NewOverlay(dispatcher state.Dispatcher, face80 *text.GoTextFace, screenWidth, screenHeight int) *Overlay {
	data := &pregameData{
		dispatcher:   dispatcher,
		textFace:     face80,
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
	}
	machine := state.NewMachine[pregameState, *pregameData](stateGraph, data)
	must.NoError(machine.Launch())
	return &Overlay{
		machine: machine,
	}
}
