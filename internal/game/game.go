package game

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/raiich/go-run-gopher/internal/game/audio"
	"github.com/raiich/go-run-gopher/internal/game/effect"
	"github.com/raiich/go-run-gopher/internal/game/gopher"
	"github.com/raiich/go-run-gopher/internal/game/pregame"
	"github.com/raiich/go-run-gopher/internal/game/proverbs"
	ui2 "github.com/raiich/go-run-gopher/internal/game/ui"
	"github.com/raiich/go-run-gopher/lib/scene"
	libui "github.com/raiich/go-run-gopher/lib/ui"
	"github.com/raiich/kazura/must"
	"github.com/raiich/kazura/state"
	"github.com/raiich/kazura/task"
)

// stateGraph defines the screen state transitions
var stateGraph = check(state.NewGraph[sceneState](
	titleState{},
	on[startPregameEvent](titleState{}, pregameState{}),
	on[gameStartEvent](pregameState{}, runningState{}),
	on[gameEndEvent](runningState{}, resultState{}),
	on[resetEvent](resultState{}, titleState{}),
))

type sceneEvent interface {
	state.Event
}

type sceneState interface {
	Entry(machine *state.EntryMachine[*sceneData], event state.Event) state.Command
	handleInput(sc *Scene, button ui2.ButtonType, justPressed bool)
	update(sc *Scene) error
	draw(sc *Scene, screen *ebiten.Image)
}

// Events

type startPregameEvent struct{}
type gameStartEvent struct{}
type gameEndEvent struct{}
type resetEvent struct{}

// titleState represents the title screen state
type titleState struct{}

func (s titleState) Entry(machine *state.EntryMachine[*sceneData], event state.Event) state.Command {
	data := machine.Value()
	data.gameData.ResetForNewGame()
	data.gopher.SpeedMax = calculateSpeedForStage(data.gameData.stage)
	data.gopher.X = data.gopher.InitialX
	data.gopher.Y = data.gopher.InitialY
	data.gopher.Direction = 1
	data.ui.ActionButton.SetText("go")
	return nil
}

func (s titleState) handleInput(sc *Scene, button ui2.ButtonType, justPressed bool) {
	// Start button - starts the game
	if button == ui2.ButtonTypeAction && justPressed {
		must.NoError(sc.machine.Trigger(startPregameEvent{}))
	}

	// About button - shows credits scene
	if button == ui2.ButtonTypeAbout && justPressed {
		sc.showCredits()
	}
}

func (s titleState) update(sc *Scene) error {
	// No update needed for title state
	return nil
}

func (s titleState) draw(sc *Scene, screen *ebiten.Image) {
	// draw UI
	sc.ui.Draw(screen, ui2.ButtonTypeAction)
	sc.ui.Draw(screen, ui2.ButtonTypeAbout)
}

// pregameState represents the pregame state before game starts
type pregameState struct{}

func (s pregameState) Entry(machine *state.EntryMachine[*sceneData], event state.Event) state.Command {
	data := machine.Value()

	data.ui.ActionButton.SetText("go")

	// Entry cannot trigger a transition, and the callback may run inside it,
	// so the event is carried out of the transition by the dispatcher.
	must.NoError(data.overlay.Start(0, func() {
		must.NoError(machine.AfterFunc(data.dispatcher, 0, func(machine *state.AfterFuncMachine[*sceneData]) {
			must.NoError(machine.Trigger(gameStartEvent{}))
		}))
	}))
	return nil
}

func (s pregameState) handleInput(sc *Scene, button ui2.ButtonType, justPressed bool) {
	// No input handling during pregame
}

func (s pregameState) update(sc *Scene) error {
	// No update needed
	return nil
}

func (s pregameState) draw(sc *Scene, screen *ebiten.Image) {
	// draw UI
	sc.ui.Draw(screen, ui2.ButtonTypeAction)
}

// runningState represents the main game running state
type runningState struct{}

func (s runningState) Entry(machine *state.EntryMachine[*sceneData], event state.Event) state.Command {
	// nothing to do
	return nil
}

func (s runningState) handleInput(sc *Scene, button ui2.ButtonType, justPressed bool) {
	// Delegate input handling to gopher (which will ignore input during bounce)
	sc.gopher.HandleInput(button, justPressed)
}

