package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/onamfc/branch-clean/internal"
	"github.com/spf13/cobra"
)

var (
	dryRun       bool
	staleDays    int
	protected    []string
	mergedOnly   bool
	staleOnly    bool
	verbose      bool
	force        bool
	assumeYes    bool
	deleteRemote bool
	outputFormat string
	version      = "dev" // Set via ldflags at build time
)

var rootCmd = &cobra.Command{
	Use:   "branch-clean",
	Short: "Safely delete merged and stale git branches",
	Long:  "Interactive tool to clean up merged and stale git branches with safety checks",
	RunE:  runCleanup,
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all branches with their status",
	RunE:  runList,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("branch-clean version %s\n", version)
	},
}

func init() {
	// Load configuration from file
	config, err := internal.LoadConfig()
	if err != nil {
		// Config load failed, use defaults
		config = internal.DefaultConfig()
	}

	rootCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "d", false, "Show what would be deleted without making changes")
	rootCmd.PersistentFlags().IntVarP(&staleDays, "stale-days", "s", config.StaleDays, "Days since last commit to consider branch stale")
	rootCmd.PersistentFlags().StringSliceVarP(&protected, "protect", "p", config.Protected, "Protected branch patterns")
	rootCmd.PersistentFlags().BoolVarP(&mergedOnly, "merged-only", "m", false, "Only show merged branches")
	rootCmd.PersistentFlags().BoolVar(&staleOnly, "stale-only", false, "Only show stale branches")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	rootCmd.PersistentFlags().BoolVarP(&force, "force", "f", false, "Skip confirmation prompt")
	rootCmd.PersistentFlags().BoolVarP(&assumeYes, "yes", "y", false, "Automatically answer yes to all prompts")
	rootCmd.PersistentFlags().BoolVar(&deleteRemote, "remote", false, "Also delete branches from remote")

	listCmd.Flags().StringVar(&outputFormat, "format", "table", "Output format: table or json")

	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(versionCmd)
}

func validateGitRepo(path string) error {
	gitDir := filepath.Join(path, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		return fmt.Errorf("not a git repository (or any parent up to mount point)\nTry running this command from within a git repository")
	}
	return nil
}

func validateFlags() error {
	if staleDays <= 0 {
		return fmt.Errorf("stale-days must be positive, got %d", staleDays)
	}
	return nil
}

func runList(cmd *cobra.Command, args []string) error {
	// Validate flags
	if err := validateFlags(); err != nil {
		return err
	}

	// Validate output format
	if outputFormat != "table" && outputFormat != "json" {
		return fmt.Errorf("invalid output format: %s (must be 'table' or 'json')", outputFormat)
	}

	repoPath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	// Validate git repository
	if err := validateGitRepo(repoPath); err != nil {
		return err
	}

	git, err := internal.NewGitRepo(repoPath)
	if err != nil {
		return err
	}

	branches, err := git.ListBranches(staleDays, protected)
	if err != nil {
		return fmt.Errorf("failed to list branches: %w", err)
	}

	// Filter branches if needed
	filtered := branches
	if mergedOnly || staleOnly {
		filtered = internal.FilterBranches(branches, mergedOnly, staleOnly)
	}

	// Output based on format
	if outputFormat == "json" {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(filtered); err != nil {
			return fmt.Errorf("failed to encode JSON: %w", err)
		}
	} else {
		internal.PrintBranches(filtered, false, false)
	}

	return nil
}

func runCleanup(cmd *cobra.Command, args []string) error {
	// Validate flags
	if err := validateFlags(); err != nil {
		return err
	}

	repoPath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	// Validate git repository
	if err := validateGitRepo(repoPath); err != nil {
		return err
	}

	git, err := internal.NewGitRepo(repoPath)
	if err != nil {
		return err
	}

	branches, err := git.ListBranches(staleDays, protected)
	if err != nil {
		return fmt.Errorf("failed to list branches: %w", err)
	}

	filtered := internal.FilterBranches(branches, mergedOnly, staleOnly)
	if len(filtered) == 0 {
		fmt.Println("No branches to clean up")
		return nil
	}

	selected, err := internal.SelectBranches(filtered)
	if err != nil {
		return fmt.Errorf("branch selection failed: %w", err)
	}

	if len(selected) == 0 {
		fmt.Println("No branches selected")
		return nil
	}

	// Skip confirmation if force or assumeYes flag is set
	if !force && !assumeYes && !internal.ConfirmDeletion(selected, dryRun) {
		fmt.Println("Cancelled")
		return nil
	}

	if dryRun {
		fmt.Println("\n[DRY RUN] Would delete:")
		for _, b := range selected {
			fmt.Printf("  - %s", b.Name)
			if deleteRemote {
				fmt.Printf(" (local and remote)")
			}
			fmt.Println()
		}
		return nil
	}

	var hasErrors bool
	var successCount int
	for _, branch := range selected {
		if verbose {
			fmt.Printf("Deleting branch: %s\n", branch.Name)
		}

		// Delete local branch
		if err := git.DeleteBranch(branch.Name); err != nil {
			fmt.Fprintf(os.Stderr, "✗ Failed to delete local branch %s: %v\n", branch.Name, err)
			hasErrors = true
			continue
		}
		fmt.Printf("✓ Deleted local branch %s\n", branch.Name)
		successCount++

		// Delete remote branch if flag is set
		if deleteRemote {
			if err := git.DeleteRemoteBranch(branch.Name); err != nil {
				fmt.Fprintf(os.Stderr, "⚠ Failed to delete remote branch %s: %v\n", branch.Name, err)
				// Don't mark as error since local deletion succeeded
			} else {
				fmt.Printf("✓ Deleted remote branch %s\n", branch.Name)
			}
		}
	}

	fmt.Printf("\nDeleted %d of %d branches\n", successCount, len(selected))

	if hasErrors {
		return fmt.Errorf("some branches failed to delete")
	}
	return nil
}

const (
	exitSuccess         = 0
	exitError           = 1
	exitProtectedBranch = 2
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		// Check for custom error types
		var protectedErr *internal.ProtectedBranchError
		if errors.As(err, &protectedErr) {
			os.Exit(exitProtectedBranch)
		}
		if errors.Is(err, internal.ErrProtectedBranch) ||
			errors.Is(err, internal.ErrCurrentBranch) ||
			errors.Is(err, internal.ErrDefaultBranch) {
			os.Exit(exitProtectedBranch)
		}
		os.Exit(exitError)
	}
}
