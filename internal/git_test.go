package internal

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func setupTestRepo(t *testing.T) (string, *git.Repository) {
	tmpDir := t.TempDir()
	repo, err := git.PlainInit(tmpDir, false)
	if err != nil {
		t.Fatalf("failed to init repo: %v", err)
	}

	w, _ := repo.Worktree()
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}
	w.Add("test.txt")
	w.Commit("initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test",
			Email: "test@test.com",
			When:  time.Now(),
		},
	})

	return tmpDir, repo
}

func TestNewGitRepo(t *testing.T) {
	tmpDir, _ := setupTestRepo(t)

	gitRepo, err := NewGitRepo(tmpDir)
	if err != nil {
		t.Fatalf("NewGitRepo failed: %v", err)
	}

	if gitRepo.repo == nil {
		t.Error("repo is nil")
	}
	if gitRepo.defaultBranch == "" {
		t.Error("defaultBranch is empty")
	}
}

func TestNewGitRepo_InvalidPath(t *testing.T) {
	_, err := NewGitRepo("/nonexistent/path")
	if err == nil {
		t.Error("expected error for invalid path")
	}
}

func TestListBranches(t *testing.T) {
	tmpDir, repo := setupTestRepo(t)

	head, _ := repo.Head()
	repo.Storer.SetReference(plumbing.NewHashReference(
		plumbing.NewBranchReferenceName("feature-branch"),
		head.Hash(),
	))

	gitRepo, _ := NewGitRepo(tmpDir)
	branches, err := gitRepo.ListBranches(30, []string{"main", "master"})

	if err != nil {
		t.Fatalf("ListBranches failed: %v", err)
	}

	if len(branches) == 0 {
		t.Error("expected at least one branch")
	}
}

func TestIsProtected(t *testing.T) {
	tests := []struct {
		name     string
		branch   string
		patterns []string
		want     bool
	}{
		{"exact match", "main", []string{"main", "master"}, true},
		{"wildcard match", "release/v1.0", []string{"release/*"}, true},
		{"no match", "feature-x", []string{"main", "master"}, false},
		{"empty patterns", "main", []string{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isProtected(tt.branch, tt.patterns)
			if got != tt.want {
				t.Errorf("isProtected() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeleteBranch(t *testing.T) {
	tmpDir, repo := setupTestRepo(t)

	head, _ := repo.Head()
	repo.Storer.SetReference(plumbing.NewHashReference(
		plumbing.NewBranchReferenceName("test-branch"),
		head.Hash(),
	))

	gitRepo, _ := NewGitRepo(tmpDir)
	err := gitRepo.DeleteBranch("test-branch")

	if err != nil {
		t.Errorf("DeleteBranch failed: %v", err)
	}

	_, err = repo.Reference(plumbing.NewBranchReferenceName("test-branch"), true)
	if err == nil {
		t.Error("branch still exists after deletion")
	}
}

func TestDeleteBranch_DefaultBranch(t *testing.T) {
	tmpDir, _ := setupTestRepo(t)
	gitRepo, _ := NewGitRepo(tmpDir)

	err := gitRepo.DeleteBranch(gitRepo.defaultBranch)
	if err == nil {
		t.Error("expected error when trying to delete default branch")
	}
}

func TestDeleteBranch_NonExistent(t *testing.T) {
	tmpDir, _ := setupTestRepo(t)
	gitRepo, _ := NewGitRepo(tmpDir)

	err := gitRepo.DeleteBranch("non-existent-branch")
	// Should not panic, but may return an error
	_ = err
}

func TestIsProtected_InvalidPattern(t *testing.T) {
	// Test with invalid glob pattern
	result := isProtected("test-branch", []string{"[invalid"})
	// Should not panic, and will use fallback logic
	if result {
		t.Error("expected false for branch with invalid pattern")
	}
}

func TestIsProtected_WildcardPrefix(t *testing.T) {
	tests := []struct {
		name     string
		branch   string
		patterns []string
		want     bool
	}{
		{"prefix match with wildcard", "release/v1.0", []string{"release/*"}, true},
		{"prefix match with wildcard", "hotfix/bug-123", []string{"hotfix/*"}, true},
		{"no prefix match", "feature/x", []string{"release/*"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isProtected(tt.branch, tt.patterns)
			if got != tt.want {
				t.Errorf("isProtected(%q, %v) = %v, want %v", tt.branch, tt.patterns, got, tt.want)
			}
		})
	}
}
