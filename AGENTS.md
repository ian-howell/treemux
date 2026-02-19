# treemux agents

treemux is organized around three runtime agents: listers, prompters, and attachers. Each agent is
an interface role with a focused responsibility. Implementations live under `internal/` and are
wired together by the CLI.

## Listers

Responsibility: enumerate candidate sessions.

Interface:

- `List() ([]treemux.Session, error)`

Notes:

- Listers can collect sessions from tmux, configuration, or other sources.
- Listers should handle missing tmux servers gracefully and return an empty list when appropriate.

Current implementation:

- `internal/listers/active_sessions.go` reads tmux state with `list-sessions` and builds session
  models.

## Prompters

Responsibility: present available sessions and return a user selection.

Interface:

- `Prompt([]treemux.Session) (treemux.Session, error)`

Notes:

- Prompters are UI-only. They should not attach or mutate tmux state.
- Prompters must handle empty input with a clear error.

Current implementation:

- `internal/prompters/huh.go` uses the `huh` TUI to render the selection list and returns the
  chosen session.

## Attachers

Responsibility: attach or switch to the selected session.

Interface:

- `Attach() error`

Notes:

- The core `treemux.Session` embeds an `Attacher`, allowing listers to return sessions that know
  how to attach themselves.
- Attachers should be safe to call when the session is already attached.

Current state:

- No concrete attachers are wired yet; the interface exists for upcoming implementations.

## Wiring

The CLI (`internal/cli/cli.go`) constructs the application by providing listers and the prompter:

- `treemux.WithListers(...)` sets the listers used for discovery.
- `treemux.WithPrompter(...)` sets the interactive selector.

The core app (`internal/treemux/app.go`) is intentionally dependency-free and only interacts with
the interfaces above.
