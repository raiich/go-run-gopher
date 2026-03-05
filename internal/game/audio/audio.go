// Package audio provides musical scale playback for the game
package audio

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

const (
	sampleRate = 48000 // Sample rate for audio playback
	noteLength = 0.8   // Duration of each note in seconds (long for flowing sound)
)

// Musical scale frequencies (C4 to C5 - do-re-mi-fa-so-la-si-do)
var scaleFrequencies = [8]float64{
	261.63, // C4 (Do)
	293.66, // D4 (Re)
	329.63, // E4 (Mi)
	349.23, // F4 (Fa)
	392.00, // G4 (So)
	440.00, // A4 (La)
	493.88, // B4 (Si)
	523.25, // C5 (Do)
}

// ScalePlayer manages musical scale playback during gameplay
type ScalePlayer struct {
	context *audio.Context
	players [8]*audio.Player // Pre-generated players for each note
}

// NewScalePlayer creates a new scale player
func NewScalePlayer() *ScalePlayer {
	ctx := audio.NewContext(sampleRate)
	sp := &ScalePlayer{
		context: ctx,
	}

	// Pre-generate audio players for all 8 notes
	for i := 0; i < 8; i++ {
		sp.players[i] = sp.createPlayer(scaleFrequencies[i])
	}

	return sp
}

// createPlayer creates an audio player for a specific frequency
// Generates a piano-like tone with harmonics and ADSR envelope
func (sp *ScalePlayer) createPlayer(frequency float64) *audio.Player {
	numSamples := int(float64(sampleRate) * noteLength)
	pcm := make([]byte, numSamples*4) // 4 bytes per sample (2 channels × 2 bytes)

	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(sampleRate)

		// Piano-like waveform: fundamental + harmonics (higher harmonics decay faster)
		sample := math.Sin(2.0*math.Pi*frequency*t) * 1.0
		sample += math.Sin(2.0*math.Pi*2*frequency*t) * 0.5 * math.Exp(-t*3)
		sample += math.Sin(2.0*math.Pi*3*frequency*t) * 0.25 * math.Exp(-t*5)
		sample += math.Sin(2.0*math.Pi*4*frequency*t) * 0.125 * math.Exp(-t*7)

		// Piano-like envelope: quick attack, exponential decay, sustain floor
		attack := math.Min(t/0.005, 1.0)          // 5ms attack
		decay := math.Exp(-t * 2.5)               // exponential decay
		envelope := attack * math.Max(decay, 0.2) // sustain floor

		// Fade out at the end to prevent clicking
		fadeOutSamples := 2000
		if i > numSamples-fadeOutSamples {
			envelope *= float64(numSamples-i) / float64(fadeOutSamples)
		}

		sample *= envelope

		// Convert to 16-bit PCM (volume reduced to account for harmonics)
		pcmValue := int16(sample * 0x7fff * 0.2)

		// Little-endian format: left channel
		pcm[i*4] = byte(pcmValue)
		pcm[i*4+1] = byte(pcmValue >> 8)
		// Right channel (same as left for mono)
		pcm[i*4+2] = byte(pcmValue)
		pcm[i*4+3] = byte(pcmValue >> 8)
	}

	player := sp.context.NewPlayerFromBytes(pcm)
	return player
}

// PlayNote plays the note at the specified index (0-7)
func (sp *ScalePlayer) PlayNote(index int) {
	if index < 0 || index >= 8 {
		return
	}

	player := sp.players[index]
	// Rewind to beginning (in case it was played before)
	player.Rewind()
	// Start playback
	player.Play()
}
