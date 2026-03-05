package game

import (
	"fmt"
	"log/slog"
	"reflect"

	"github.com/raiich/kazura/must"
	"github.com/raiich/kazura/state"
	"github.com/raiich/kazura/state/graph"
)

// on is a helper function to simplify state.On[sceneState, sceneEvent] calls
func on[T sceneEvent](from, to sceneState) state.Edge[sceneState] {
	return state.On[sceneState, T](from, to)
}

func check(g *graph.Graph[sceneState, reflect.Type], err error) *graph.Graph[sceneState, reflect.Type] {
	must.NoError(err)
	slog.Debug(fmt.Sprintf("state: %v", graph.Dump(g)))
	return g
}
