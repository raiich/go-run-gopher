// Package effect provides temporary text overlay effects for gameplay events.
package effect

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	libui "github.com/raiich/go-run-gopher/lib/ui"
)

// Kind identifies the type of preset effect to display.
type Kind int

const (
	KindPanic   Kind = iota // "panic" — red blinking, ~30 frames
	KindRecover             // "recover" — green fade-in, ~30 frames
)

type effectDef struct {
	text     string
	color    color.RGBA
	duration int  // total frames
	blink    bool // toggle alpha every few frames instead of fade
}

var defs = [...]effectDef{
	KindPanic:   {text: "panic", color: color.RGBA{R: 255, G: 60, B: 60, A: 255}, duration: 30, blink: true},
	KindRecover: {text: "recover", color: color.RGBA{R: 100, G: 255, B: 100, A: 255}, duration: 30},
}

// Overlay manages a single fire-and-forget text effect displayed on screen.
type Overlay struct {
	face         *text.GoTextFace
	screenWidth  int
	screenHeight int

	active bool
	def    effectDef
	frame  int
}

// New creates an Overlay for the given screen dimensions.
func New(face *text.GoTextFace, screenWidth, screenHeight int) *Overlay {
	return &Overlay{
		face:         face,
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
	}
}

// Show starts displaying a preset effect kind. Replaces any active effect.
func (o *Overlay) Show(kind Kind) {
	o.active = true
	o.def = defs[kind]
	o.frame = 0
}

// Update advances the effect by one frame.
func (o *Overlay) Update() {
	if !o.active {
		return
	}
	o.frame++
	if o.frame >= o.def.duration {
		o.active = false
	}
}

// Draw renders the current effect onto the screen.
func (o *Overlay) Draw(screen *ebiten.Image) {
	if !o.active {
		return
	}

	// Calculate alpha
	var alpha uint8
	if o.def.blink {
		// Toggle every 4 frames
		if (o.frame/4)%2 == 0 {
			alpha = 255
		} else {
			alpha = 60
		}
	} else {
		// Fade out over duration
		progress := float64(o.frame) / float64(o.def.duration)
		alpha = uint8((1.0 - progress) * 255)
	}

	c := o.def.color
	c.A = alpha

	t := libui.Text{
		Text:  o.def.text,
		Face:  o.face,
		Color: c,
	}
	img := t.Image()
	w := img.Rectangle().Dx()
	h := img.Rectangle().Dy()
	x := o.screenWidth/2 - w/2
	y := o.screenHeight/3 - h/2
	img.Draw(screen, image.Point{X: x, Y: y})
}
