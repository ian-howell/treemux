package treemux

import (
	"os"
	"testing"
)

func TestResolveRootDirPrefersDir(t *testing.T) {
	app := New()
	result, err := app.resolveRootDir("/tmp/example", "branch")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result != "/tmp/example" {
		t.Fatalf("expected dir to be returned, got %q", result)
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
	root, err := app.resolveChildRoot("")
	if err == nil {
		t.Fatalf("expected error, got root %q", root)
	}
}
