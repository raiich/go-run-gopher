package animation

import (
	"bytes"
	"fmt"
	"image"
	"log/slog"
	"reflect"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/raiich/go-run-gopher/resources/ingame"
	"github.com/raiich/kazura/must"
	"github.com/raiich/kazura/state"
	"github.com/raiich/kazura/state/graph"
)

var gopherImages = loadImages()

// loadImages loads and returns the three gopher animation images
func loadImages() []*ebiten.Image {
	gopherImg1, _, err := image.Decode(bytes.NewReader(ingame.CharaImageBytes1))
	if err != nil {
		panic(err)
	}

	gopherImg2, _, err := image.Decode(bytes.NewReader(ingame.CharaImageBytes2))
	if err != nil {
		panic(err)
	}

	gopherImg3, _, err := image.Decode(bytes.NewReader(ingame.CharaImageBytes3))
	if err != nil {
		panic(err)
	}

	return []*ebiten.Image{
		ebiten.NewImageFromImage(gopherImg1),
		ebiten.NewImageFromImage(gopherImg2),
		ebiten.NewImageFromImage(gopherImg3),
	}
}

// on is a helper function to simplify state.On[animationState, animationEvent] calls
func on[T animationEvent](from, to animationState) state.Edge[animationState] {
	return state.On[animationState, T](from, to)
}

func check(g *graph.Graph[animationState, reflect.Type], err error) *graph.Graph[animationState, reflect.Type] {
	must.NoError(err)
	slog.Debug(fmt.Sprintf("animation state: %v", graph.Dump(g)))
	return g
}
