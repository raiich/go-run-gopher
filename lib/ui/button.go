package ui

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Button represents a clickable UI button with hit detection and rendering
type Button interface {
	Contains(x, y float64) bool
	SetPressed(bool)
	IsPressed() bool
	Draw(screen *ebiten.Image, textFace *text.GoTextFace)
}

// RectButton is a rectangular button with solid color background
type RectButton struct {
	X, Y, W, H   float64
	Text         Text
	DefaultColor color.Color
	PressedColor color.Color
	pressed      bool
}

func (b *RectButton) Contains(x, y float64) bool {
	return x >= b.X && x <= b.X+b.W && y >= b.Y && y <= b.Y+b.H
}

func (b *RectButton) SetPressed(pressed bool) {
	b.pressed = pressed
}

func (b *RectButton) IsPressed() bool {
	return b.pressed
}

func (b *RectButton) Draw(screen *ebiten.Image, textFace *text.GoTextFace) {
	btnColor := b.DefaultColor
	if b.pressed {
		btnColor = b.PressedColor
	}
	vector.DrawFilledRect(screen, float32(b.X), float32(b.Y), float32(b.W), float32(b.H), btnColor, false)

	b.Text.Image().Draw(screen, image.Point{X: int(b.X + 15), Y: int(b.Y + 15)})
}

// RoundedRectButton is a rectangular button with rounded corners
type RoundedRectButton struct {
	X, Y, W, H   float64
	Radius       float64 // corner radius
	Text         Text
	DefaultColor color.Color
	PressedColor color.Color
	pressed      bool
}

func (b *RoundedRectButton) Contains(x, y float64) bool {
	return x >= b.X && x <= b.X+b.W && y >= b.Y && y <= b.Y+b.H
}

func (b *RoundedRectButton) SetPressed(pressed bool) {
	b.pressed = pressed
}

func (b *RoundedRectButton) IsPressed() bool {
	return b.pressed
}

func (b *RoundedRectButton) Draw(screen *ebiten.Image, textFace *text.GoTextFace) {
	btnColor := b.DefaultColor
	if b.pressed {
		btnColor = b.PressedColor
	}

	// Draw rounded rectangle using filled rects and circles
	x, y, w, h, r := float32(b.X), float32(b.Y), float32(b.W), float32(b.H), float32(b.Radius)

	// Draw center rectangle
	vector.DrawFilledRect(screen, x+r, y+r, w-2*r, h-2*r, btnColor, false)

	// Draw edge rectangles
	vector.DrawFilledRect(screen, x+r, y, w-2*r, r, btnColor, false)     // top
	vector.DrawFilledRect(screen, x+r, y+h-r, w-2*r, r, btnColor, false) // bottom
	vector.DrawFilledRect(screen, x, y+r, r, h-2*r, btnColor, false)     // left
	vector.DrawFilledRect(screen, x+w-r, y+r, r, h-2*r, btnColor, false) // right

	// Draw corner circles
	vector.DrawFilledCircle(screen, x+r, y+r, r, btnColor, false)     // top-left
	vector.DrawFilledCircle(screen, x+w-r, y+r, r, btnColor, false)   // top-right
	vector.DrawFilledCircle(screen, x+r, y+h-r, r, btnColor, false)   // bottom-left
	vector.DrawFilledCircle(screen, x+w-r, y+h-r, r, btnColor, false) // bottom-right

	tw, th := text.Measure(b.Text.Text, b.Text.Face, b.Text.Face.Size*1.3)
	b.Text.Image().Draw(screen, image.Point{X: int(b.X + b.W/2 - tw/2), Y: int(b.Y + b.H/2 - th/2)})
}

// CircleButton is a circular button with solid color background
type CircleButton struct {
	X, Y, R      float64
	Text         Text
	DefaultColor color.Color
	PressedColor color.Color
	pressed      bool
}

func (b *CircleButton) Contains(x, y float64) bool {
	dx := x - b.X
	dy := y - b.Y
	return dx*dx+dy*dy <= b.R*b.R
}

func (b *CircleButton) SetPressed(pressed bool) {
	b.pressed = pressed
}

func (b *CircleButton) IsPressed() bool {
	return b.pressed
}

func (b *CircleButton) SetText(text string) {
	b.Text.Text = text
}

func (b *CircleButton) Draw(screen *ebiten.Image, textFace *text.GoTextFace) {
	btnColor := b.DefaultColor
	if b.pressed {
		btnColor = b.PressedColor
	}
	cx, cy, r := float32(b.X), float32(b.Y), float32(b.R)
	vector.DrawFilledCircle(screen, cx, cy, r, btnColor, false)

	// Center text using measurement (shift up slightly for visual balance with lowercase)
	w, h := text.Measure(b.Text.Text, b.Text.Face, b.Text.Face.Size*1.3)
	b.Text.Image().Draw(screen, image.Point{X: int(b.X - w/2), Y: int(b.Y - h/2 - 4)})
}
