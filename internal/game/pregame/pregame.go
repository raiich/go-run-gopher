package pregame

import (
	"fmt"
	"image"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	libui "github.com/raiich/go-run-gopher/lib/ui"
	"github.com/raiich/kazura/must"
	"github.com/raiich/kazura/state"
	"github.com/raiich/kazura/task"
)

// stateGraph defines the pregame state transitions
var stateGraph = check(state.NewGraph[pregameState](
	noneState{},
	on[tickEvent](noneState{}, displayState{}),
	on[tickEvent](displayState{}, displayState{}),
	on[stopEvent](displayState{}, noneState{}),
))

type pregameEvent interface {
	state.Event
}

type pregameState interface {
	Entry(machine *state.EntryMachine[*pregameData], event state.Event) state.Command
	Draw(data *pregameData, screen *ebiten.Image)
}

type tickEvent struct {
	nextCount time.Duration
	callback  func()
}

type stopEvent struct{}

type noneState struct{}

func (s noneState) Entry(machine *state.EntryMachine[*pregameData], event state.Event) state.Command {
	// nothing to do
	return nil
}

func (s noneState) Draw(data *pregameData, screen *ebiten.Image) {
	// nothing to do
}

// displayState represents the pregame display state before game starts
type displayState struct{}

func (s displayState) Entry(machine *state.EntryMachine[*pregameData], event state.Event) state.Command {
	value := machine.Value()
	e := event.(tickEvent)

	if e.nextCount > 0 {
		value.text = fmt.Sprintf("%v", int(e.nextCount/time.Second))
		must.NoError(machine.AfterFunc(value.dispatcher, 300*time.Millisecond, func(machine *state.AfterFuncMachine[*pregameData]) {
			must.NoError(machine.Trigger(tickEvent{
				nextCount: e.nextCount - 1*time.Second,
				callback:  e.callback,
			}))
		}))
	} else {
		value.text = "GO"
		e.callback()
		must.NoError(machine.AfterFunc(value.dispatcher, 300*time.Millisecond, func(machine *state.AfterFuncMachine[*pregameData]) {
			must.NoError(machine.Trigger(stopEvent{}))
		}))
	}
	return nil
}

func (s displayState) Draw(data *pregameData, screen *ebiten.Image) {
	// Create large semi-transparent text
	textUI := libui.Text{
		Text:    data.text,
		Face:    data.textFace,
		Color:   color.RGBA{R: 255, G: 255, B: 255, A: 180},
		BGColor: color.RGBA{R: 0, G: 0, B: 0, A: 0},
	}
	textImage := textUI.Image()
	textWidth := textImage.Rectangle().Dx()
	textHeight := textImage.Rectangle().Dy()

	// Center on screen
	textX := data.screenWidth/2 - textWidth/2
	textY := data.screenHeight/2 - textHeight/2

	textImage.Draw(screen, image.Point{X: textX, Y: textY})
}

type pregameData struct {
	dispatcher   task.Dispatcher
	textFace     *text.GoTextFace
	text         string
	screenWidth  int
	screenHeight int
}
