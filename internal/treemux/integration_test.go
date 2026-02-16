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
	if name == "" {
		return
	}
	cmd := exec.Command("tmux", "kill-session", "-t", name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		stderr := string(output)
		if strings.Contains(stderr, "can't find session") {
			return
		}
		if strings.Contains(stderr, "failed to connect to server") || strings.Contains(stderr, "no server running") {
			return
		}
		if stderr != "" {
			t.Fatalf("failed to kill session %q: %s", name, stderr)
		}
		t.Fatalf("failed to kill session %q: %v", name, err)
	}
}

func TestIntegration(t *testing.T) {
	if os.Getenv("TREEMUX_INTEGRATION") == "" {
		t.Skip("requires TREEMUX_INTEGRATION")
	}

	t.Run("show-roots-mru", func(t *testing.T) {
		client := tmux.New()
		app := New()
		app.attach = func(string) error { return nil }
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
	})

	t.Run("show-children-output", func(t *testing.T) {
		client := tmux.New()
		app := New()
		app.attach = func(string) error { return nil }
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
	})

	t.Run("show-children-hide-current", func(t *testing.T) {
		client := tmux.New()
		app := New()
		app.attach = func(string) error { return nil }
		root := "treemux-hide-current-root"
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
		app.isInsideTmux = func() bool { return true }
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
		app.currentSessionName = func() (string, error) {
			return childB, nil
		}

		lines, err := app.ShowChildren(ShowChildrenRequest{Root: root, HideCurrent: true})
		if err != nil {
			t.Fatalf("show-children failed: %v", err)
		}
		for _, line := range lines {
			if strings.Contains(line, "beta") {
				t.Fatalf("expected current child to be hidden, got %v", lines)
			}
		}
		cleanupSession(t, childA)
		cleanupSession(t, childB)
		cleanupSession(t, root)
	})

	t.Run("attach-child-current-root", func(t *testing.T) {
		client := tmux.New()
		app := New()
		app.attach = func(string) error { return nil }
		app.isInsideTmux = func() bool { return true }
		app.currentSessionName = func() (string, error) {
			return "treemux-current-root", nil
		}
		root := "treemux-current-root"
		rootDir := os.TempDir()
		separator := tmuxSeparator()
		childName := "nvim"
		childSession := root + separator + childName

		cleanupSession(t, childSession)
		cleanupSession(t, root)
		t.Cleanup(func() {
			cleanupSession(t, childSession)
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
		if err := app.AttachChild(AttachChildRequest{Name: childName}); err != nil {
			t.Fatalf("attach-child failed: %v", err)
		}
		if ok, err := client.HasSession(childSession); err != nil || !ok {
			t.Fatalf("expected child session %q to exist", childSession)
		}
		cleanupSession(t, childSession)
		cleanupSession(t, root)
	})
}
