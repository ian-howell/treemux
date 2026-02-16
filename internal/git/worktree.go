// Package git provides minimal git helpers for treemux.
package git

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// RepositoryRoot returns the repository root directory for the given path.
func RepositoryRoot(dir string) (string, error) {
	output, err := runGit(dir, "rev-parse", "--show-toplevel")
	if err != nil {
		return "", fmt.Errorf("not inside a git repository")
	}
	root := strings.TrimSpace(output)
	if root == "" {
		return "", fmt.Errorf("not inside a git repository")
	}
	return root, nil
}

// EnsureWorktree creates a worktree if it does not already exist.
func EnsureWorktree(repoRoot, branch, path string) error {
	if repoRoot == "" || branch == "" || path == "" {
		return fmt.Errorf("missing worktree parameters")
	}
	if info, err := os.Stat(path); err == nil && info.IsDir() {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("failed to create worktree parent dir: %w", err)
	}
	if _, err := runGit(repoRoot, "worktree", "add", path, branch); err != nil {
		return fmt.Errorf("failed to create worktree: %w", err)
	}
	return nil
}

// runGit executes git with the provided arguments.
func runGit(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	if dir != "" {
		cmd.Dir = dir
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		errText := strings.TrimSpace(stderr.String())
		if errText == "" {
			errText = err.Error()
		}
		return "", fmt.Errorf("git %s: %s", strings.Join(args, " "), errText)
	}
	return stdout.String(), nil
}
