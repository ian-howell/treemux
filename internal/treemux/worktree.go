package treemux

import (
	"fmt"
	"path/filepath"

	"github.com/ian-howell/treemux/internal/git"
)

// resolveWorktree ensures a worktree exists and returns its path.
func resolveWorktree(branch string) (string, error) {
	if branch == "" {
		return "", fmt.Errorf("missing worktree branch")
	}
	if isWhitespaceOnly(branch) {
		return "", fmt.Errorf("worktree branch cannot be blank")
	}
	root, err := git.RepositoryRoot("")
	if err != nil {
		return "", err
	}
	name := branch
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
	if branch == "" {
		return "", fmt.Errorf("missing worktree branch")
	}
	if isWhitespaceOnly(branch) {
		return "", fmt.Errorf("worktree branch cannot be blank")
	}
	name := branch
	return filepath.Join(root, ".worktrees", name), nil
}
