package prompters

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/ian-howell/treemux/internal/treemux"
)

type Huh struct{}

type sessionChoice struct {
	label   string
	session treemux.Session
}

func (c sessionChoice) String() string {
	return c.label
}

func NewHuh() *Huh {
	return &Huh{}
}

func (p *Huh) Prompt(sessions []treemux.Session) (treemux.Session, error) {
	if len(sessions) == 0 {
		return nil, fmt.Errorf("no sessions available")
	}

	choices := make([]huh.Option[sessionChoice], 0, len(sessions))
	for _, session := range sessions {
		label := session.String()
		choices = append(choices, huh.NewOption(label, sessionChoice{
			label:   label,
			session: session,
		}))
	}

	selected := sessionChoice{}
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[sessionChoice]().
				Title("Select a session").
				Options(choices...).
				Value(&selected),
		),
	)

	if err := form.Run(); err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return nil, fmt.Errorf("prompt canceled")
		}
		return nil, err
	}

	return selected.session, nil
}
