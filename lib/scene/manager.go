package scene

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/raiich/kazura/state"
	"github.com/raiich/kazura/task/eventloop"
)

// Manager manages game scenes and provides scene transition functionality.
// It implements the ebiten.Game interface and delegates to the current scene.
type Manager struct {
	dispatcher *eventloop.Dispatcher
	now        time.Time
	delta      time.Duration
	current    Scene
	stack      []Scene
}

// Update updates the current scene and advances the event dispatcher.
// Implements ebiten.Game interface.
func (m *Manager) Update() error {
	// Advance dispatcher time to current time to execute scheduled tasks
	m.now = m.now.Add(m.delta)
	if err := m.dispatcher.FastForward(m.now); err != nil {
		return err
	}
	return m.current.Update()
}

// Draw draws the current scene to the screen.
// Implements ebiten.Game interface.
func (m *Manager) Draw(screen *ebiten.Image) {
	m.current.Draw(screen)
}

// Layout delegates to the current scene's Layout method.
// Implements ebiten.Game interface.
func (m *Manager) Layout(outsideWidth, outsideHeight int) (int, int) {
	return m.current.Layout(outsideWidth, outsideHeight)
}

// SetScene sets the current scene, replacing any existing scene and clearing the scene stack.
func (m *Manager) SetScene(scene Scene) error {
	m.current = scene
	m.stack = nil
	return nil
}

// PushScene pushes a new scene onto the stack, making it the current scene.
// The previous scene can be restored by calling PopScene.
func (m *Manager) PushScene(scene Scene) error {
	m.stack = append(m.stack, m.current)
	m.current = scene
	return nil
}

// PopScene pops the current scene from the stack and restores the previous scene.
// Returns an error if the stack is empty.
func (m *Manager) PopScene() error {
	if len(m.stack) == 0 {
		return fmt.Errorf("no scene to pop")
	}
	head, tail := m.stack[len(m.stack)-1], m.stack[:len(m.stack)-1]
	m.current = head
	m.stack = tail
	return nil
}

// Dispatcher returns the event dispatcher for scheduling tasks.
func (m *Manager) Dispatcher() state.Dispatcher {
	return m.dispatcher
}

// NewSceneManager creates a new scene manager with an event dispatcher.
func NewSceneManager() *Manager {
	now := time.Now()
	delta := float64(1*time.Second) / float64(ebiten.TPS())
	return &Manager{
		dispatcher: eventloop.NewDispatcher(now),
		now:        now,
		delta:      time.Duration(delta),
	}
}
