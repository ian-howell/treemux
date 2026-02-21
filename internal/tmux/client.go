// Package tmux wraps tmux interactions for treemux.
package tmux

import (
	"fmt"

	gotmux "github.com/jubnzv/go-tmux"
)

// Client provides tmux operations used by treemux.
type Client struct{}

// New returns a new tmux client.
func New() *Client {
	return &Client{}
}

// RunCmd runs a tmux command and returns its output.
func (c *Client) RunCmd(args []string) (stdout string, err error) {
	stdout, stderr, err := gotmux.RunCmd(args)
	if err != nil {
		return "", fmt.Errorf("tmux command failed: %w", err)
	}

	if stderr != "" {
		return "", fmt.Errorf("tmux error: %s", stderr)
	}

	return stdout, nil
}

// Attach attaches to the named session, creating it if it doesn't exist.
func (c *Client) AttachOrSwitch(name string) error {
	return (&gotmux.Session{Name: name}).AttachSession()
}
