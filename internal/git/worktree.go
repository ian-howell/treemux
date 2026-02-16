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
	if err := EnsureBranch(repoRoot, branch); err != nil {
		return err
	}
	if _, err := runGit(repoRoot, "worktree", "add", path, branch); err != nil {
		return fmt.Errorf("failed to create worktree: %w", err)
	}
	if err := ensureGitignoreWorktree(repoRoot); err != nil {
		return err
	}
	return nil
}

// EnsureBranch ensures a local branch exists, creating or tracking as needed.
func EnsureBranch(repoRoot, branch string) error {
	branch = strings.TrimSpace(branch)
	if repoRoot == "" || branch == "" {
		return fmt.Errorf("missing branch parameters")
	}
	if ok, err := hasLocalBranch(repoRoot, branch); err != nil {
		return err
	} else if ok {
		return nil
	}
	remotes, err := listRemotes(repoRoot)
	if err != nil {
		return err
	}
	if len(remotes) > 0 {
		if _, err := runGit(repoRoot, "fetch", "--all", "--prune"); err != nil {
			return err
		}
		remote, err := firstRemoteWithBranch(repoRoot, branch)
		if err != nil {
			return err
		}
		if remote != "" {
			if _, err := runGit(repoRoot, "branch", "--track", branch, remote+"/"+branch); err != nil {
				return err
			}
			return nil
		}
	}
	if _, err := runGit(repoRoot, "branch", branch); err != nil {
		return err
	}
	return nil
}

func hasLocalBranch(repoRoot, branch string) (bool, error) {
	_, err := runGit(repoRoot, "show-ref", "--verify", "refs/heads/"+branch)
	if err == nil {
		return true, nil
	}
	if strings.Contains(err.Error(), "not a valid ref") || strings.Contains(err.Error(), "fatal") {
		return false, nil
	}
	return false, err
}

func hasRemoteBranch(repoRoot, remote, branch string) (bool, error) {
	_, err := runGit(repoRoot, "show-ref", "--verify", "refs/remotes/"+remote+"/"+branch)
	if err == nil {
		return true, nil
	}
	if strings.Contains(err.Error(), "not a valid ref") || strings.Contains(err.Error(), "fatal") {
		return false, nil
	}
	return false, err
}

func listRemotes(repoRoot string) ([]string, error) {
	output, err := runGit(repoRoot, "remote")
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(output), "\n")
	remotes := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if line == "origin" {
			if _, err := runGit(repoRoot, "remote", "get-url", "origin"); err != nil {
				continue
			}
		}
		remotes = append(remotes, line)
	}
	return remotes, nil
}

func firstRemoteWithBranch(repoRoot, branch string) (string, error) {
	remotes, err := listRemotes(repoRoot)
	if err != nil {
		return "", err
	}
	if len(remotes) == 0 {
		return "", nil
	}
	if hasOrigin, err := hasRemoteBranch(repoRoot, "origin", branch); err != nil {
		return "", err
	} else if hasOrigin {
		return "origin", nil
	}
	for _, remote := range remotes {
		if remote == "origin" {
			continue
		}
		ok, err := hasRemoteBranch(repoRoot, remote, branch)
		if err != nil {
			return "", err
		}
		if ok {
			return remote, nil
		}
	}
	return "", nil
}

func ensureGitignoreWorktree(repoRoot string) error {
	const (
		ignoreEntry  = ".worktree"
		commentEntry = "# Added by treemux worktree creation"
	)
	path := filepath.Join(repoRoot, ".gitignore")
	contents, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to read .gitignore: %w", err)
	}
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == ignoreEntry {
			return nil
		}
	}
	if len(lines) > 0 && lines[len(lines)-1] != "" {
		lines = append(lines, "")
	}
	lines = append(lines, commentEntry, ignoreEntry)
	updated := strings.Join(lines, "\n")
	if err := os.WriteFile(path, []byte(updated), 0o644); err != nil {
		return fmt.Errorf("failed to update .gitignore: %w", err)
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
