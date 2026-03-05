# go run ./gopher

<p align="center">[ <b>English</b> | <a href="README.ja.md">日本語</a> ]</p>

A 2D shuttle run game starring Go's mascot "Gopher".

Tap `go` to run, hit the wall to `panic`, and bounce back to `recover`. Results scroll in `go test -v` style, ending with a Go Proverb.

![Gameplay screenshot](docs/images/running.png)

## How to Play

- Tap the `go` button to accelerate Gopher
- Reach the goal line before time runs out to complete a lap
- Every 7 laps the stage advances, increasing speed and difficulty
- Two consecutive misses and the game is over

Controls are touch-first (mobile-first). Falls back to mouse input when touch is unavailable.

## Build & Run

```bash
make run          # Run the game (native)
make wasm-dist    # Build WebAssembly distribution package
make serve        # Start local server (localhost:8080)
make test         # Run tests
```

### Prerequisites

- Go 1.25+
- [Ebiten v2.9.8](https://ebitengine.org/) system requirements

## Documentation

- [Game Design](docs/game-design.md) — Screens, controls, game rules, balance
- [Architecture](docs/architecture.md) — State machine, package structure, design patterns

## License

[MIT](LICENSE)
