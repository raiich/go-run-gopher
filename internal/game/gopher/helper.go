package gopher

import (
	"fmt"
	"log/slog"
	"reflect"

	"github.com/raiich/kazura/must"
	"github.com/raiich/kazura/state"
	"github.com/raiich/kazura/state/graph"
)

const (
	defaultSize          = 100.0
	defaultSpeedMax      = 10.0
	defaultKickSpeed     = 1.0 // speedX increment per kick
	defaultStopThreshold = 2.0
	defaultHalvingFrames = 60
)

func check(g *graph.Graph[gopherState, reflect.Type], err error) *graph.Graph[gopherState, reflect.Type] {
	must.NoError(err)
	slog.Debug(fmt.Sprintf("bouncing state: %v", graph.Dump(g)))
	return g
}

// on is a helper function to simplify state.On[gopherState, gopherEvent] calls
func on[T gopherEvent](from, to gopherState) state.Edge[gopherState] {
	return state.On[gopherState, T](from, to)
}
