package treemux

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/ian-howell/treemux/internal/tmux"
)

const lastAttachedOption = "@tree_last_attached_time"

type sessionInfo struct {
	Name       string
	RootName   string
	RootDir    string
	LastUsedNs int64
}

// ShowRoots returns formatted root session lines.
func (a *App) ShowRoots(req ShowRootsRequest) ([]string, error) {
	if err := validateSortBy(req.SortBy); err != nil {
		return nil, err
	}
	roots, _, currentRoot, err := a.collectSessions()
	if err != nil {
		return nil, err
	}
	entries := make([]sessionInfo, 0, len(roots))
	for _, info := range roots {
		if req.HideCurrent && currentRoot != "" && info.Name == currentRoot {
			continue
		}
		entries = append(entries, info)
	}
	entries = sortSessions(entries, req.SortBy)
	lines := make([]string, 0, len(entries))
	for _, info := range entries {
		prefix := " "
		if info.Name == currentRoot && currentRoot != "" {
			prefix = "*"
		}
		lines = append(lines, fmt.Sprintf("%s %s", prefix, info.Name))
	}
	return lines, nil
}

// ShowChildren returns formatted child session lines for a root.
func (a *App) ShowChildren(req ShowChildrenRequest) ([]string, error) {
	if err := validateSortBy(req.SortBy); err != nil {
		return nil, err
	}
	rootName, err := a.resolveChildRoot(req.Root)
	if err != nil {
		return nil, err
	}
	roots, children, _, err := a.collectSessions()
	if err != nil {
		return nil, err
	}
	if _, ok := roots[rootName]; !ok {
		return nil, fmt.Errorf("root session does not exist: %s", rootName)
	}
	currentSession := ""
	if a.isInsideTmux() {
		currentSession, _ = a.currentSessionName()
		currentSession = strings.TrimSpace(currentSession)
	}
	childEntries := make([]sessionInfo, 0, len(children[rootName]))
	for _, info := range children[rootName] {
		if req.HideCurrent && currentSession != "" && info.Name == currentSession {
			continue
		}
		childEntries = append(childEntries, info)
	}
	childEntries = sortSessions(childEntries, req.SortBy)
	lines := make([]string, 0, len(childEntries))
	for _, info := range childEntries {
		prefix := " "
		if info.Name == currentSession && currentSession != "" {
			prefix = "*"
		}
		label := childLabel(rootName, info.Name)
		lines = append(lines, fmt.Sprintf("%s %s", prefix, label))
	}
	return lines, nil
}

func (a *App) collectSessions() (map[string]sessionInfo, map[string][]sessionInfo, string, error) {
	sessions, err := a.tmux.ListSessions()
	if err != nil {
		return nil, nil, "", err
	}
	separator := tmuxSeparator()
	roots := map[string]sessionInfo{}
	children := map[string][]sessionInfo{}
	for _, session := range sessions {
		rootName, err := a.tmux.ShowOption(session, rootNameOption)
		if err != nil || rootName == "" {
			continue
		}
		rootDir, err := a.tmux.ShowOption(session, rootDirOption)
		if err != nil || rootDir == "" {
			continue
		}
		lastUsed := readLastUsed(a.tmux, session)
		info := sessionInfo{
			Name:       session,
			RootName:   rootName,
			RootDir:    rootDir,
			LastUsedNs: lastUsed,
		}
		if rootName == session {
			roots[rootName] = info
			continue
		}
		children[rootName] = append(children[rootName], info)
		if strings.Contains(session, separator) {
			rootPart := strings.SplitN(session, separator, 2)[0]
			if rootPart != "" {
				if _, ok := roots[rootPart]; !ok {
					roots[rootPart] = sessionInfo{Name: rootPart, RootName: rootPart, RootDir: rootDir}
				}
			}
		}
	}
	currentRoot := ""
	if a.isInsideTmux() {
		currentSession, err := a.currentSessionName()
		if err == nil {
			currentSession = strings.TrimSpace(currentSession)
			if currentSession != "" {
				currentRoot, _ = a.tmux.ShowOption(currentSession, rootNameOption)
				currentRoot = strings.TrimSpace(currentRoot)
			}
		}
	}
	return roots, children, currentRoot, nil
}

func sortSessions(entries []sessionInfo, sortBy string) []sessionInfo {
	mode := strings.TrimSpace(strings.ToLower(sortBy))
	switch mode {
	case "most-recently-used", "mru":
		sort.Slice(entries, func(i, j int) bool {
			if entries[i].LastUsedNs == entries[j].LastUsedNs {
				return entries[i].Name < entries[j].Name
			}
			return entries[i].LastUsedNs > entries[j].LastUsedNs
		})
	default:
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Name < entries[j].Name
		})
	}
	return entries
}

// validateSortBy ensures the sort mode is recognized.
func validateSortBy(sortBy string) error {
	mode := strings.TrimSpace(strings.ToLower(sortBy))
	switch mode {
	case "", "alphabetic", "most-recently-used", "mru":
		return nil
	default:
		return fmt.Errorf("unknown sort mode: %s", sortBy)
	}
}

func childLabel(rootName, sessionName string) string {
	separator := tmuxSeparator()
	childName := strings.TrimPrefix(sessionName, rootName+separator)
	childName = strings.TrimSpace(childName)
	if childName == "" {
		return sessionName
	}
	return childName
}

func readLastUsed(client *tmux.Client, session string) int64 {
	value, err := client.ShowOption(session, lastAttachedOption)
	if err != nil {
		return 0
	}
	value = strings.TrimSpace(value)
	if value == "" {
		return 0
	}
	if parsed, err := strconv.ParseInt(value, 10, 64); err == nil {
		return parsed
	}
	return 0
}

func setLastAttached(client *tmux.Client, session string) error {
	return client.SetOption(session, lastAttachedOption, fmt.Sprintf("%d", time.Now().UTC().UnixNano()))
}
