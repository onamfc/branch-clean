package internal

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
)

var (
	// ErrCancelled is returned when the user cancels an operation
	ErrCancelled = errors.New("cancelled by user")
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorGray   = "\033[90m"
)

func PrintBranches(branches []Branch, mergedOnly, staleOnly bool) {
	fmt.Printf("\n%-30s %-10s %-10s %s\n", "Branch", "Status", "Age", "Last Commit")
	fmt.Println(strings.Repeat("-", 80))

	for _, b := range branches {
		if mergedOnly && !b.IsMerged {
			continue
		}
		if staleOnly && !b.IsStale {
			continue
		}

		status := getStatusString(b)
		age := getAgeString(b.LastCommit)
		date := b.LastCommit.Format("2006-01-02")

		if b.Protected {
			fmt.Printf("%s%-30s%s %s %s %s\n", colorGray, b.Name, colorReset, status, age, date)
		} else {
			fmt.Printf("%-30s %s %s %s\n", b.Name, status, age, date)
		}
	}
}

func getStatusString(b Branch) string {
	if b.IsMerged {
		return colorGreen + "merged" + colorReset + "   "
	}
	if b.IsStale {
		return colorYellow + "stale" + colorReset + "    "
	}
	return colorBlue + "active" + colorReset + "   "
}

func getAgeString(t time.Time) string {
	days := int(time.Since(t).Hours() / 24)
	if days == 0 {
		return "today     "
	}
	if days == 1 {
		return "1 day ago "
	}
	return fmt.Sprintf("%d days ago", days)
}

// FilterBranches filters branches based on merge and stale status.
// Protected branches are always excluded.
//
// Filtering logic:
// - If mergedOnly is true: only include merged branches
// - If staleOnly is true: only include stale branches
// - If both are false: include branches that are either merged OR stale (exclude active branches)
// - If both are true: include branches that are both merged AND stale
func FilterBranches(branches []Branch, mergedOnly, staleOnly bool) []Branch {
	var filtered []Branch
	for _, b := range branches {
		// Always skip protected branches
		if b.Protected {
			continue
		}

		// Apply filters based on flags
		if mergedOnly && staleOnly {
			// Both filters: must be merged AND stale
			if b.IsMerged && b.IsStale {
				filtered = append(filtered, b)
			}
		} else if mergedOnly {
			// Merged filter only: must be merged
			if b.IsMerged {
				filtered = append(filtered, b)
			}
		} else if staleOnly {
			// Stale filter only: must be stale
			if b.IsStale {
				filtered = append(filtered, b)
			}
		} else {
			// No specific filter: show branches that are either merged OR stale (exclude active)
			if b.IsMerged || b.IsStale {
				filtered = append(filtered, b)
			}
		}
	}
	return filtered
}

func SelectBranches(branches []Branch) ([]Branch, error) {
	if len(branches) == 0 {
		return nil, nil
	}

	selectedMap := make(map[int]bool)

	for {
		items := make([]string, len(branches)+1)
		for i, b := range branches {
			status := ""
			if b.IsMerged {
				status = "[merged]"
			} else if b.IsStale {
				status = "[stale]"
			}

			checkbox := "[ ]"
			if selectedMap[i] {
				checkbox = "[✓]"
			}

			items[i] = fmt.Sprintf("%s %s %s", checkbox, b.Name, status)
		}
		items[len(branches)] = colorGreen + "✓ Confirm selection" + colorReset

		prompt := promptui.Select{
			Label: "Select branches to delete (↑/↓ to navigate, enter to toggle/confirm)",
			Items: items,
			Size:  15,
			Templates: &promptui.SelectTemplates{
				Active:   "→ {{ . | cyan }}",
				Inactive: "  {{ . }}",
			},
		}

		idx, _, err := prompt.Run()
		if err != nil {
			// Handle interrupts (Ctrl+C)
			if err == promptui.ErrInterrupt || strings.Contains(err.Error(), "^C") {
				return nil, ErrCancelled
			}
			return nil, fmt.Errorf("selection failed: %w", err)
		}

		// If user selected the confirm option
		if idx == len(branches) {
			break
		}

		// Toggle selection
		selectedMap[idx] = !selectedMap[idx]
	}

	// Build the result slice
	var selected []Branch
	for i, branch := range branches {
		if selectedMap[i] {
			selected = append(selected, branch)
		}
	}

	return selected, nil
}

func ConfirmDeletion(branches []Branch, dryRun bool) bool {
	action := "delete"
	if dryRun {
		action = "would delete"
	}

	fmt.Printf("\n%sYou are about to %s %d branch(es):%s\n", colorRed, action, len(branches), colorReset)
	for _, b := range branches {
		fmt.Printf("  - %s\n", b.Name)
	}

	prompt := promptui.Prompt{
		Label:     "Continue",
		IsConfirm: true,
	}

	_, err := prompt.Run()
	return err == nil
}
