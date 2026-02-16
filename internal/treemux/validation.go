package treemux

import "strings"

func isWhitespaceOnly(value string) bool {
	return value != "" && strings.TrimSpace(value) == ""
}
