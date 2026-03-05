package game

import "github.com/raiich/go-run-gopher/internal/screen"

const (
	Title = "go run ./gopher"

	ScreenWidth  = screen.Width
	ScreenHeight = screen.Height

	// Gameplay constants
	// boundaryMargin defines the goal zone width (like a shuttle run line)
	// Used for both visual line position and hit detection
	boundaryMargin = 130

	// Stage progression
	lapsPerStage = 7
)
