package listers

import (
	"fmt"
	"strings"

	"github.com/ian-howell/treemux/internal/models"
	"github.com/ian-howell/treemux/internal/treemux"
)

type tmuxClient interface {
	RunCmd(args []string) (stdout string, err error)
}

type ActiveSessions struct {
	tmuxClient tmuxClient
}

func NewActiveSessions(tmuxClient tmuxClient) *ActiveSessions {
	return &ActiveSessions{tmuxClient: tmuxClient}
}

// List returns all active sessions.
func (s *ActiveSessions) List() ([]treemux.Session, error) {
	args := []string{"list-sessions", "-F", "#{session_name} #{session_last_attached} #{?session_attached,true,false}"}

	output, err := s.tmuxClient.RunCmd(args)
	if err != nil {
		return []treemux.Session{}, nil
	}
	lines := strings.Split(strings.TrimSpace(output), "\n")

	sessions := make([]treemux.Session, 0, len(lines))
	for _, line := range lines {
		var name string
		var lastAttachedTime int64
		var attachedStr string
		fmt.Sscanf(line, "%s %d %s", &name, &lastAttachedTime, &attachedStr)
		attached := attachedStr == "true"
		sessions = append(sessions, treemux.Session{
			Attacher: &ActiveSessionAttacher{tmuxClient: s.tmuxClient},
			Session: models.Session{
				Name:             name,
				LastAttachedTime: lastAttachedTime,
				Attached:         attached,
			},
		})
	}
	return sessions, nil
}

type ActiveSessionAttacher struct {
	tmuxClient tmuxClient
}

func (a *ActiveSessionAttacher) Attach() error {
	return nil
}
