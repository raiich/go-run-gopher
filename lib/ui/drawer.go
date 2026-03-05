package ui

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Drawer interface {
	Rectangle() image.Rectangle
	Draw(screen *ebiten.Image, translate image.Point)
}

type DrawerTranslate struct {
	Drawer
	Translate image.Point
}

func (d *DrawerTranslate) Draw(screen *ebiten.Image, translate image.Point) {
	d.Drawer.Draw(screen, d.Translate.Add(translate))
}

type CompoundDrawer []Drawer

func (c CompoundDrawer) Rectangle() image.Rectangle {
	if len(c) == 0 {
		return image.Rectangle{}
	}
	first := c[0].Rectangle()
	minX, minY, maxX, maxY := first.Min.X, first.Min.Y, first.Max.X, first.Max.Y
	for _, d := range c[1:] {
		r := d.Rectangle()
		minX = min(r.Min.X, minX)
		minY = min(r.Min.Y, minY)
		maxX = max(r.Max.X, maxX)
		maxY = max(r.Max.Y, maxY)
	}
	return image.Rect(minX, minY, maxX, maxY)
}

func (c CompoundDrawer) Draw(screen *ebiten.Image, translate image.Point) {
	for _, e := range c {
		e.Draw(screen, translate)
	}
}

type Text struct {
	Text    string
	Face    *text.GoTextFace
	Color   color.Color
	BGColor color.Color
	Padding Gap
	Spacing float64
	Bounds  image.Rectangle
}

func (t Text) Image() *TextImage {
	spacing := t.Spacing
	if spacing == 0 {
		spacing = t.Face.Size * 1.3
	}
	ret := &TextImage{
		text: &t,
	}
	ret.fit()
	return ret
}

type TextImage struct {
	text      *Text
	width     int
	height    int
	translate image.Point
}

func (t *TextImage) Fit(line string) {
	t.text.Text = line
	t.fit()
}

func (t *TextImage) fit() {
	p := t.text.Padding
	t.translate.X, t.translate.Y = p.Left, p.Top
	w, h := text.Measure(t.text.Text, t.text.Face, t.text.Spacing)

	innerWidth, innerHeight := int(w)+p.Left+p.Right, int(h)+p.Top+p.Bottom
	outerWidth, outerHeight := t.text.Bounds.Dx(), t.text.Bounds.Dy()
	if outerWidth > 0 {
		if innerWidth < outerWidth {
			t.translate.X = (outerWidth - innerWidth) / 2
		}
	} else {
		outerWidth = innerWidth
	}
	if outerHeight > 0 {
		if innerHeight < outerHeight {
			t.translate.Y = (outerHeight - innerHeight) / 2
		}
	} else {
		outerHeight = innerHeight
	}

	t.width, t.height = outerWidth, outerHeight
}

func (t *TextImage) Rectangle() image.Rectangle {
	return image.Rect(0, 0, t.width, t.height)
}

func (t *TextImage) Draw(screen *ebiten.Image, translate image.Point) {
	if t.text.BGColor != nil {
		x, y := translate.X, translate.Y
		w, h := float32(t.width), float32(t.height)
		vector.DrawFilledRect(screen, float32(x), float32(y), w, h, t.text.BGColor, false)
	}

	op := &text.DrawOptions{}
	tt := translate.Add(t.translate)
	op.GeoM.Translate(float64(tt.X), float64(tt.Y))

	op.ColorScale.ScaleWithColor(t.text.Color)
	op.LineSpacing = t.text.Spacing
	text.Draw(screen, t.text.Text, t.text.Face, op)
}

type Rect struct {
	Color  color.Color
	Width  int
	Height int
}

func (b Rect) Image() *RectImage {
	img := ebiten.NewImage(b.Width, b.Height)
	img.Fill(b.Color)
	return (*RectImage)(img)
}

type RectImage ebiten.Image

func (i *RectImage) Rectangle() image.Rectangle {
	return (*ebiten.Image)(i).Bounds()
}

func (i *RectImage) Draw(screen *ebiten.Image, translate image.Point) {
	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(translate.X), float64(translate.Y))
	screen.DrawImage((*ebiten.Image)(i), &op)
}

type Line struct {
	Color color.Color
	Width int
	Start image.Point
	End   image.Point
}

func (b Line) Image() *LineImage {
	return &LineImage{
		color: b.Color,
		width: float32(b.Width),
		x0:    float32(b.Start.X),
		y0:    float32(b.Start.Y),
		x1:    float32(b.End.X),
		y1:    float32(b.End.Y),
	}
}

type LineImage struct {
	x0    float32
	y0    float32
	x1    float32
	y1    float32
	color color.Color
	width float32
}

func (i *LineImage) Rectangle() image.Rectangle {
	minX := int(min(i.x0, i.x1) - i.width/2)
	minY := int(min(i.y0, i.y1) - i.width/2)
	maxX := int(max(i.x0, i.x1) + i.width/2)
	maxY := int(max(i.y0, i.y1) + i.width/2)
	return image.Rect(minX, minY, maxX, maxY)
}

func (i *LineImage) Draw(screen *ebiten.Image, translate image.Point) {
	x0, y0 := i.x0+float32(translate.X), i.y0+float32(translate.Y)
	x1, y1 := i.x1+float32(translate.X), i.y1+float32(translate.Y)
	vector.StrokeLine(screen, x0, y0, x1, y1, i.width, i.color, false)
}

/////////////////////////

type Circle struct {
	Color  color.Color
	Width  int
	Height int
}

func (b Circle) Image() *CircleImage {
	length := max(b.Width, b.Height)
	screen := ebiten.NewImage(length, length)
	r := length / 2
	c := b.Color
	vector.DrawFilledCircle(screen, float32(r), float32(r), float32(r), c, false)
	return &CircleImage{
		circle: screen,
		length: float64(length),
		scaleX: float64(b.Width) / float64(length),
		scaleY: float64(b.Height) / float64(length),
		rect:   image.Rect(0, 0, b.Width, b.Height),
	}
}

type CircleImage struct {
	circle *ebiten.Image
	length float64
	scaleX float64
	scaleY float64
	rect   image.Rectangle
}

func (i *CircleImage) Rectangle() image.Rectangle {
	return i.rect
}

func (i *CircleImage) Draw(screen *ebiten.Image, translate image.Point) {
	op := ebiten.DrawImageOptions{}
	scaleX, scaleY := i.scaleX, i.scaleY
	op.GeoM.Scale(scaleX, scaleY)
	op.GeoM.Translate(float64(translate.X), float64(translate.Y))
	screen.DrawImage(i.circle, &op)
}
