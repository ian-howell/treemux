package treemux

import (
	"fmt"
	"sort"
	"strings"
)

// Show returns formatted lines describing treemux sessions.
func (a *App) Show() ([]string, error) {
	sessions, err := a.tmux.ListSessions()
	if err != nil {
		return nil, err
	}
	roots := map[string]string{}
	children := map[string][]string{}

	separator := tmuxSeparator()
	for _, session := range sessions {
		rootName, err := a.tmux.ShowOption(session, rootNameOption)
		if err != nil || rootName == "" {
			continue
		}
		rootDir, err := a.tmux.ShowOption(session, rootDirOption)
		if err != nil || rootDir == "" {
			continue
		}
		if rootName == session {
			roots[rootName] = rootName
			continue
		}
		children[rootName] = append(children[rootName], session)
		if strings.Contains(session, separator) {
			rootPart := strings.SplitN(session, separator, 2)[0]
			if rootPart != "" {
				roots[rootPart] = rootPart
			}
		}
	}

	rootNames := make([]string, 0, len(roots))
	for name := range roots {
		rootNames = append(rootNames, name)
	}
	sort.Strings(rootNames)

	lines := make([]string, 0, len(rootNames))
	for _, root := range rootNames {
		lines = append(lines, root)
		childSessions := children[root]
		if len(childSessions) == 0 {
			continue
		}
		sort.Strings(childSessions)
		for _, child := range childSessions {
			childName := strings.TrimSpace(strings.TrimPrefix(child, root))
			childName = strings.TrimPrefix(childName, separator)
			childName = strings.TrimSpace(childName)
			if childName == "" {
				childName = child
			}
			lines = append(lines, fmt.Sprintf("%s  %s", root, childName))
		}
	}

	return lines, nil
}