func (s runningState) update(sc *Scene) error {
	sc.gameData.totalFrames++

	// Increment lap timer (tracks time since last lap completion)
	sc.gameData.lapTimer++

	// Update gopher physics and animation
	g := sc.gopher
	g.Update()

	// Play musical scale synced to lap timer
	// Each note plays at (beepIndex+1) * (maxLapTime/8) frames into the lap
	beepInterval := sc.gameData.maxLapTime / 8
	if sc.gameData.beepIndex < 8 && sc.gameData.lapTimer >= sc.gameData.beepIndex*beepInterval {
		noteIndex := sc.gameData.beepIndex
		if sc.gameData.scaleDirection == -1 {
			noteIndex = 7 - sc.gameData.beepIndex
		}
		sc.scalePlayer.PlayNote(noteIndex)
		sc.gameData.beepIndex++
		sc.gameData.scaleCount = 8 - sc.gameData.beepIndex // 7, 6, 5, 4, 3, 2, 1, 0
	}

	// Screen edges: absolute boundaries that trigger bounce if hit while moving
	atLeftEdge := g.X <= 0
	atRightEdge := g.X >= g.FieldBounds.Width-g.Size

	// Boundary lines: goal zones within the boundary margin
	crossedLeftLine := g.X <= boundaryMargin
	crossedRightLine := g.X >= g.FieldBounds.Width-g.Size-boundaryMargin

	// Bounce penalty: hitting edge while moving triggers bounce animation
	if (atLeftEdge || atRightEdge) && g.GetSpeed() > 0 {
		sc.effect.Show(effect.KindPanic)
		return sc.gopher.StartBounce()
	}

	// Keep gopher stopped at edge if already stopped (no bounce)
	if (atLeftEdge || atRightEdge) && g.GetSpeed() == 0 {
		g.AdjustSpeed(0) // Reset halving timer
	}

	// Track last reached goal based on current position in boundary zones
	if crossedLeftLine {
		sc.gameData.lastReachedGoal = -1
	} else if crossedRightLine {
		sc.gameData.lastReachedGoal = 1
	}

	// After all 8 notes have played, evaluate success or miss
	// Check if the LAST reached goal matches the expected goal side
	if sc.gameData.beepIndex >= 8 && sc.gameData.lapTimer > sc.gameData.maxLapTime {
		lapSec := float64(sc.gameData.lapTimer) / 60.0
		testName := fmt.Sprintf("TestLap%d", sc.gameData.laps+1)
		if sc.gameData.lastReachedGoal == sc.gameData.goalSide {
			// Success: last reached goal matches expected side
			sc.gameData.testLog = append(sc.gameData.testLog, testResult{passed: true, name: testName, sec: lapSec})
			sc.gameData.misses = 0
			sc.gameData.laps++
			sc.gameData.goalSide = -sc.gameData.goalSide
			sc.gameData.lapTimer = 0
			sc.gameData.beepIndex = 0
			sc.gameData.scaleCount = 8
			sc.gameData.scaleDirection = sc.gameData.goalSide // right=ascending, left=descending
			sc.gameData.maxLapTime = calculateMaxLapTime(sc.gameData.laps)

			// Stage progression: every 7 laps increases difficulty
			newStage := calculateStage(sc.gameData.laps)
			if newStage > sc.gameData.stage {
				sc.gameData.stage = newStage
				sc.gopher.SpeedMax = calculateSpeedForStage(sc.gameData.stage)
			}
		} else {
			// Miss: continue to next lap (shuttle run rules — music doesn't stop)
			sc.gameData.testLog = append(sc.gameData.testLog, testResult{passed: false, name: testName, sec: lapSec})
			sc.gameData.misses++
			if sc.gameData.misses >= 2 {
				return sc.machine.Trigger(gameEndEvent{})
			}
			// Lap counts even on miss (shuttle run rules)
			sc.gameData.laps++
			sc.gameData.goalSide = -sc.gameData.goalSide
			sc.gameData.lapTimer = 0
			sc.gameData.beepIndex = 0
			sc.gameData.scaleCount = 8
			sc.gameData.scaleDirection = sc.gameData.goalSide
			sc.gameData.maxLapTime = calculateMaxLapTime(sc.gameData.laps)

			// Stage progression: every 7 laps increases difficulty
			newStage := calculateStage(sc.gameData.laps)
			if newStage > sc.gameData.stage {
				sc.gameData.stage = newStage
				sc.gopher.SpeedMax = calculateSpeedForStage(sc.gameData.stage)
			}
		}
	}

	// Auto-reverse: if stopped at wrong boundary facing wrong way, flip direction
	if (crossedLeftLine || crossedRightLine) && g.GetSpeed() == 0 {
		shouldReverse := (crossedLeftLine && g.Direction == -1) || (crossedRightLine && g.Direction == 1)
		if shouldReverse {
			sc.gopher.Direction = -sc.gopher.Direction
		}
	}

	return nil
}

