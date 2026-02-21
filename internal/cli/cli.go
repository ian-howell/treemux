// Package cli defines the command-line interface for treemux.
//
// It parses CLI arguments, loads configuration, and runs the treemux application.
package cli

import (
	"fmt"

	"github.com/ian-howell/treemux/internal/listers"
	"github.com/ian-howell/treemux/internal/prompters"
	"github.com/ian-howell/treemux/internal/tmux"
	"github.com/ian-howell/treemux/internal/treemux"
)

// Run parses CLI args and executes the requested command.
func Run(config Config) error {
	// TODO: Handle config
	tmuxClient := tmux.New()
	app, err := treemux.New(
		treemux.WithPrompter(&prompters.Huh{
			FullScreen: config.FullScreen,
		}),
		treemux.WithListers([]treemux.Lister{
			listers.NewActiveSessions(tmuxClient),
		}),
	)
	if err != nil {
		return fmt.Errorf("creating app: %w", err)
	}
	return app.Run()
}
