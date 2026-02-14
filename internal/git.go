package internal

import (
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// Custom error types for better error handling
var (
	ErrProtectedBranch = errors.New("protected branch")
	ErrCurrentBranch   = errors.New("cannot delete current branch")
	ErrDefaultBranch   = errors.New("cannot delete default branch")
)

// ProtectedBranchError represents an error when trying to delete a protected branch
type ProtectedBranchError struct {
	BranchName string
}

func (e *ProtectedBranchError) Error() string {
	return fmt.Sprintf("cannot delete protected branch '%s'", e.BranchName)
}

func (e *ProtectedBranchError) Is(target error) bool {
	return target == ErrProtectedBranch
}

type GitRepo struct {
	repo          *git.Repository
	repoPath      string
	defaultBranch string
}

type Branch struct {
	Name       string    `json:"name"`
	IsMerged   bool      `json:"is_merged"`
	IsStale    bool      `json:"is_stale"`
	LastCommit time.Time `json:"last_commit"`
	Protected  bool      `json:"protected"`
}

// NewGitRepo opens a git repository at the given path and detects the default branch.
// Returns an error if the path is not a valid git repository.
func NewGitRepo(path string) (*GitRepo, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open git repository at %s: %w\nIs this a git repository? Try running 'git status'", path, err)
	}

	defaultBranch, err := detectDefaultBranch(repo)
	if err != nil {
		return nil, err
	}

	return &GitRepo{
		repo:          repo,
		repoPath:      path,
		defaultBranch: defaultBranch,
	}, nil
}

func detectDefaultBranch(repo *git.Repository) (string, error) {
	// First, try to get the default branch from remote HEAD
	remote, err := repo.Remote("origin")
	if err == nil {
		refs, listErr := remote.List(&git.ListOptions{})
		if listErr == nil {
			for _, ref := range refs {
				if ref.Name().String() == "HEAD" {
					// This is a symbolic ref pointing to the default branch
					target := ref.Target()
					if target != "" {
						return target.Short(), nil
					}
				}
			}
		}
	}

	// Try common default branch names
	branches := []string{"main", "master", "develop"}
	for _, name := range branches {
		if _, refErr := repo.Reference(plumbing.NewBranchReferenceName(name), true); refErr == nil {
			return name, nil
		}
	}

	// If no common branches exist, list all branches and use the first one
	branchRefs, err := repo.Branches()
	if err != nil {
		return "", fmt.Errorf("could not determine default branch: %w", err)
	}

	var firstBranch string
	_ = branchRefs.ForEach(func(ref *plumbing.Reference) error {
		if firstBranch == "" {
			firstBranch = ref.Name().Short()
		}
		return nil
	})

	if firstBranch != "" {
		return firstBranch, nil
	}

	return "", fmt.Errorf("repository has no branches")
}

func (g *GitRepo) ListBranches(staleDays int, protectedPatterns []string) ([]Branch, error) {
	branchRefs, err := g.repo.Branches()
	if err != nil {
		return nil, err
	}

	var branches []Branch
	staleThreshold := time.Now().AddDate(0, 0, -staleDays)

	err = branchRefs.ForEach(func(ref *plumbing.Reference) error {
		name := ref.Name().Short()
		if name == g.defaultBranch {
			return nil
		}

		commit, commitErr := g.repo.CommitObject(ref.Hash())
		if commitErr != nil {
			return commitErr
		}

		isMerged, mergeErr := g.isMerged(name)
		if mergeErr != nil {
			return mergeErr
		}

		branch := Branch{
			Name:       name,
			IsMerged:   isMerged,
			IsStale:    commit.Committer.When.Before(staleThreshold),
			LastCommit: commit.Committer.When,
			Protected:  isProtected(name, protectedPatterns),
		}

		branches = append(branches, branch)
		return nil
	})

	return branches, err
}

// isMerged checks if a branch has been merged into the default branch using git CLI.
// This handles all merge strategies (merge commits, squash merges, rebase merges) correctly.
func (g *GitRepo) isMerged(branchName string) (bool, error) {
	// Use git merge-base --is-ancestor to check if the branch is merged
	// This works for all merge strategies
	cmd := exec.Command("git", "merge-base", "--is-ancestor", branchName, g.defaultBranch)
	cmd.Dir = g.repoPath

	err := cmd.Run()
	if err != nil {
		// Exit code 1 means not an ancestor (not merged)
		// Exit code 0 means it is an ancestor (merged)
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() == 1 {
				return false, nil
			}
		}
		return false, fmt.Errorf("failed to check merge status: %w", err)
	}

	return true, nil
}

// DeleteBranch deletes a branch by name.
// Returns ErrCurrentBranch if trying to delete the currently checked out branch.
// Returns ErrDefaultBranch if trying to delete the default branch.
func (g *GitRepo) DeleteBranch(name string) error {
	// Check if trying to delete current branch
	head, err := g.repo.Head()
	if err == nil && head.Name().Short() == name {
		return fmt.Errorf("%w: '%s'", ErrCurrentBranch, name)
	}

	// Check if trying to delete default branch
	if name == g.defaultBranch {
		return fmt.Errorf("%w: '%s'", ErrDefaultBranch, name)
	}

	return g.repo.Storer.RemoveReference(plumbing.NewBranchReferenceName(name))
}

// DeleteRemoteBranch deletes a branch from the remote repository.
func (g *GitRepo) DeleteRemoteBranch(name string) error {
	cmd := exec.Command("git", "push", "origin", "--delete", name)
	cmd.Dir = g.repoPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to delete remote branch: %w\nOutput: %s", err, output)
	}

	return nil
}

func isProtected(name string, patterns []string) bool {
	for _, pattern := range patterns {
		matched, err := filepath.Match(pattern, name)
		if err != nil {
			// Invalid pattern, try prefix matching as fallback
			if strings.Contains(pattern, "*") {
				prefix := strings.TrimSuffix(pattern, "*")
				if strings.HasPrefix(name, prefix) {
					return true
				}
			}
			continue
		}
		if matched {
			return true
		}
	}
	return false
}
