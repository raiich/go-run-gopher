package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/raiich/go-run-gopher/internal/credits"
	"github.com/raiich/go-run-gopher/internal/game"
	"github.com/raiich/go-run-gopher/lib/scene"
	"github.com/raiich/kazura/must"
)

var (
	sceneFlag = flag.String("scene", "game", "Initial scene to load (game, credits)")
)

func main() {
	flag.Parse()

	ebiten.SetWindowTitle(game.Title)
	ebiten.SetWindowSize(game.ScreenWidth, game.ScreenHeight)

	sceneManager := scene.NewSceneManager()

	// Select initial scene based on flag
	var initialScene scene.Scene
	switch *sceneFlag {
	case "game":
		initialScene = game.NewScene(sceneManager)
	case "credits":
		initialScene = credits.NewScene(sceneManager)
	default:
		fmt.Fprintf(os.Stderr, "Unknown scene: %s\n", *sceneFlag)
		fmt.Fprintf(os.Stderr, "Available scenes: game, credits\n")
		os.Exit(1)
	}

	must.NoError(sceneManager.SetScene(initialScene))
	must.NoError(ebiten.RunGame(sceneManager))
}
