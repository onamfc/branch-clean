package internal

import (
	"strings"
	"testing"
	"time"
)

func TestGetStatusString(t *testing.T) {
	tests := []struct {
		name   string
		branch Branch
		want   string
	}{
		{"merged", Branch{IsMerged: true}, "merged"},
		{"stale", Branch{IsStale: true}, "stale"},
		{"active", Branch{}, "active"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getStatusString(tt.branch)
			if !strings.Contains(got, tt.want) {
				t.Errorf("getStatusString() = %v, want to contain %v", got, tt.want)
			}
		})
	}
}

func TestGetAgeString(t *testing.T) {
	tests := []struct {
		name string
		time time.Time
		want string
	}{
		{"today", time.Now(), "today"},
		{"1 day", time.Now().AddDate(0, 0, -1), "1 day ago"},
		{"30 days", time.Now().AddDate(0, 0, -30), "30 days ago"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getAgeString(tt.time)
			if !strings.Contains(got, tt.want) {
				t.Errorf("getAgeString() = %v, want to contain %v", got, tt.want)
			}
		})
	}
}

func TestFilterBranches(t *testing.T) {
	branches := []Branch{
		{Name: "merged-branch", IsMerged: true, Protected: false},
		{Name: "stale-branch", IsStale: true, Protected: false},
		{Name: "protected-branch", IsMerged: true, Protected: true},
		{Name: "active-branch", Protected: false},
	}

	t.Run("merged only", func(t *testing.T) {
		filtered := FilterBranches(branches, true, false)
		if len(filtered) != 1 || filtered[0].Name != "merged-branch" {
			t.Errorf("expected 1 merged branch, got %d", len(filtered))
		}
	})

	t.Run("stale only", func(t *testing.T) {
		filtered := FilterBranches(branches, false, true)
		if len(filtered) != 1 || filtered[0].Name != "stale-branch" {
			t.Errorf("expected 1 stale branch, got %d", len(filtered))
		}
	})

	t.Run("no protected", func(t *testing.T) {
		filtered := FilterBranches(branches, true, false)
		for _, b := range filtered {
			if b.Protected {
				t.Error("protected branch in filtered results")
			}
		}
	})

	t.Run("all filters", func(t *testing.T) {
		filtered := FilterBranches(branches, false, false)
		if len(filtered) != 2 {
			t.Errorf("expected 2 branches (merged+stale), got %d", len(filtered))
		}
	})
}

func TestFilterBranches_EmptyInput(t *testing.T) {
	filtered := FilterBranches([]Branch{}, false, false)
	if len(filtered) != 0 {
		t.Error("expected empty result for empty input")
	}
}

func TestFilterBranches_BothFilters(t *testing.T) {
	branches := []Branch{
		{Name: "merged-only", IsMerged: true, IsStale: false, Protected: false},
		{Name: "stale-only", IsMerged: false, IsStale: true, Protected: false},
		{Name: "both", IsMerged: true, IsStale: true, Protected: false},
		{Name: "neither", IsMerged: false, IsStale: false, Protected: false},
	}

	// Both filters: should only get branches that are BOTH merged AND stale
	filtered := FilterBranches(branches, true, true)
	if len(filtered) != 1 || filtered[0].Name != "both" {
		t.Errorf("expected 1 branch (merged AND stale), got %d", len(filtered))
	}
}

func TestFilterBranches_NoFilters(t *testing.T) {
	branches := []Branch{
		{Name: "merged", IsMerged: true, Protected: false},
		{Name: "stale", IsStale: true, Protected: false},
		{Name: "active", IsMerged: false, IsStale: false, Protected: false},
	}

	// No filters: should get branches that are merged OR stale (but not active)
	filtered := FilterBranches(branches, false, false)
	if len(filtered) != 2 {
		t.Errorf("expected 2 branches (merged OR stale), got %d", len(filtered))
	}

	// Verify active branch is excluded
	for _, b := range filtered {
		if b.Name == "active" {
			t.Error("active branch should not be in filtered results")
		}
	}
}

func TestFilterBranches_AllProtected(t *testing.T) {
	branches := []Branch{
		{Name: "main", IsMerged: true, Protected: true},
		{Name: "master", IsMerged: true, Protected: true},
	}

	filtered := FilterBranches(branches, false, false)
	if len(filtered) != 0 {
		t.Error("expected no branches when all are protected")
	}
}

func TestSelectBranches_EmptyList(t *testing.T) {
	selected, err := SelectBranches([]Branch{})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if selected != nil {
		t.Error("expected nil result for empty input")
	}
}
