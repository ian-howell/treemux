package treemux

import (
	"os"
	"testing"
)

func TestAttachRootRejectsDirAndWorktree(t *testing.T) {
	app := New()
	err := app.AttachRoot(AttachRootRequest{
		Name:     "example",
		Dir:      "/tmp/example",
		Worktree: "branch",
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestResolveRootDirDefaultsToWorktreePath(t *testing.T) {
	repo := newTestRepo(t)
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}
	if err := os.Chdir(repo); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(cwd)
	})
	app := New()
	result, err := app.resolveRootDir("", "main")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == "" {
		t.Fatalf("expected worktree default path, got empty string")
	}
}

func TestResolveChildRootRequiresRootOutsideTmux(t *testing.T) {
	app := New()
	app.isInsideTmux = func() bool { return false }
	root, err := app.resolveChildRoot("")
	if err == nil {
		t.Fatalf("expected error, got root %q", root)
	}
}
