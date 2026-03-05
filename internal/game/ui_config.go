package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
	ui2 "github.com/raiich/go-run-gopher/internal/game/ui"
)

// UIConfig holds layout and style configuration for in-game HUD elements.
type UIConfig struct {
	Stats  StatsConfig
	Status StatusConfig
	Goal   GoalConfig
	Timer  TimerConfig
	Result ResultConfig
}

// StatsConfig configures the stage/laps display.
type StatsConfig struct {
	Y     int
	Face  *text.GoTextFace
	Color color.RGBA
}

// StatusConfig configures the test status display.
type StatusConfig struct {
	Y        int
	Face     *text.GoTextFace
	RunColor color.RGBA
}

// GoalConfig configures the GOAL marker display.
type GoalConfig struct {
	Y     int
	Face  *text.GoTextFace
	Color color.RGBA
}

// TimerConfig configures the scale count (timer) display.
type TimerConfig struct {
	Y               int
	Face            *text.GoTextFace
	ColorNormal     color.RGBA
	ColorWarn       color.RGBA
	ColorDanger     color.RGBA
	WarnThreshold   int
	DangerThreshold int
}

// ResultConfig configures the game over message display.
type ResultConfig struct {
	Face  *text.GoTextFace
	Color color.RGBA
}

// defaultUIConfig returns the default UI configuration.
func defaultUIConfig() UIConfig {
	return UIConfig{
		Stats: StatsConfig{
			Y:     20,
			Face:  ui2.TextFace16,
			Color: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		},
		Status: StatusConfig{
			Y:        40,
			Face:     ui2.TextFace14,
			RunColor: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		},
		Goal: GoalConfig{
			Y:     60,
			Face:  ui2.TextFace24,
			Color: color.RGBA{R: 180, G: 255, B: 100, A: 255},
		},
		Timer: TimerConfig{
			Y:               8,
			Face:            ui2.TextFace36,
			ColorNormal:     color.RGBA{R: 255, G: 255, B: 255, A: 255},
			ColorWarn:       color.RGBA{R: 255, G: 255, B: 0, A: 255},
			ColorDanger:     color.RGBA{R: 255, G: 0, B: 0, A: 255},
			WarnThreshold:   4,
			DangerThreshold: 2,
		},
		Result: ResultConfig{
			Face:  ui2.TextFace16,
			Color: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		},
	}
}
