// Package treemux contains core treemux behavior.
package treemux

import (
	"fmt"
)

type Session interface {
	Attach() error
	fmt.Stringer
}

type Lister interface {
	List() ([]Session, error)
}

// The prompter is needed to provide a UI for the user to select a Session.
// Callers should assume that a returned nil Session implies that the user canceled the prompt
type prompter interface {
	Prompt(sessions []Session) (Session, error)
}

// App bundles core treemux dependencies.
type App struct {
	// listers provide sessions to display and attach to.
	listers []Lister

	// prompter provides session selection UI.
	prompter prompter
}

type Option func(*App)

func WithListers(listers []Lister) Option {
	return func(app *App) {
		app.listers = listers
	}
}

func WithPrompter(prompter prompter) Option {
	return func(app *App) {
		app.prompter = prompter
	}
}

// New returns a new App instance.
func New(opts ...Option) (*App, error) {
	app := &App{}

	for _, opt := range opts {
		opt(app)
	}

	if app.prompter == nil {
		return nil, fmt.Errorf("no prompter configured")
	}

	if len(app.listers) == 0 {
		return nil, fmt.Errorf("no listers configured")
	}

	return app, nil
}

func (a *App) Run() error {
	sessions, err := a.listSessions()
	if err != nil {
		return fmt.Errorf("failed to list sessions: %w", err)
	}

	session, err := a.prompter.Prompt(sessions)
	if err != nil {
		return fmt.Errorf("failed to prompt for session: %w", err)
	}
	if session == nil {
		// Assume that the user cancelled the selection
		return nil
	}

	if err := session.Attach(); err != nil {
		return fmt.Errorf("failed to attach to session: %w", err)
	}

	return nil
}

func (a *App) listSessions() ([]Session, error) {
	// TODO: Handle duplicates and sorting
	var allSessions []Session
	for _, lister := range a.listers {
		sessions, err := lister.List()
		if err != nil {
			return nil, err
		}
		allSessions = append(allSessions, sessions...)
	}
	return allSessions, nil
}
