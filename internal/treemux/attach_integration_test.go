package treemux

import (
	"os"
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
