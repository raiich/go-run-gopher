// Package ui provides UI utilities for text rendering and formatting.
package ui

import (
	"bytes"
	_ "embed"
	"strings"

	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/raiich/kazura/must"
)

const (
	// DefaultTextOffset is the default text offset in pixels
	DefaultTextOffset = 64.0
	// DefaultLineWrapLen is the default maximum line length for text wrapping
	DefaultLineWrapLen = 100
	// DefaultLineSpacing is the default line spacing multiplier
	DefaultLineSpacing = 1.3
)

// DefaultGoTextFaceSource is the default font face source using MPlus1p Regular font
var DefaultGoTextFaceSource = must.Must(text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf)))

// TextFormatter handles text wrapping and formatting
type TextFormatter struct {
	MaxLineLength int
}

// NewTextFormatter creates a new text formatter with default settings
func NewTextFormatter() *TextFormatter {
	return &TextFormatter{
		MaxLineLength: DefaultLineWrapLen,
	}
}

// WrapLines wraps multiple lines of text according to the maximum line length
func (f *TextFormatter) WrapLines(lines []string) []string {
	var result []string
	for _, line := range lines {
		result = append(result, f.wrapLine(line)...)
	}
	return result
}

func (f *TextFormatter) wrapLine(line string) []string {
	if len(line) <= f.MaxLineLength {
		return []string{line}
	}

	var result []string
	remaining := line

	for len(remaining) > f.MaxLineLength {
		tokens := strings.Split(remaining, " ")
		var currentLine string

		for i, token := range tokens {
			nextLength := len(currentLine)
			if currentLine != "" {
				nextLength++
			}
			nextLength += len(token)

			if nextLength > f.MaxLineLength && currentLine != "" {
				result = append(result, currentLine)
				remaining = strings.Join(tokens[i:], " ")
				break
			}

			if currentLine == "" {
				currentLine = token
			} else {
				currentLine += " " + token
			}

			if i == len(tokens)-1 {
				result = append(result, currentLine)
				return result
			}
		}
	}

	if remaining != "" {
		result = append(result, remaining)
	}

	return result
}
