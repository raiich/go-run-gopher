// Package gopher implements the gopher character with movement, animation, and bounce mechanics.
package gopher

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/raiich/go-run-gopher/internal/game/gopher/animation"
	"github.com/raiich/go-run-gopher/internal/game/ui"
	"github.com/raiich/kazura/must"
	"github.com/raiich/kazura/state"
	"github.com/raiich/kazura/task"
)

// Bounds represents the field boundaries for gopher movement
type Bounds struct {
	Width  float64
	Height float64
}

// Config holds the configuration parameters for Gopher behavior and physics
type Config struct {
	Size          float64
	SpeedMax      float64
	KickSpeed     float64 // speedX increment per kick
	StopThreshold float64
	HalvingFrames int
	FieldBounds   Bounds
	InitialX      float64
	InitialY      float64
}

// Data holds the gopher's runtime state including position, speed, and animation
type Data struct {
	Config
	animation          *animation.Animation
	dispatcher         task.Dispatcher
	X                  float64
	Y                  float64
	Direction          int // 1 = right, -1 = left
	speedX             float64
	speedY             float64
	HalvingTimer       int
	lastBounceStrength float64 // Stores the bounce speedY for stun duration calculation
	OnRecover          func()  // Called when gopher recovers from stun
}

// Gopher represents the gopher character, managing its state machine and behavior
type Gopher struct {
	*Data
	machine *state.Machine[gopherState, *Data]
}

// New creates and initializes a new gopher character with animation and state machine.
// The gopher is positioned at the left side of the screen facing right.
func New(dispatcher task.Dispatcher, screenWidth, screenHeight float64) *Gopher {
	// Create gopher config with screen bounds
	config := Config{
		Size:          defaultSize,
		SpeedMax:      defaultSpeedMax,
		KickSpeed:     defaultKickSpeed,
		StopThreshold: defaultStopThreshold,
		HalvingFrames: defaultHalvingFrames,
		InitialX:      50.0,
		InitialY:      screenHeight*9/20 - defaultSize/2,
		FieldBounds: Bounds{
			Width:  screenWidth,
			Height: screenHeight,
		},
	}

	// Create gopher data with animation and bouncing
	data := &Data{
		Config:     config,
		animation:  animation.New(dispatcher),
		dispatcher: dispatcher,
		X:          config.InitialX,
		Y:          config.InitialY,
		Direction:  1, // Start facing right
	}

	machine := state.NewMachine[gopherState, *Data](stateGraph, data)
	must.NoError(machine.Launch())

	// Create gopher instance
	return &Gopher{
		Data:    data,
		machine: machine,
	}
}

// AdjustSpeed sets the gopher's speedX with threshold and cap handling
// All speedX changes must go through this method
//
// Behavior:
// - Negative values or values below StopThreshold are treated as 0
// - Values above SpeedMax are capped at SpeedMax
// - Setting speedX to 0 will reset HalvingTimer and stop animation
// - Setting speedX > 0 will start animation if currently stopped
func (d *Data) AdjustSpeed(speed float64) {
	previousSpeed := d.speedX

	// Check if speedX is below threshold (including negative values)
	// Only apply threshold when speedX is decreasing (halving), not when accelerating
	if speed < d.StopThreshold && speed < previousSpeed {
		d.speedX = 0
		d.HalvingTimer = 0
		// Stop animation if transitioning from moving to stopped
		if previousSpeed > 0 {
			must.NoError(d.animation.Stop())
		}
		return
	}
	// Cap at maximum speedX
	if speed > d.SpeedMax {
		d.speedX = d.SpeedMax
	} else {
		d.speedX = speed
	}

	// Start animation if transitioning from stopped to moving
	if previousSpeed == 0 && d.speedX > 0 {
		must.NoError(d.animation.Start())
	}
}

// GetSpeed returns the current speedX
func (d *Data) GetSpeed() float64 {
	return d.speedX
}

// Update updates the gopher's position and speedX
func (g *Gopher) Update() {
	must.Must(g.machine.CurrentState()).update(g)
}

// Kick accelerates the gopher by incrementing speedX
func (d *Data) Kick() {
	// Add kick speedX to current speedX (AdjustSpeed will cap at maximum and start animation if needed)
	d.AdjustSpeed(d.speedX + d.KickSpeed)
	// Reset halving timer
	d.HalvingTimer = d.HalvingFrames
}

// StartBounce initiates a bounce animation in the given direction
// direction: -1 = left (bouncing away from left wall), 1 = right (bouncing away from right wall)
func (g *Gopher) StartBounce() error {
	// Start bounce animation in the new direction (which is opposite to the wall)
	return g.machine.Trigger(crashEvent{})
}

// HandleInput handles user input for the gopher
func (g *Gopher) HandleInput(button ui.ButtonType, justPressed bool) {
	must.Must(g.machine.CurrentState()).handleInput(g, button, justPressed)
}

// Draw draws the gopher on the screen with bounce offsets if active
func (g *Gopher) Draw(screen *ebiten.Image) {
	gopherOpts := &ebiten.DrawImageOptions{}

	currentGopherImage := g.animation.CurrentImage()
	gopherBounds := currentGopherImage.Bounds()
	gopherScale := g.Size / float64(gopherBounds.Dx())
	gopherOpts.GeoM.Scale(gopherScale*float64(g.Direction), gopherScale)
	if g.Direction == -1 {
		gopherOpts.GeoM.Translate(g.X+g.Size, g.Y)
	} else {
		gopherOpts.GeoM.Translate(g.X, g.Y)
	}
	screen.DrawImage(currentGopherImage, gopherOpts)
}
