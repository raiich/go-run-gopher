// Package ui provides UI components for the game including buttons and controllers.
package ui

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/raiich/go-run-gopher/lib/ui"
)

// ButtonType represents the type of button that was pressed
type ButtonType int

const (
	// ButtonTypeNone indicates no button was pressed
	ButtonTypeNone ButtonType = iota
	// ButtonTypeAbout indicates the About button was pressed
	ButtonTypeAbout
	// ButtonTypeAction indicates the main Action button was pressed
	ButtonTypeAction
)

// Controller manages all UI elements in the game scene including buttons and fonts
type Controller struct {
	AboutButton  *ui.RoundedRectButton
	ActionButton *ui.CircleButton
}

// NewController creates a new UI controller
func NewController(screenWidth, screenHeight int) *Controller {
	return &Controller{
		AboutButton: &ui.RoundedRectButton{
			X:      20,
			Y:      20,
			W:      125,
			H:      40,
			Radius: 10,
			Text: ui.Text{
				Text:    "Credits",
				Face:    TextFace16,
				Color:   color.RGBA{R: 255, G: 255, B: 255, A: 255},
				BGColor: color.RGBA{R: 0, G: 0, B: 0, A: 0},
			},
			DefaultColor: color.RGBA{R: 100, G: 100, B: 200, A: 255},
			PressedColor: color.RGBA{R: 150, G: 150, B: 250, A: 255},
		},
		ActionButton: &ui.CircleButton{
			X: float64(screenWidth) - 120,
			Y: float64(screenHeight) - 90,
			R: 70,
			Text: ui.Text{
				Text:    "go",
				Face:    TextFace36,
				Color:   color.RGBA{R: 255, G: 255, B: 255, A: 255},
				BGColor: color.RGBA{R: 0, G: 0, B: 0, A: 0},
			},
			DefaultColor: color.RGBA{R: 200, G: 100, B: 100, A: 255},
			PressedColor: color.RGBA{R: 250, G: 150, B: 150, A: 255},
		},
	}
}

func (c *Controller) GetButtonState() (ButtonType, bool) {
	// Handle touch input first (mobile-first approach)
	touchIDs := ebiten.AppendTouchIDs(nil)
	if len(touchIDs) > 0 {
		// Only process the first touch to avoid multiple simultaneous actions
		x, y := ebiten.TouchPosition(touchIDs[0])
		return c.HandleInput(float64(x), float64(y))
	} else {
		// Handle mouse input only if no touch is active
		mousePressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
		if mousePressed {
			x, y := ebiten.CursorPosition()
			return c.HandleInput(float64(x), float64(y))
		} else {
			// Reset button states when no input is active
			c.ResetButtonStates()
		}
	}

	return ButtonTypeNone, false
}

// HandleInput handles input for UI elements and returns which button was pressed and if it was just pressed
// Returns: (ButtonType, justPressed) where justPressed is true if the button was just pressed this frame
func (c *Controller) HandleInput(x, y float64) (ButtonType, bool) {
	var pressedButton ButtonType
	var justPressed bool

	if c.AboutButton.Contains(x, y) {
		if !c.AboutButton.IsPressed() {
			c.AboutButton.SetPressed(true)
			justPressed = true
		}
		if c.AboutButton.IsPressed() {
			pressedButton = ButtonTypeAbout
		}
	}

	if c.ActionButton.Contains(x, y) {
		if !c.ActionButton.IsPressed() {
			c.ActionButton.SetPressed(true)
			justPressed = true
		}
		if c.ActionButton.IsPressed() {
			pressedButton = ButtonTypeAction
		}
	}

	return pressedButton, justPressed
}

// ResetButtonStates resets all button states
func (c *Controller) ResetButtonStates() {
	if c.AboutButton.IsPressed() {
		c.AboutButton.SetPressed(false)
	}
	if c.ActionButton.IsPressed() {
		c.ActionButton.SetPressed(false)
	}
}

// Draw draws UI elements based on visibility settings
func (c *Controller) Draw(screen *ebiten.Image, button ButtonType) {
	switch button {
	case ButtonTypeAbout:
		c.AboutButton.Draw(screen, TextFace16)
	case ButtonTypeAction:
		c.ActionButton.Draw(screen, TextFace36)
	default:
		panic(fmt.Errorf("unknown button type: %v", button))
	}
}
