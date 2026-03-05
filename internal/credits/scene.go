// Package credits implements the credits/license display scene
package credits

import (
	"bytes"
	"image"
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/raiich/go-run-gopher/internal/game/ui"
	"github.com/raiich/go-run-gopher/internal/screen"
	"github.com/raiich/go-run-gopher/lib/scene"
	libui "github.com/raiich/go-run-gopher/lib/ui"
	"github.com/raiich/go-run-gopher/resources/credits"
)

const (
	// Screen dimensions
	screenWidth  = 800
	screenHeight = 450

	// Scroll settings
	scrollSpeed     = 20  // pixels per scroll tick
	lineHeight      = 18  // pixels per line
	padding         = 20  // padding from screen edges
	buttonHeight    = 40  // height of close button
	contentStartY   = 20  // Y position where content starts (top of screen)
	maxContentWidth = 760 // maximum width for text content
)

// Scene represents the credits display scene with scrollable license text
type Scene struct {
	scenes       *scene.Manager
	closeButton  *libui.RoundedRectButton
	scrollOffset int    // current scroll offset in pixels
	contentText  string // combined license text
	textHeight   int    // total height of rendered text

	// Touch drag state
	touchID       ebiten.TouchID
	touchActive   bool
	touchLastX    int
	touchLastY    int
	touchStartY   int // Y position at touch start (for tap detection)
	touchStartOff int // scrollOffset at touch start (for tap detection)
}

// Update handles input for scrolling and button presses
func (s *Scene) Update() error {
	// Handle mouse wheel scrolling
	_, dy := ebiten.Wheel()
	if dy != 0 {
		s.scrollOffset -= int(dy * scrollSpeed)
		s.clampScrollOffset()
	}

	// Handle touch input (drag to scroll + tap for button)
	touchIDs := ebiten.AppendTouchIDs(nil)
	if len(touchIDs) > 0 {
		if !s.touchActive {
			// Touch started
			s.touchID = touchIDs[0]
			x, y := ebiten.TouchPosition(s.touchID)
			s.touchActive = true
			s.touchLastX = x
			s.touchLastY = y
			s.touchStartY = y
			s.touchStartOff = s.scrollOffset
		} else {
			// Touch is ongoing — drag to scroll
			x, y := ebiten.TouchPosition(s.touchID)
			s.touchLastX = x
			deltaY := s.touchLastY - y
			if deltaY != 0 {
				s.scrollOffset += deltaY
				s.clampScrollOffset()
				s.touchLastY = y
			}
		}
	} else if s.touchActive {
		// Touch ended — check if it was a tap (minimal scroll)
		s.touchActive = false
		scrolledDistance := s.scrollOffset - s.touchStartOff
		if scrolledDistance < 0 {
			scrolledDistance = -scrolledDistance
		}
		if scrolledDistance < 10 {
			if s.closeButton.Contains(float64(s.touchLastX), float64(s.touchLastY)) {
				return s.scenes.PopScene()
			}
		}
	}

	// Handle mouse input for button (desktop)
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if s.closeButton.Contains(float64(x), float64(y)) {
			return s.scenes.PopScene()
		}
	}

	return nil
}

// Draw renders the credits scene with scrollable text
func (s *Scene) Draw(screen *ebiten.Image) {
	// Fill background with black
	screen.Fill(color.RGBA{R: 0, G: 0, B: 0, A: 255})

	// Draw close button at top right
	s.closeButton.Draw(screen, ui.TextFace16)

	// Draw scrollable license text (including # Fonts, # Go Gopher, etc.)
	contentY := contentStartY
	contentHeight := screenHeight - contentStartY - padding

	// Draw content text with scroll offset
	s.drawScrollableText(screen, padding, contentY, maxContentWidth, contentHeight)
}

