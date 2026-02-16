package treemux

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ian-howell/treemux/internal/git"
)

func TestWorktreeDefaultPath(t *testing.T) {
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

func TestEnsureBranchCreatesLocalBranch(t *testing.T) {
	repo := newTestRepo(t)
	branch := "example"
	if err := git.EnsureBranch(repo, branch); err != nil {
		t.Fatalf("expected branch creation, got %v", err)
	}
	if ok, err := hasLocalBranch(repo, branch); err != nil {
		t.Fatalf("expected branch check to succeed, got %v", err)
	} else if !ok {
		t.Fatalf("expected local branch to exist")
	}
}

func TestEnsureBranchTracksRemoteBranch(t *testing.T) {
	repo := newTestRepo(t)
	branch := "remote-branch"
	if err := createRemoteBranch(t, repo, branch); err != nil {
		t.Fatalf("failed to create remote branch: %v", err)
	}
	if err := git.EnsureBranch(repo, branch); err != nil {
		t.Fatalf("expected branch tracking, got %v", err)
	}
	if ok, err := hasLocalBranch(repo, branch); err != nil {
		t.Fatalf("expected branch check to succeed, got %v", err)
	} else if !ok {
		t.Fatalf("expected local branch to exist")
	}
}

func newTestRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if err := runGit(dir, "init"); err != nil {
		t.Fatalf("git init failed: %v", err)
	}
	if err := runGit(dir, "remote", "remove", "origin"); err != nil {
		// ignore if origin doesn't exist
	}
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("test"), 0o644); err != nil {
		t.Fatalf("write file failed: %v", err)
	}
	if err := runGit(dir, "add", "."); err != nil {
		t.Fatalf("git add failed: %v", err)
	}
	if err := runGit(dir, "commit", "-m", "init"); err != nil {
		t.Fatalf("git commit failed: %v", err)
	}
	return dir
}

func createRemoteBranch(t *testing.T, repo, branch string) error {
	t.Helper()
	remoteDir := t.TempDir()
	if err := runGit(remoteDir, "init", "--bare"); err != nil {
		return err
	}
	if err := runGit(repo, "remote", "add", "origin", remoteDir); err != nil {
		return err
	}
	if err := runGit(repo, "checkout", "-b", branch); err != nil {
		return err
	}
	if err := runGit(repo, "push", "-u", "origin", branch); err != nil {
		return err
	}
	if err := runGit(repo, "checkout", "main"); err != nil {
		return err
	}
	if err := runGit(repo, "branch", "-D", branch); err != nil {
		return err
	}
	if err := runGit(repo, "fetch", "--all", "--prune"); err != nil {
		return err
	}
	return nil
}

func hasLocalBranch(repoRoot, branch string) (bool, error) {
	cmd := exec.Command("git", "show-ref", "--verify", "refs/heads/"+branch)
	cmd.Dir = repoRoot
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			stderr := strings.TrimSpace(string(exitErr.Stderr))
			if strings.Contains(stderr, "not a valid ref") || strings.Contains(stderr, "fatal") {
				return false, nil
			}
		}
		return false, err
	}
	return true, nil
}

func runGit(dir string, args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=treemux",
		"GIT_AUTHOR_EMAIL=treemux@example.com",
		"GIT_COMMITTER_NAME=treemux",
		"GIT_COMMITTER_EMAIL=treemux@example.com",
	)
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			stderr := strings.TrimSpace(string(exitErr.Stderr))
			if stderr != "" {
				return fmt.Errorf("git %s: %s", strings.Join(args, " "), stderr)
			}
		}
		return err
	}
	return nil
}
