// Package treemux contains core treemux behavior.
package treemux

import "github.com/ian-howell/treemux/internal/tmux"

// App bundles core treemux dependencies.
type App struct {
	// tmux provides session operations.
	tmux *tmux.Client
	// attach overrides session attachment for tests.
	attach func(string) error
	// isInsideTmux overrides tmux environment detection for tests.
	isInsideTmux func() bool
	// currentSessionName overrides tmux current session lookup for tests.
	currentSessionName func() (string, error)
}

// New returns a new App instance.
func New() *App {
	app := &App{tmux: tmux.New()}
	app.attach = app.tmux.AttachOrSwitch
	app.isInsideTmux = app.tmux.IsInsideTmux
	app.currentSessionName = app.tmux.CurrentSessionName
	return app
}

// AttachRootRequest holds options for root attachment.
type AttachRootRequest struct {
	// Name specifies the root session name.
	Name string
	// Dir sets the session start directory.
	Dir string
	// Worktree selects a git worktree branch.
	Worktree string
}

// AttachChildRequest holds options for child attachment.
type AttachChildRequest struct {
	// Root specifies the root session name.
	Root string
	// Name specifies the child session name.
	Name string
	// Command specifies a command to run in the session.
	Command string
}

// ShowRootsRequest holds options for root listing.
type ShowRootsRequest struct {
	// SortBy sets the ordering mode.
	SortBy string
	// HideCurrent controls whether to hide the current root.
	HideCurrent bool
}

// ShowChildrenRequest holds options for child listing.
type ShowChildrenRequest struct {
	// Root specifies the root session name.
	Root string
	// SortBy sets the ordering mode.
	SortBy string
	// HideCurrent controls whether to hide the current session.
	HideCurrent bool
}