func (s runningState) draw(sc *Scene, screen *ebiten.Image) {
	sc.ui.Draw(screen, ui2.ButtonTypeAction)
	drawGameUI(sc, screen)
}

// resultState represents the game over/result screen state
type resultState struct{}

func (s resultState) Entry(machine *state.EntryMachine[*sceneData], event state.Event) state.Command {
	data := machine.Value()
	data.gopher.AdjustSpeed(0)
	data.gameData.proverb = proverbs.Random()
	data.gameData.resultReady = false

	// Enable input after 3 seconds cooldown
	must.NoError(machine.AfterFunc(data.dispatcher, 3*time.Second, func(machine *state.AfterFuncMachine[*sceneData]) {
		d := machine.Value()
		d.gameData.resultReady = true
		d.ui.ActionButton.SetText("return")
	}))
	return nil
}

func (s resultState) handleInput(sc *Scene, button ui2.ButtonType, justPressed bool) {
	if sc.gameData.resultReady && button == ui2.ButtonTypeAction && justPressed {
		must.NoError(sc.machine.Trigger(resetEvent{}))
	}
}

func (s resultState) update(sc *Scene) error {
	// No update needed for result state (the timer in Entry enables input)
	return nil
}

func (s resultState) draw(sc *Scene, screen *ebiten.Image) {
	// Semi-transparent dark overlay for terminal-style result display
	vector.DrawFilledRect(screen, 0, 0, float32(ScreenWidth), float32(ScreenHeight), color.RGBA{A: 180}, false)

	if sc.gameData.resultReady {
		sc.ui.Draw(screen, ui2.ButtonTypeAction)
	}

	cfg := sc.uiConfig
	totalSec := float64(sc.gameData.totalFrames) / 60.0
	passColor := color.RGBA{R: 100, G: 255, B: 100, A: 255}
	failColor := color.RGBA{R: 255, G: 80, B: 80, A: 255}
	proverbColor := color.RGBA{R: 180, G: 180, B: 180, A: 200}

	// Build terminal-style output lines
	type line struct {
		s    string
		c    color.RGBA
		skip bool // true = empty line (just advance Y)
	}
	var lines []line

	// Test log (last 3 results)
	log := sc.gameData.testLog
	start := 0
	if len(log) > 3 {
		start = len(log) - 3
	}
	for _, r := range log[start:] {
		if r.passed {
			lines = append(lines, line{
				s: fmt.Sprintf("--- PASS: %s (%.2fs)", r.name, r.sec),
				c: passColor,
			})
		} else {
			lines = append(lines, line{
				s: fmt.Sprintf("--- FAIL: %s (%.2fs)", r.name, r.sec),
				c: failColor,
			})
		}
	}

	// Benchmark + summary (subtract 1 for the failed lap that ended the game)
	passedLaps := sc.gameData.laps - 1
	if passedLaps > 0 {
		perLap := totalSec / float64(passedLaps)
		lines = append(lines, line{
			s: fmt.Sprintf("BenchmarkGopher    %d laps    %.2fs/lap", passedLaps, perLap),
			c: cfg.Result.Color,
		})
	} else {
		lines = append(lines, line{s: "BenchmarkGopher    0 laps", c: cfg.Result.Color})
	}
	lines = append(lines, line{
		s: fmt.Sprintf("FAIL"),
		c: failColor,
	})

	// Measure max width from non-proverb lines to center the block
	maxW := 0
	for _, l := range lines {
		if l.skip {
			continue
		}
		t := libui.Text{Text: l.s, Face: cfg.Result.Face, Color: l.c}
		if w := t.Image().Rectangle().Dx(); w > maxW {
			maxW = w
		}
	}

	// Proverb: wrap to fit within the block width
	lines = append(lines, line{skip: true})
	proverbPrefix := "// "
	proverbText := sc.gameData.proverb
	words := strings.Fields(proverbText)
	current := proverbPrefix
	for _, w := range words {
		test := current + w
		t := libui.Text{Text: test, Face: cfg.Result.Face, Color: proverbColor}
		if t.Image().Rectangle().Dx() > maxW && current != proverbPrefix {
			lines = append(lines, line{s: strings.TrimRight(current, " "), c: proverbColor})
			current = "// " + w + " "
		} else {
			current = test + " "
		}
	}
	if strings.TrimSpace(current) != "//" {
		lines = append(lines, line{s: strings.TrimRight(current, " "), c: proverbColor})
	}

	lineHeight := 22
	blockX := ScreenWidth/2 - maxW/2
	startY := ScreenHeight/2 - lineHeight*len(lines)/2

	for i, l := range lines {
		if l.skip {
			continue
		}
		t := libui.Text{Text: l.s, Face: cfg.Result.Face, Color: l.c}
		t.Image().Draw(screen, image.Point{X: blockX, Y: startY + i*lineHeight})
	}
}

