package pregame

import (
	"fmt"
	"log/slog"
	"reflect"

	"github.com/raiich/kazura/must"
	"github.com/raiich/kazura/state"
	"github.com/raiich/kazura/state/graph"
)

func check(g *graph.Graph[pregameState, reflect.Type], err error) *graph.Graph[pregameState, reflect.Type] {
	must.NoError(err)
	slog.Debug(fmt.Sprintf("state: %v", graph.Dump(g)))
	return g
}

// on is a helper function to simplify state.On[pregameState, pregameEvent] calls
func on[T pregameEvent](from, to pregameState) state.Edge[pregameState] {
	return state.On[pregameState, T](from, to)
}
