package treemux

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ian-howell/treemux/internal/git"
)

// resolveWorktree ensures a worktree exists and returns its path.
func resolveWorktree(branch string) (string, error) {
	if branch == "" {
		return "", fmt.Errorf("missing worktree branch")
	}
	root, err := git.RepositoryRoot("")
	if err != nil {
		return "", err
	}
	name := strings.TrimSpace(branch)
	if name == "" {
		return "", fmt.Errorf("missing worktree branch")
	}
	path := filepath.Join(root, ".worktrees", name)
	if err := git.EnsureWorktree(root, name, path); err != nil {
		return "", err
	}
	return path, nil
}

// worktreeDefaultPath returns the default worktree path for a branch.
func worktreeDefaultPath(branch string) (string, error) {
	root, err := git.RepositoryRoot("")
	if err != nil {
		return "", err
	}
	name := strings.TrimSpace(branch)
	if name == "" {
		return "", fmt.Errorf("missing worktree branch")
	}
	return filepath.Join(root, ".worktrees", name), nil
}
