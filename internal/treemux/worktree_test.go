package treemux

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestWorktreeDefaultPath(t *testing.T) {
	root, err := worktreeDefaultPath("main")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if filepath.Base(root) != "main" {
		t.Fatalf("expected worktree path to end with branch name, got %q", root)
	}
	if !strings.Contains(root, string(filepath.Separator)+".worktrees"+string(filepath.Separator)) {
		t.Fatalf("expected worktree path to include .worktrees, got %q", root)
	}
}