// game holds the scene state data
type sceneData struct {
	scenes      *scene.Manager
	dispatcher  task.Dispatcher
	ui          *ui2.Controller
	gopher      *gopher.Gopher
	effect      *effect.Overlay
	gameData    gameData
	overlay     *pregame.Overlay
	scalePlayer *audio.ScalePlayer
	uiConfig    UIConfig
}

// testResult stores the outcome of a single lap evaluation.
type testResult struct {
	passed bool
	name   string
	sec    float64
}

// gameData holds game-specific data managed by the Scene state machine.
// This is separate from Gopher's animation data to maintain clear separation of concerns.
type gameData struct {
	// Number of completed laps
	laps int
	// Current stage number (starts at 1)
	stage int
	// Number of consecutive misses
	misses int
	// Timer for current lap (frames)
	lapTimer int
	// Maximum frames allowed for current lap
	maxLapTime int
	// Which side is the goal for next lap: -1 = left, 1 = right
	goalSide int
	// Last goal side the gopher reached: 0 = none, -1 = left, 1 = right
	// Continuously updated based on gopher's position in boundary zones
	lastReachedGoal int
	// Current beep index in the scale (0-7)
	beepIndex int
	// Scale direction: 1 = ascending (do-re-mi...), -1 = descending (do-si-la...)
	scaleDirection int
	// Remaining scale count displayed to player (8..0), driven by note events
	scaleCount int
	// Total frames elapsed during running state (for bench-style result)
	totalFrames int
	// Go proverb displayed on result screen
	proverb string
	// Accumulated test results for HUD log display
	testLog []testResult
	// Whether the result screen accepts input (after cooldown)
	resultReady bool
}

// ResetForNewGame resets game data for a new game
func (d *gameData) ResetForNewGame() {
	*d = gameData{
		stage:          1,
		maxLapTime:     calculateMaxLapTime(0),
		goalSide:       1, // First goal is right side
		beepIndex:      0,
		scaleDirection: 1, // Ascending (right goal = do-re-mi...)
		scaleCount:     8,
	}
}

// calculateStage calculates the stage based on lap count
func calculateStage(lapCount int) int {
	return (lapCount / lapsPerStage) + 1
}

// calculateMaxLapTime calculates the maximum time allowed for a lap based on lap count
// Gets faster each lap
// Base time is 8 seconds to match shuttle run music (do-re-mi-fa-so-la-si-do = 8 counts)
func calculateMaxLapTime(lapCount int) int {
	baseTime := 480            // 8 seconds at 60fps
	reduction := lapCount * 10 // Reduce by 10 frames (1/6 second) per lap
	maxTime := baseTime - reduction
	if maxTime < 120 { // Minimum 2 seconds
		maxTime = 120
	}
	return maxTime
}

