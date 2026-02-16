package treemux

import "testing"

func TestSortSessionsMostRecentlyUsed(t *testing.T) {
	entries := []sessionInfo{
		{Name: "bar", LastUsedNs: 2},
		{Name: "foo", LastUsedNs: 3},
		{Name: "baz", LastUsedNs: 1},
	}
	sorted := sortSessions(entries, "most-recently-used")
	if sorted[0].Name != "foo" {
		t.Fatalf("expected most-recently-used first, got %s", sorted[0].Name)
	}
}

func TestSortSessionsAlphabetic(t *testing.T) {
	entries := []sessionInfo{
		{Name: "bar", LastUsedNs: 2},
		{Name: "foo", LastUsedNs: 3},
		{Name: "baz", LastUsedNs: 1},
	}
	sorted := sortSessions(entries, "alphabetic")
	if sorted[0].Name != "bar" {
		t.Fatalf("expected alphabetic first, got %s", sorted[0].Name)
	}
}

func TestShowChildrenOutputShape(t *testing.T) {
	root := "root"
	sep := " 🌿 "
	entries := []sessionInfo{
		{Name: root + sep + "alpha"},
		{Name: root + sep + "beta"},
	}
	childLabels := make([]string, 0, len(entries))
	for _, entry := range entries {
		childLabels = append(childLabels, childLabel(root, entry.Name))
	}
	if childLabels[0] != "alpha" || childLabels[1] != "beta" {
		t.Fatalf("expected child labels without separator, got %v", childLabels)
	}
}

func TestChildLabelStripsSeparator(t *testing.T) {
	root := "treemux"
	sep := " 🌿 "
	name := root + sep + "nvim"
	label := childLabel(root, name)
	if label != "nvim" {
		t.Fatalf("expected child label without separator, got %s", label)
	}
}

func TestValidateSortByRejectsUnknown(t *testing.T) {
	if err := validateSortBy("newest"); err == nil {
		t.Fatalf("expected error for unknown sort mode")
	}
}
