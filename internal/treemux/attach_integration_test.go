package treemux

import (
	"os"
	"path/filepath"
	"testing"
)

func newTestApp() *App {
	app := New()
	app.attach = func(string) error { return nil }
	return app
}

func TestAttachRootSetsLastAttached(t *testing.T) {
	app := newTestApp()
	rootDir := os.TempDir()
	name := "treemux-attach-root"
	cleanupSession(t, name)
	t.Cleanup(func() {
		cleanupSession(t, name)
	})

	if err := app.tmux.NewSession(name, rootDir, nil); err == nil {
		_ = app.tmux.SetOption(name, rootNameOption, name)
		_ = app.tmux.SetOption(name, rootDirOption, rootDir)
	}

	if err := app.AttachRoot(AttachRootRequest{Name: name, Dir: rootDir}); err != nil {
		t.Fatalf("attach-root failed: %v", err)
	}

	value, err := app.tmux.ShowOption(name, lastAttachedOption)
	if err != nil {
		t.Fatalf("failed to read last attached: %v", err)
	}
	if value == "" {
		t.Fatalf("expected last attached to be set")
	}
}

func TestAttachChildSetsLastAttached(t *testing.T) {
	app := newTestApp()
	rootDir := os.TempDir()
	root := "treemux-attach-root"
	child := "treemux-attach-child"
	childSession := root + tmuxSeparator() + child
	cleanupSession(t, childSession)
	cleanupSession(t, root)
	t.Cleanup(func() {
		cleanupSession(t, childSession)
	})
	t.Cleanup(func() {
		cleanupSession(t, root)
	})

	if err := app.tmux.NewSession(root, rootDir, nil); err != nil {
		t.Fatalf("failed to create root session: %v", err)
	}
	if err := app.tmux.SetOption(root, rootNameOption, root); err != nil {
		t.Fatalf("failed to set root name: %v", err)
	}
	if err := app.tmux.SetOption(root, rootDirOption, rootDir); err != nil {
		t.Fatalf("failed to set root dir: %v", err)
	}

	if err := app.AttachChild(AttachChildRequest{Root: root, Name: child}); err != nil {
		t.Fatalf("attach-child failed: %v", err)
	}

	childValue, err := app.tmux.ShowOption(childSession, lastAttachedOption)
	if err != nil {
		t.Fatalf("failed to read child last attached: %v", err)
	}
	if childValue == "" {
		t.Fatalf("expected child last attached to be set")
	}
	rootValue, err := app.tmux.ShowOption(root, lastAttachedOption)
	if err != nil {
		t.Fatalf("failed to read root last attached: %v", err)
	}
	if rootValue == "" {
		t.Fatalf("expected root last attached to be set")
	}
}

func TestAttachRootFromCwdDerivesNameAndDir(t *testing.T) {
	app := newTestApp()
	rootDir := t.TempDir()
	current, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}
	if err := os.Chdir(rootDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(current)
	})

	rootName := rootNameFromDir(rootDir)
	cleanupSession(t, rootName)
	t.Cleanup(func() {
		cleanupSession(t, rootName)
	})

	if err := app.AttachRoot(AttachRootRequest{}); err != nil {
		t.Fatalf("attach-root failed: %v", err)
	}

	if ok, err := app.tmux.HasSession(rootName); err != nil || !ok {
		t.Fatalf("expected root session %q to exist", rootName)
	}
	storedRoot, err := app.tmux.ShowOption(rootName, rootDirOption)
	if err != nil {
		t.Fatalf("failed to read root dir option: %v", err)
	}
	if storedRoot != rootDir {
		t.Fatalf("expected root dir %q, got %q", rootDir, storedRoot)
	}
}

func TestAttachRootWithDirDerivesName(t *testing.T) {
	app := newTestApp()
	rootDir := t.TempDir()
	rootName := rootNameFromDir(rootDir)
	cleanupSession(t, rootName)
	t.Cleanup(func() {
		cleanupSession(t, rootName)
	})

	if err := app.AttachRoot(AttachRootRequest{Dir: rootDir}); err != nil {
		t.Fatalf("attach-root failed: %v", err)
	}

	if ok, err := app.tmux.HasSession(rootName); err != nil || !ok {
		t.Fatalf("expected root session %q to exist", rootName)
	}
}

func TestAttachRootWithDirAndNameOverrides(t *testing.T) {
	app := newTestApp()
	rootDir := t.TempDir()
	rootName := "work"
	cleanupSession(t, rootName)
	t.Cleanup(func() {
		cleanupSession(t, rootName)
	})

	if err := app.AttachRoot(AttachRootRequest{Name: rootName, Dir: rootDir}); err != nil {
		t.Fatalf("attach-root failed: %v", err)
	}

	if ok, err := app.tmux.HasSession(rootName); err != nil || !ok {
		t.Fatalf("expected root session %q to exist", rootName)
	}
	storedRoot, err := app.tmux.ShowOption(rootName, rootDirOption)
	if err != nil {
		t.Fatalf("failed to read root dir option: %v", err)
	}
	if storedRoot != rootDir {
		t.Fatalf("expected root dir %q, got %q", rootDir, storedRoot)
	}
}

func TestAttachRootWithWorktreeCreatesDir(t *testing.T) {
	app := newTestApp()
	repo := newTestRepo(t)
	current, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}
	if err := os.Chdir(repo); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(current)
	})

	branch := "feature-x"
	rootDir := filepath.Join(repo, ".worktrees", branch)
	rootName := rootNameFromDir(rootDir)
	cleanupSession(t, rootName)
	t.Cleanup(func() {
		cleanupSession(t, rootName)
	})

	if err := app.AttachRoot(AttachRootRequest{Worktree: branch}); err != nil {
		t.Fatalf("attach-root failed: %v", err)
	}
	if info, err := os.Stat(rootDir); err != nil || !info.IsDir() {
		t.Fatalf("expected worktree dir %q to exist", rootDir)
	}
	if ok, err := app.tmux.HasSession(rootName); err != nil || !ok {
		t.Fatalf("expected root session %q to exist", rootName)
	}
}

func TestAttachRootWithWorktreeAndNameOverrides(t *testing.T) {
	app := newTestApp()
	repo := newTestRepo(t)
	current, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}
	if err := os.Chdir(repo); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(current)
	})

	branch := "feature-x"
	rootDir := filepath.Join(repo, ".worktrees", branch)
	rootName := "work"
	cleanupSession(t, rootName)
	t.Cleanup(func() {
		cleanupSession(t, rootName)
	})

	if err := app.AttachRoot(AttachRootRequest{Name: rootName, Worktree: branch}); err != nil {
		t.Fatalf("attach-root failed: %v", err)
	}
	if info, err := os.Stat(rootDir); err != nil || !info.IsDir() {
		t.Fatalf("expected worktree dir %q to exist", rootDir)
	}
	if ok, err := app.tmux.HasSession(rootName); err != nil || !ok {
		t.Fatalf("expected root session %q to exist", rootName)
	}
}
