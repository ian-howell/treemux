package listers

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ian-howell/treemux/internal/models"
	"github.com/ian-howell/treemux/internal/treemux"
)

type tmuxClient interface {
	RunCmd(args []string) (stdout string, err error)
	AttachOrSwitch(name string) error
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

	sessions := make([]ActiveSession, 0, len(lines))
	for _, line := range lines {
		var name string
		var lastAttachedTime int64
		var attachedStr string
		fmt.Sscanf(line, "%s %d %s", &name, &lastAttachedTime, &attachedStr)
		attached := attachedStr == "true"
		sessions = append(sessions, ActiveSession{
			tmuxClient: s.tmuxClient,
			Session: models.Session{
				Name:             name,
				LastAttachedTime: lastAttachedTime,
				Attached:         attached,
			},
		})
	}

	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].Session.LastAttachedTime > sessions[j].Session.LastAttachedTime
	})

	treemuxSessions := make([]treemux.Session, 0, len(sessions))
	for _, session := range sessions {
		treemuxSessions = append(treemuxSessions, session)
	}

	return treemuxSessions, nil
}

type ActiveSession struct {
	tmuxClient tmuxClient
	Session    models.Session
}

// Attach attaches to the session, creating it if it doesn't exist.
func (a ActiveSession) Attach() error {
	return a.tmuxClient.AttachOrSwitch(a.Session.Name)
}

// String returns the session as a string for display in a prompter.
func (a ActiveSession) String() string {
	if a.Session.Attached {
		return fmt.Sprintf("* %s", a.Session.Name)
	}
	return fmt.Sprintf("  %s", a.Session.Name)
}
