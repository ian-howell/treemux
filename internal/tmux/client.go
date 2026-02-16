// Package tmux wraps tmux interactions for treemux.
package tmux

import (
	"fmt"
	"strings"

	gotmux "github.com/jubnzv/go-tmux"
)

// Client provides tmux operations used by treemux.
type Client struct {
}

// New returns a new tmux client.
func New() *Client {
	return &Client{}
}

// HasSession reports whether the named session exists.
func (c *Client) HasSession(name string) (bool, error) {
	args := []string{"has-session", "-t", "=" + name}
	_, stderr, err := gotmux.RunCmd(args)
	if err != nil {
		if strings.Contains(stderr, "can't find session") {
			return false, nil
		}
		if strings.Contains(stderr, "failed to connect to server") || strings.Contains(stderr, "no server running") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// ListSessions returns all session names.
func (c *Client) ListSessions() ([]string, error) {
	args := []string{"list-sessions", "-F", "#{session_name}"}
	output, stderr, err := gotmux.RunCmd(args)
	if err != nil {
		if strings.Contains(stderr, "failed to connect to server") || strings.Contains(stderr, "no server running") {
			return []string{}, nil
		}
		return nil, err
	}
	output = strings.TrimSpace(output)
	if output == "" {
		return []string{}, nil
	}
	lines := strings.Split(output, "\n")
	names := make([]string, 0, len(lines))
	for _, line := range lines {
		if line == "" {
			continue
		}
		names = append(names, line)
	}
	return names, nil
}

// ShowOption reads a tmux option from the target session.
func (c *Client) ShowOption(target, name string) (string, error) {
	args := []string{"show-option", "-qv"}
	if target != "" {
		args = append(args, "-t", target)
	}
	args = append(args, name)
	output, _, err := gotmux.RunCmd(args)
	if err != nil {
		return "", fmt.Errorf("failed to read tmux option %s", name)
	}
	return strings.TrimSpace(output), nil
}

// SetOption writes a tmux option for the target session.
func (c *Client) SetOption(target, name, value string) error {
	args := []string{"set-option", "-q"}
	if target != "" {
		args = append(args, "-t", target)
	}
	args = append(args, name, value)
	_, _, err := gotmux.RunCmd(args)
	return err
}

// NewSession creates a detached tmux session.
func (c *Client) NewSession(name, dir string, command []string) error {
	args := []string{"new-session", "-d", "-s", name, "-c", dir}
	if len(command) > 0 {
		args = append(args, "--")
		args = append(args, command...)
	}
	_, _, err := gotmux.RunCmd(args)
	return err
}

// AttachOrSwitch attaches to or switches to the named session.
func (c *Client) AttachOrSwitch(name string) error {
	args := []string{"attach-session", "-t", name}
	if gotmux.IsInsideTmux() {
		args = []string{"switch-client", "-t", name}
	}
	return gotmux.ExecCmd(args)
}

// IsInsideTmux reports whether the current process runs inside tmux.
func (c *Client) IsInsideTmux() bool {
	return gotmux.IsInsideTmux()
}

// CurrentSessionName returns the attached tmux session name.
func (c *Client) CurrentSessionName() (string, error) {
	name, err := gotmux.GetAttachedSessionName()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(name), nil
}