// calculateSpeedForStage calculates the required speed for a given stage
func calculateSpeedForStage(stage int) float64 {
	baseSpeed := 5.0
	increment := 0.5
	return baseSpeed + float64(stage-1)*increment
}

// drawGameUI draws the common game UI elements (game info)
func drawGameUI(sc *Scene, screen *ebiten.Image) {
	cfg := sc.uiConfig

	// draw go test -v style log (top-aligned, flows downward)
	logX := boundaryMargin + 10
	lineHeight := 20
	maxVisibleResults := 3

	logTopY := cfg.Stats.Y

	// Show recent test results (top to bottom, oldest first)
	log := sc.gameData.testLog
	start := 0
	if len(log) > maxVisibleResults {
		start = len(log) - maxVisibleResults
	}
	visible := log[start:]
	for i, r := range visible {
		var label string
		var c color.RGBA
		if r.passed {
			label = fmt.Sprintf("--- PASS: %s (%.2fs)", r.name, r.sec)
			c = color.RGBA{R: 100, G: 255, B: 100, A: 255}
		} else {
			label = fmt.Sprintf("--- FAIL: %s (%.2fs)", r.name, r.sec)
			c = color.RGBA{R: 255, G: 80, B: 80, A: 255}
		}
		y := logTopY + i*lineHeight
		t := libui.Text{Text: label, Face: cfg.Status.Face, Color: c}
		t.Image().Draw(screen, image.Point{X: logX, Y: y})
	}

	// Current running test (below results)
	runY := logTopY + len(visible)*lineHeight
	runLabel := fmt.Sprintf("=== RUN   TestLap%d", sc.gameData.laps+1)
	runText := libui.Text{
		Text:  runLabel,
		Face:  cfg.Status.Face,
		Color: cfg.Status.RunColor,
	}
	runTextImage := runText.Image()
	runTextImage.Draw(screen, image.Point{X: logX, Y: runY})

	// draw goal marker text with pulsing alpha
	pulse := math.Sin(float64(sc.gameData.lapTimer) * 0.1)
	var goalLabel string
	if sc.gameData.goalSide == -1 {
		goalLabel = "goto L"
	} else {
		goalLabel = "goto R"
	}
	textAlpha := uint8(160 + 95*pulse) // 65..255
	goalColor := color.RGBA{R: cfg.Goal.Color.R, G: cfg.Goal.Color.G, B: cfg.Goal.Color.B, A: textAlpha}
	goalMarker := libui.Text{
		Text:  goalLabel,
		Face:  cfg.Goal.Face,
		Color: goalColor,
	}
	goalMarkerImage := goalMarker.Image()
	markerW := goalMarkerImage.Rectangle().Dx()
	markerH := goalMarkerImage.Rectangle().Dy()
	// Fixed position: as if maxVisibleResults+1 lines are always shown
	fixedRunLineBottom := logTopY + (maxVisibleResults+1)*lineHeight
	goalY := fixedRunLineBottom - markerH
	if sc.gameData.goalSide == -1 {
		goalMarkerImage.Draw(screen, image.Point{X: boundaryMargin/2 - markerW/2, Y: goalY})
	} else {
		goalMarkerImage.Draw(screen, image.Point{X: ScreenWidth - boundaryMargin/2 - markerW/2, Y: goalY})
	}

	// draw scale count (7..0) in the top-right corner; hidden before first note
	scaleCount := sc.gameData.scaleCount
	if scaleCount <= 7 {
		timerColor := cfg.Timer.ColorNormal
		if scaleCount <= cfg.Timer.DangerThreshold {
			timerColor = cfg.Timer.ColorDanger
		} else if scaleCount <= cfg.Timer.WarnThreshold {
			timerColor = cfg.Timer.ColorWarn
		}
		timerText := libui.Text{
			Text:  fmt.Sprintf("%d", scaleCount),
			Face:  cfg.Timer.Face,
			Color: timerColor,
		}
		timerTextImage := timerText.Image()
		timerX := ScreenWidth - timerTextImage.Rectangle().Dx() - 20
		timerTextImage.Draw(screen, image.Point{X: timerX, Y: cfg.Timer.Y})
	}
}
