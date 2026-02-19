// Package treemux contains core treemux behavior.
package treemux

import (
	"fmt"

	"github.com/ian-howell/treemux/internal/models"
)

type Session struct {
	Attacher

	models.Session
}

type Attacher interface {
	Attach() error
}

type SesssionAttacher interface {
	Attach(name string) error
}

type Lister interface {
	List() ([]Session, error)
	SesssionAttacher
}

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
func New(opts ...Option) *App {
	app := &App{}

	for _, opt := range opts {
		opt(app)
	}

	return app
}

func (a *App) Run() error {
	if a.prompter == nil {
		return fmt.Errorf("no prompter configured")
	}

	sessions, err := a.listSessions()
	if err != nil {
		return fmt.Errorf("failed to list sessions: %w", err)
	}

	session, err := a.prompter.Prompt(sessions)
	if err != nil {
		return fmt.Errorf("failed to prompt for session: %w", err)
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
