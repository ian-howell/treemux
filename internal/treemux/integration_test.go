package treemux

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/ian-howell/treemux/internal/tmux"
)

func cleanupSession(t *testing.T, name string) {
	t.Helper()
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

func TestShowRootsMostRecentlyUsedIntegration(t *testing.T) {
	client := tmux.New()
	app := New()
	rootA := "treemux-root-a"
	rootB := "treemux-root-b"
	rootDir := os.TempDir()

	cleanupSession(t, rootA)
	cleanupSession(t, rootB)
	t.Cleanup(func() {
		cleanupSession(t, rootA)
	})
	t.Cleanup(func() {
		cleanupSession(t, rootB)
	})

	if err := client.NewSession(rootA, rootDir, nil); err != nil {
		t.Fatalf("failed to create root session: %v", err)
	}
	if err := client.NewSession(rootB, rootDir, nil); err != nil {
		t.Fatalf("failed to create root session: %v", err)
	}
	if err := client.SetOption(rootA, rootNameOption, rootA); err != nil {
		t.Fatalf("failed to set root name option: %v", err)
	}
	if err := client.SetOption(rootA, rootDirOption, rootDir); err != nil {
		t.Fatalf("failed to set root dir option: %v", err)
	}
	if err := client.SetOption(rootB, rootNameOption, rootB); err != nil {
		t.Fatalf("failed to set root name option: %v", err)
	}
	if err := client.SetOption(rootB, rootDirOption, rootDir); err != nil {
		t.Fatalf("failed to set root dir option: %v", err)
	}
	if err := client.SetOption(rootA, lastAttachedOption, fmt.Sprintf("%d", time.Now().Add(-time.Minute).UnixNano())); err != nil {
		t.Fatalf("failed to set attached time: %v", err)
	}
	if err := client.SetOption(rootB, lastAttachedOption, fmt.Sprintf("%d", time.Now().UnixNano())); err != nil {
		t.Fatalf("failed to set attached time: %v", err)
	}

	lines, err := app.ShowRoots(ShowRootsRequest{SortBy: "most-recently-used"})
	if err != nil {
		t.Fatalf("show-roots failed: %v", err)
	}
	if len(lines) == 0 || !strings.Contains(lines[0], rootB) {
		t.Fatalf("expected most recent root first, got %v", lines)
	}
}

func TestShowChildrenOutputIntegration(t *testing.T) {
	client := tmux.New()
	app := New()
	root := "treemux-root"
	rootDir := os.TempDir()
	separator := tmuxSeparator()
	childA := root + separator + "alpha"
	childB := root + separator + "beta"

	cleanupSession(t, childA)
	cleanupSession(t, childB)
	cleanupSession(t, root)
	t.Cleanup(func() {
		cleanupSession(t, childA)
	})
	t.Cleanup(func() {
		cleanupSession(t, childB)
	})
	t.Cleanup(func() {
		cleanupSession(t, root)
	})

	if err := client.NewSession(root, rootDir, nil); err != nil {
		t.Fatalf("failed to create root session: %v", err)
	}
	if err := client.SetOption(root, rootNameOption, root); err != nil {
		t.Fatalf("failed to set root name option: %v", err)
	}
	if err := client.SetOption(root, rootDirOption, rootDir); err != nil {
		t.Fatalf("failed to set root dir option: %v", err)
	}
	if err := client.NewSession(childA, rootDir, nil); err != nil {
		t.Fatalf("failed to create child session: %v", err)
	}
	if err := client.SetOption(childA, rootNameOption, root); err != nil {
		t.Fatalf("failed to set child root name option: %v", err)
	}
	if err := client.SetOption(childA, rootDirOption, rootDir); err != nil {
		t.Fatalf("failed to set child root dir option: %v", err)
	}
	if err := client.NewSession(childB, rootDir, nil); err != nil {
		t.Fatalf("failed to create child session: %v", err)
	}
	if err := client.SetOption(childB, rootNameOption, root); err != nil {
		t.Fatalf("failed to set child root name option: %v", err)
	}
	if err := client.SetOption(childB, rootDirOption, rootDir); err != nil {
		t.Fatalf("failed to set child root dir option: %v", err)
	}
	if err := client.SetOption(childA, lastAttachedOption, fmt.Sprintf("%d", time.Now().Add(-time.Minute).UnixNano())); err != nil {
		t.Fatalf("failed to set attached time: %v", err)
	}
	if err := client.SetOption(childB, lastAttachedOption, fmt.Sprintf("%d", time.Now().UnixNano())); err != nil {
		t.Fatalf("failed to set attached time: %v", err)
	}

	lines, err := app.ShowChildren(ShowChildrenRequest{Root: root, SortBy: "most-recently-used"})
	if err != nil {
		t.Fatalf("show-children failed: %v", err)
	}
	if len(lines) == 0 {
		t.Fatalf("expected child output, got empty")
	}
	for _, line := range lines {
		if strings.Contains(line, root) {
			t.Fatalf("expected child-only output, got %v", lines)
		}
		if strings.Contains(line, "🌿") {
			t.Fatalf("expected no separator output, got %v", lines)
		}
	}
	if !strings.Contains(lines[0], "beta") {
		t.Fatalf("expected most recent child first, got %v", lines)
	}
}