// wrapLine wraps a single line to fit within maxWidth
func wrapLine(line string, face *text.GoTextFace, maxWidth float64) []string {
	if line == "" {
		return []string{""}
	}

	// Measure the line width
	lineWidth, _ := text.Measure(line, face, 0)
	if lineWidth <= maxWidth {
		return []string{line}
	}

	// Split by words
	words := strings.Fields(line)
	if len(words) == 0 {
		return []string{line}
	}

	var wrapped []string
	var currentLine string

	for _, word := range words {
		testLine := currentLine
		if testLine != "" {
			testLine += " "
		}
		testLine += word

		testWidth, _ := text.Measure(testLine, face, 0)
		if testWidth <= maxWidth {
			currentLine = testLine
		} else {
			// Current word doesn't fit, save current line and start new one
			if currentLine != "" {
				wrapped = append(wrapped, currentLine)
				currentLine = word
			} else {
				// Single word is too long, just add it anyway
				wrapped = append(wrapped, word)
				currentLine = ""
			}
		}
	}

	if currentLine != "" {
		wrapped = append(wrapped, currentLine)
	}

	return wrapped
}

// drawScrollableText renders the license text with scroll offset applied
func (s *Scene) drawScrollableText(screen *ebiten.Image, x, y, width, height int) {
	// Split text into lines and wrap long lines
	originalLines := strings.Split(s.contentText, "\n")
	var wrappedLines []string
	for _, line := range originalLines {
		wrapped := wrapLine(line, ui.TextFace14, float64(width))
		wrappedLines = append(wrappedLines, wrapped...)
	}

	textFace := ui.TextFace14

	// Calculate line height
	lineSpacing := textFace.Size * 1.2
	s.textHeight = int(float64(len(wrappedLines)) * lineSpacing)

	// Create clipping region
	subImage := screen.SubImage(image.Rect(x, y, x+width, y+height)).(*ebiten.Image)

	// Draw each line
	textColor := color.RGBA{R: 200, G: 200, B: 200, A: 255}
	currentY := float64(y - s.scrollOffset)

	for _, line := range wrappedLines {
		// Skip lines that are completely above the visible area
		if currentY+lineSpacing < float64(y) {
			currentY += lineSpacing
			continue
		}

		// Stop drawing lines that are completely below the visible area
		if currentY > float64(y+height) {
			break
		}

		// Draw the line
		op := &text.DrawOptions{}
		op.GeoM.Translate(float64(x), currentY)
		op.ColorScale.ScaleWithColor(textColor)
		op.LineSpacing = lineSpacing
		text.Draw(subImage, line, textFace, op)

		currentY += lineSpacing
	}
}

// clampScrollOffset ensures scroll offset stays within valid bounds
func (s *Scene) clampScrollOffset() {
	maxScroll := s.textHeight - (screenHeight - contentStartY - padding)
	if maxScroll < 0 {
		maxScroll = 0
	}

	if s.scrollOffset < 0 {
		s.scrollOffset = 0
	} else if s.scrollOffset > maxScroll {
		s.scrollOffset = maxScroll
	}
}

// Layout returns the logical screen size
func (s *Scene) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screen.Width, screen.Height
}

// NewScene creates a new credits scene
func NewScene(scenes *scene.Manager) *Scene {
	// Combine all license texts
	var buf bytes.Buffer

	// Fonts
	buf.Write(credits.Fonts())
	buf.WriteString("\n\n---\n\n")

	// Go Gopher
	buf.Write(credits.GoGopher())
	buf.WriteString("\n\n---\n\n")

	// Go Programming Language
	buf.Write(credits.GoProgrammingLanguage())
	buf.WriteString("\n\n---\n\n")

	// Third-party licenses
	buf.Write(credits.Licenses())

	return &Scene{
		scenes: scenes,
		closeButton: &libui.RoundedRectButton{
			X:      screenWidth - 100 - padding,
			Y:      padding,
			W:      100,
			H:      buttonHeight,
			Radius: 10,
			Text: libui.Text{
				Text:    "return",
				Face:    ui.TextFace16,
				Color:   color.RGBA{R: 255, G: 255, B: 255, A: 255},
				BGColor: color.RGBA{R: 0, G: 0, B: 0, A: 0},
			},
			DefaultColor: color.RGBA{R: 100, G: 100, B: 200, A: 255},
			PressedColor: color.RGBA{R: 150, G: 150, B: 250, A: 255},
		},
		scrollOffset: 0,
		contentText:  buf.String(),
	}
}
