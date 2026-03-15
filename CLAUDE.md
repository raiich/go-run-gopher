# CLAUDE.md

2D shuttle run game built with Ebitengine. Run with `make run`, test with `make test`.

## Package structure

Dependency: `gopher/` → `internal/` → `lib/`

| Package | Role |
|---|---|
| `gopher` | Entry point |
| `internal/game` | Game scene, state machine, constants |
| `internal/game/gopher` | Character movement, bounce physics |
| `internal/game/pregame` | Countdown overlay |
| `internal/game/audio` | Scale playback (sine wave) |
| `internal/game/effect` | Text effects (panic/recover) |
| `internal/game/proverbs` | Go Proverbs (result screen) |
| `internal/game/ui` | Game-specific UI |
| `internal/credits` | Credits scene |
| `internal/screen` | Screen size constants |
| `lib/scene` | Scene framework (ebiten.Game, stack) |
| `lib/ui` | UI library (button, text, layout) |

## Design notes

- All game behavior is organized as states (Entry/handleInput/update/draw)
- `eventloop.Dispatcher` syncs time management with the 60fps game loop
- New behavior with distinct modes/phases should be a state in the relevant state machine
- Init errors: `must.NoError()` (panic). Runtime errors: handle `machine.Trigger()` return

## Docs

- Game spec → `docs/game-design.md`
