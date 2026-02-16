// Package treemux contains core treemux behavior.
package treemux

import "github.com/ian-howell/treemux/internal/tmux"

// App bundles core treemux dependencies.
type App struct {
	// tmux provides session operations.
	tmux *tmux.Client
}

// New returns a new App instance.
func New() *App {
	return &App{tmux: tmux.New()}
}

// AttachRootRequest holds options for root attachment.
type AttachRootRequest struct {
	// Name specifies the root session name.
	Name string
	// Command specifies a command to run in the session.
	Command string
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
