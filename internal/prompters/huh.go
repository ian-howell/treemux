package prompters

import (
	"errors"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/x/term"

	"github.com/ian-howell/treemux/internal/treemux"
)

type Huh struct {
	FullScreen bool
}

type sessionChoice struct {
	label   string
	session treemux.Session
}

func (c sessionChoice) String() string {
	return c.label
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

	if p.FullScreen {
		form = form.WithWidth(screenWidth()).WithHeight(screenHeight())
	}

	keymap := huh.NewDefaultKeyMap()
	keymap.Quit = key.NewBinding(key.WithKeys(
		"ctrl+c",
		"esc",
	))
	form = form.WithKeyMap(keymap)

	if err := form.Run(); err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			// If the user aborted, send a nil Session with a nil err to indicate that
			// no selection was made.
			return nil, nil
		}
		return nil, err
	}

	return selected.session, nil
}

// screenWidth returns the width of the terminal screen. If it cannot be determined, it returns a default width.
func screenWidth() int {
	// Huh uses stderr for its output, so we get the terminal size from stderr.
	width, _, err := term.GetSize(os.Stderr.Fd())
	if err != nil || width <= 0 {
		return 0
	}
	return width
}

// screenHeight returns the height of the terminal screen. If it cannot be determined, it returns a default height.
func screenHeight() int {
	// Huh uses stderr for its output, so we get the terminal size from stderr.
	_, height, err := term.GetSize(os.Stderr.Fd())
	if err != nil || height <= 0 {
		return 0
	}
	return height
}
