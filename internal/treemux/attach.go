package treemux

import (
	"fmt"
	"os"
	"strings"
)

// AttachRoot ensures and attaches to a root session.
func (a *App) AttachRoot(req AttachRootRequest) error {
	if req.Dir != "" && req.Worktree != "" {
		return fmt.Errorf("--dir and --worktree are mutually exclusive")
	}
	rootDir, err := a.resolveRootDir(req.Dir, req.Worktree)
	if err != nil {
		return err
	}
	if req.Worktree != "" {
		if _, err := resolveWorktree(req.Worktree, req.Dir); err != nil {
			return err
		}
	}
	rootDir, err = normalizePath(rootDir)
	if err != nil {
		return err
	}
	if info, err := os.Stat(rootDir); err != nil || !info.IsDir() {
		return fmt.Errorf("root directory does not exist: %s", rootDir)
	}
	rootName := sanitizeSessionName(req.Name)
	if rootName == "" {
		rootName = rootNameFromDir(rootDir)
	}
	if rootName == "" {
		return fmt.Errorf("root session name is empty")
	}

	exists, err := a.tmux.HasSession(rootName)
	if err != nil {
		return err
	}
	if !exists {
		if err := a.tmux.NewSession(rootName, rootDir, nil); err != nil {
			return fmt.Errorf("failed to create tmux session '%s': %w", rootName, err)
		}
		if err := a.tmux.SetOption(rootName, rootNameOption, rootName); err != nil {
			return fmt.Errorf("failed to set tmux option %s for session '%s': %w", rootNameOption, rootName, err)
		}
		if err := a.tmux.SetOption(rootName, rootDirOption, rootDir); err != nil {
			return fmt.Errorf("failed to set tmux option %s for session '%s': %w", rootDirOption, rootName, err)
		}
	}

	if err := a.attach(rootName); err != nil {
		return err
	}
	return setLastAttached(a.tmux, rootName)
}

// AttachChild ensures and attaches to a child session.
func (a *App) AttachChild(req AttachChildRequest) error {
	rootName, err := a.resolveChildRoot(req.Root)
	if err != nil {
		return err
	}
	exists, err := a.tmux.HasSession(rootName)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("root session does not exist: %s", rootName)
	}

	rootDir, err := a.tmux.ShowOption(rootName, rootDirOption)
	if err != nil {
		return err
	}
	if rootDir == "" {
		return fmt.Errorf("missing root metadata (%s) for session: %s", rootDirOption, rootName)
	}

	// TODO: Simplify this. I'm not convinced that resolceWorktree needs a dir argument
	childDir := rootDir
	childName := sanitizeSessionName(req.Name)
	if childName == "" {
		return fmt.Errorf("child session name is empty")
	}
	childSession := rootName + tmuxSeparator() + childName

	exists, err = a.tmux.HasSession(childSession)
	if err != nil {
		return err
	}
	if !exists {
		if err := a.tmux.NewSession(childSession, childDir, commandArgs(req.Command)); err != nil {
			return fmt.Errorf("failed to create child tmux session '%s': %w", childSession, err)
		}
		if err := a.tmux.SetOption(childSession, rootNameOption, rootName); err != nil {
			return fmt.Errorf("failed to set tmux option %s for child session '%s': %w", rootNameOption, childSession, err)
		}
		if err := a.tmux.SetOption(childSession, rootDirOption, rootDir); err != nil {
			return fmt.Errorf("failed to set tmux option %s for child session '%s': %w", rootDirOption, childSession, err)
		}
	}

	if err := a.attach(childSession); err != nil {
		return err
	}
	if err := setLastAttached(a.tmux, childSession); err != nil {
		return err
	}
	return setLastAttached(a.tmux, rootName)
}

// resolveChildRoot determines the root session name for a child attach.
func (a *App) resolveChildRoot(root string) (string, error) {
	rootName := sanitizeSessionName(root)
	if rootName != "" {
		return rootName, nil
	}
	if !a.tmux.IsInsideTmux() {
		return "", fmt.Errorf("--root is required outside tmux")
	}
	current, err := a.tmux.CurrentSessionName()
	if err != nil {
		return "", err
	}
	current = strings.TrimSpace(current)
	if current == "" {
		return "", fmt.Errorf("failed to determine current tmux session")
	}
	rootName, err = a.tmux.ShowOption(current, rootNameOption)
	if err != nil {
		return "", err
	}
	rootName = strings.TrimSpace(rootName)
	if rootName == "" {
		return "", fmt.Errorf("current session is not a treemux session")
	}
	return sanitizeSessionName(rootName), nil
}

// resolveRootDir determines the root directory for the session.

func (a *App) resolveRootDir(dir, worktree string) (string, error) {
	if dir != "" {
		return dir, nil
	}
	if worktree != "" {
		return worktreeDefaultPath(worktree)
	}
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}
	return cwd, nil
}

// commandArgs builds a command invocation for tmux.
func commandArgs(command string) []string {
	command = strings.TrimSpace(command)
	if command == "" {
		return nil
	}
	return []string{"sh", "-lc", command}
}
