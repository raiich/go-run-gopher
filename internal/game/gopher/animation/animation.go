package animation

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/raiich/kazura/must"
	"github.com/raiich/kazura/state"
)

// animationData holds the animation state data
type animationData struct {
	dispatcher state.Dispatcher
	images     []*ebiten.Image
	current    int
}

// Animation represents the gopher animation state machine
type Animation struct {
	*animationData
	machine *state.Machine[animationState, *animationData]
}

// New creates a new animation state machine
func New(dispatcher state.Dispatcher) *Animation {
	data := &animationData{
		dispatcher: dispatcher,
		images:     gopherImages,
		current:    0,
	}

	machine := state.NewMachine(stateGraph, data)
	must.NoError(machine.Launch())

	return &Animation{
		animationData: data,
		machine:       machine,
	}
}

// Start starts the animation
func (a *Animation) Start() error {
	return a.machine.Trigger(startEvent{})
}

// Stop stops the animation and resets to frame 0
func (a *Animation) Stop() error {
	return a.machine.Trigger(stopEvent{})
}

func (a *Animation) CurrentImage() *ebiten.Image {
	return a.images[a.animationData.current]
}

// stateGraph defines the animation state transitions
var stateGraph = check(state.NewGraph[animationState](
	stoppedState{},
	on[startEvent](stoppedState{}, runningState{}),
	on[tickEvent](runningState{}, runningState{}),
	on[stopEvent](nil, stoppedState{}),
))

type animationEvent interface {
	state.Event
}

type animationState interface {
	Entry(machine *state.EntryMachine[*animationData], event state.Event)
}

// Events

type startEvent struct{}
type stopEvent struct{}
type tickEvent struct{}

// States

type stoppedState struct{}

func (s stoppedState) Entry(machine *state.EntryMachine[*animationData], event state.Event) {
	data := machine.Value()
	data.current = 0
}

type runningState struct{}

func (s runningState) Entry(machine *state.EntryMachine[*animationData], event state.Event) {
	data := machine.Value()

	// tickEvent: increment frame and wrap around
	data.current = (data.current + 1) % len(data.images)

	// Schedule next frame transition
	machine.AfterFunc(data.dispatcher, 100*time.Millisecond, func(machine *state.AfterFuncMachine[*animationData]) {
		must.NoError(machine.Trigger(tickEvent{}))
	})
}
