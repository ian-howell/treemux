// Package tmux tests integration with tmux.
package tmux

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestMain skips integration tests unless TREEMUX_TMUX is set.
func TestMain(m *testing.M) {
	if os.Getenv("TREEMUX_TMUX") == "" {
		os.Exit(0)
	}
	os.Exit(m.Run())
}

// TestAttachOrSwitch validates attach/switch for a test session.
func TestAttachOrSwitch(t *testing.T) {
	client := New()
	name := "treemux-test-attach"

	cleanupSession(t, name)
	if err := client.NewSession(name, os.TempDir(), nil); err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	defer cleanupSession(t, name)

	if err := client.AttachOrSwitch(name); err != nil {
		t.Fatalf("attach/switch failed: %v", err)
	}
}

// TestCurrentSessionName validates access to the attached session name.
func TestCurrentSessionName(t *testing.T) {
	client := New()
	name, err := client.CurrentSessionName()
	if err != nil {
		t.Fatalf("failed to get current session name: %v", err)
	}
	if strings.TrimSpace(name) == "" {
		t.Fatalf("expected current session name, got empty")
	}
}

// cleanupSession removes a session if it exists.
func cleanupSession(t *testing.T, name string) {
	cmd := exec.Command("tmux", "kill-session", "-t", name)
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			stderr := strings.TrimSpace(string(exitErr.Stderr))
			if strings.Contains(stderr, "can't find session") {
				return
			}
		}
	}
}
