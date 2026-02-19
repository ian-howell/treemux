// Package cli defines the command-line interface for treemux.
//
// It parses CLI arguments, loads configuration, and runs the treemux application.
package cli

import (
	"github.com/ian-howell/treemux/internal/sessionizers"
	"github.com/ian-howell/treemux/internal/tmux"
	"github.com/ian-howell/treemux/internal/treemux"
)

// Run parses CLI args and executes the requested command.
func Run(config Config) error {
	// TODO: Handle config
	tmuxClient := tmux.New()
	return treemux.New(
		treemux.WithSessionizers([]treemux.Sessionizer{
			sessionizers.NewActiveSessions(tmuxClient),
		}),
	).Run()
}
