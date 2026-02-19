# treemux

treemux is a tmux session selector with a small, composable core. The runtime is built around three
roles: listers that discover sessions, prompters that present choices, and attachers that join the
selected session. Each role is a narrow interface so new behaviors can be added without rewriting
the rest of the system.

## Architecture overview

treemux flows through a simple pipeline:

1. The CLI wires dependencies and constructs the app.
2. Listers return candidate sessions.
3. The prompter renders a selection UI and returns a chosen session.
4. The chosen session attaches through its attacher.

Most of the code lives under `internal/` and is intentionally split by role. The app only cares
about interfaces, so the wiring in `internal/cli` is the main place where implementations are
selected.

## Listers

Listers enumerate sessions and return a slice of `treemux.Session`. They can draw from tmux state,
local metadata, or other sources without forcing the UI or attach logic to change.

Contract:

- `List() ([]treemux.Session, error)` returns sessions ready to display.
- Each session includes metadata for prompt display (name, last attached time, attached state).

Current implementation:

- `internal/listers/active_sessions.go` uses `tmux list-sessions` to collect active sessions and
  populate the core session model.

Design notes:

- Listers can return overlapping sessions; the app currently does not deduplicate.
- Listers are expected to tolerate tmux not running and return an empty list rather than failing.

## Prompters

Prompters present the session list and return a single `treemux.Session` selection. They should
focus on user interaction only, and avoid tmux-specific logic.

Contract:

- `Prompt([]treemux.Session) (treemux.Session, error)` returns the chosen session.
- If no sessions are available, the prompter should return a clear error.

Current implementation:

- `internal/prompters/huh.go` uses `github.com/charmbracelet/huh` to render a TUI selection list.
  Attached sessions are prefixed with `* ` in the label.

Design notes:

- Prompt cancellation is treated as a clean error (`prompt canceled`).
- The prompter is configured in `internal/cli/cli.go` via `treemux.WithPrompter`.

## Attachers

Attachers connect to the chosen session. The core `treemux.Session` embeds an `Attacher` interface,
so a lister can return sessions that know how to attach themselves.

Contract:

- `Attach() error` joins or switches to the session.
- Attachers should be safe to call even if the session is already attached.

Current state:

- The active sessions lister does not yet populate a concrete attacher. The interface exists so
  session models can carry attach behavior when implementations are added.

## Data flow

The app assembles dependencies in the CLI and runs a short-lived pipeline:

1. `treemux.New(...)` initializes the app with listers and a prompter.
2. `List()` is called on each lister and results are concatenated.
3. The prompter returns a `treemux.Session`.
4. `Session.Attach()` is invoked on the chosen session.

Errors bubble up with context, so callers can report where the pipeline failed.

## Extensibility

Adding new behavior typically means implementing one of the role interfaces and wiring it in the
CLI. For example, a future lister could provide session metadata from a config file, while a new
prompter might render a different TUI or a non-interactive selector.

## Development

Run tests:

```
go test ./...
```

Integration tests require tmux and must be run inside tmux:

```
TREEMUX_TMUX=1 go test ./...
```
