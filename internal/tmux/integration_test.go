// Package tmux tests integration with tmux.
package tmux

import (
	"errors"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestMain skips integration tests unless TREEMUX_INTEGRATION is set.
func TestMain(m *testing.M) {
	if os.Getenv("TREEMUX_INTEGRATION") == "" {
		os.Exit(0)
	}
	os.Exit(m.Run())
}

// TestSessionLifecycle validates tmux session operations without attaching.
func TestSessionLifecycle(t *testing.T) {
	client := New()
	name := "treemux-test-session"

	cleanupSession(t, name)
	if err := client.NewSession(name, os.TempDir(), nil); err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	t.Cleanup(func() {
		cleanupSession(t, name)
	})

	exists, err := client.HasSession(name)
	if err != nil {
		t.Fatalf("failed to check session: %v", err)
	}
	if !exists {
		t.Fatalf("expected session to exist")
	}

	if err := client.SetOption(name, "@treemux_test", "1"); err != nil {
		t.Fatalf("failed to set option: %v", err)
	}
	value, err := client.ShowOption(name, "@treemux_test")
	if err != nil {
		t.Fatalf("failed to read option: %v", err)
	}
	if value != "1" {
		t.Fatalf("expected option value, got %q", value)
	}
}

// cleanupSession removes a session if it exists.
func cleanupSession(t *testing.T, name string) {
	cmd := exec.Command("tmux", "kill-session", "-t", name)
	if err := cmd.Run(); err != nil {
		if exitErr, ok := errors.AsType[*exec.ExitError](err); ok {
			stderr := string(exitErr.Stderr)
			if strings.Contains(stderr, "can't find session") {
				return
			}
		}
	}
}
