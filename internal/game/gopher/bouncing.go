package gopher

import (
	"math"
	"time"

	"github.com/raiich/go-run-gopher/internal/game/ui"
	"github.com/raiich/kazura/must"
	"github.com/raiich/kazura/state"
)

// Bounce parameters for speed-dependent bounce calculations
const (
	// Base values: current fixed values used as reference for medium speed
	bounceBaseSpeed  = 7.5   // Reference collision speed (approximately half of SpeedMax)
	bounceBaseSpeedX = 3.0   // Base horizontal bounce speed
	bounceBaseSpeedY = -12.0 // Base initial upward velocity

	// Minimum bounce values (ensure bounce even at low speeds)
	bounceMinSpeedX = 1.5
	bounceMinSpeedY = -8.0

	// Maximum bounce values (cap bounce at high speeds)
	bounceMaxSpeedX = 6.0
	bounceMaxSpeedY = -24.0
)

// stateGraph defines the bouncing state transitions
var stateGraph = check(state.NewGraph[gopherState](
	normalState{},
	on[crashEvent](normalState{}, boundState{}),
	on[landEvent](boundState{}, stunState{}),
	on[recoverEvent](stunState{}, normalState{}),
))

type gopherEvent interface {
	state.Event
}

type gopherState interface {
	Entry(machine *state.EntryMachine[*Data], event state.Event) state.Command
	handleInput(g *Gopher, button ui.ButtonType, justPressed bool)
	update(g *Gopher)
}

type crashEvent struct {
}

type landEvent struct{}

type recoverEvent struct{}

type normalState struct{}

func (s normalState) Entry(machine *state.EntryMachine[*Data], event state.Event) state.Command {
	g := machine.Value()
	g.speedX = 0
	g.Y = g.InitialY

	if _, ok := event.(recoverEvent); ok && g.OnRecover != nil {
		g.OnRecover()
	}
	return nil
}

func (s normalState) handleInput(g *Gopher, button ui.ButtonType, justPressed bool) {
	// Run button - accelerate gopher
	if button == ui.ButtonTypeAction && justPressed {
		g.Kick()
	}
}

func (s normalState) update(g *Gopher) {
	// Handle speedX halving timer
	if g.HalvingTimer > 0 {
		g.HalvingTimer--
		if g.HalvingTimer == 0 {
			// Halve the speedX after timer expires
			g.AdjustSpeed(g.speedX / 2)
			// Set timer for next halving if speedX is still above threshold (AdjustSpeed handles 0 case)
			if g.speedX > 0 {
				g.HalvingTimer = g.HalvingFrames
			}
		}
	}

	if g.speedX == 0 {
		return
	}

	// Only move if speedX > 0
	g.X += g.speedX * float64(g.Direction)

	// Boundary checking - clamp position only (don't stop here, let game state handle it)
	if g.X <= 0 {
		g.X = 0
	} else if g.X >= g.FieldBounds.Width-g.Size {
		g.X = g.FieldBounds.Width - g.Size
	}
}

// boundState represents the bouncing animation state
type boundState struct{}

func (s boundState) Entry(machine *state.EntryMachine[*Data], event state.Event) state.Command {
	g := machine.Value()

	// Calculate bounce strength based on collision speed (before stopping)
	collisionSpeed := g.speedX
	// Use linear scaling: factor = collisionSpeed / baseSpeed
	// This provides stronger bounce at high speeds
	factor := collisionSpeed / bounceBaseSpeed

	// Calculate bounce speeds with min/max clamping
	bounceSpeedX := math.Max(bounceMinSpeedX, math.Min(bounceMaxSpeedX, bounceBaseSpeedX*factor))
	bounceSpeedY := math.Max(bounceMaxSpeedY, math.Min(bounceMinSpeedY, bounceBaseSpeedY*factor))

	// Stop gopher movement (AdjustSpeed will reset HalvingTimer when setting to 0)
	g.AdjustSpeed(0)

	// Reverse direction - gopher will run in opposite direction after bounce
	g.Direction = -g.Direction

	g.speedX = bounceSpeedX
	g.speedY = bounceSpeedY
	g.lastBounceStrength = bounceSpeedY // Store bounce strength for stun duration

	must.NoError(machine.OnExit(func(state.Event) *state.Guarded {
		// Reset bouncing offsets
		g.Y = g.InitialY

		return nil
	}))
	return nil
}

func (s boundState) handleInput(g *Gopher, button ui.ButtonType, justPressed bool) {
	// ignore input
}

func (s boundState) update(g *Gopher) {
	// Apply gravity
	g.speedY += 0.6
	g.Y += g.speedY

	// Update horizontal position (bounce backwards) - directly modify X
	g.X += float64(g.Direction) * g.speedX

	// Check if landed (returned to ground level)
	if g.Y >= g.InitialY {
		g.Y = g.InitialY
		must.NoError(g.machine.Trigger(landEvent{}))
	}
}

// stunState represents the stun/recovery state after landing
type stunState struct{}

func (s stunState) Entry(machine *state.EntryMachine[*Data], event state.Event) state.Command {
	g := machine.Value()

	// Calculate stun duration based on bounce strength (lastBounceStrength)
	// Weak bounce (speedY > -10): ~250ms
	// Medium bounce (speedY -10 to -17): 250ms to 500ms (linear interpolation)
	// Strong bounce (speedY < -17): 500ms to 700ms (linear interpolation)
	var stunDuration time.Duration
	bounceStrength := -g.lastBounceStrength // Convert to positive for easier calculation

	if bounceStrength < 10.0 {
		// Weak bounce - short stun
		stunDuration = 250 * time.Millisecond
	} else if bounceStrength < 17.0 {
		// Medium bounce - interpolate between 250ms and 500ms
		t := (bounceStrength - 10.0) / 7.0 // normalize to 0-1
		stunDuration = time.Duration(250+t*250) * time.Millisecond
	} else {
		// Strong bounce - interpolate between 500ms and 700ms
		t := math.Min(1.0, (bounceStrength-17.0)/7.0) // normalize to 0-1, cap at 1
		stunDuration = time.Duration(500+t*200) * time.Millisecond
	}

	must.NoError(machine.AfterFunc(g.dispatcher, stunDuration, func(machine *state.AfterFuncMachine[*Data]) {
		must.NoError(machine.Trigger(recoverEvent{}))
	}))
	return nil
}

func (s stunState) handleInput(g *Gopher, button ui.ButtonType, justPressed bool) {
	// ignore input
}

func (s stunState) update(g *Gopher) {
	// nothing to do during stun
}
