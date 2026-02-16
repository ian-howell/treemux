package treemux

import (
	"os"
	"path/filepath"
	"strings"
)

const (
	// rootNameOption stores the session root name.
	rootNameOption = "@tree_root_name"
	// rootDirOption stores the session root directory.
	rootDirOption = "@tree_root_dir"
)

// sanitizeSessionName normalizes names to be tmux-safe.
func sanitizeSessionName(name string) string {
	if isWhitespaceOnly(name) {
		return ""
	}
	for strings.HasPrefix(name, ".") {
		name = strings.TrimPrefix(name, ".")
	}
	replacer := strings.NewReplacer(
		".", "_",
		":", "_",
		"/", "_",
		" ", "_",
	)
	name = replacer.Replace(name)
	return name
}

// rootNameFromDir derives a session name from a directory path.
func rootNameFromDir(dir string) string {
	base := filepath.Base(dir)
	if base == "." || base == string(os.PathSeparator) || base == "" {
		base = "root"
	}
	return sanitizeSessionName(base)
}

// tmuxSeparator returns the configured session name separator.
func tmuxSeparator() string {
	if sep := os.Getenv("TMUX_TREE_SEPARATOR"); sep != "" {
		return sep
	}
	return " 🌿 "
}
