package cmd

import (
	"fmt"
	"sort"

	"github.com/DucChau/driftmap/internal/differ"
	"github.com/DucChau/driftmap/internal/store"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff <snapshot-a> <snapshot-b>",
	Short: "Diff two snapshots and show what changed",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := store.Open()
		if err != nil {
			return err
		}

		a, err := s.Resolve(args[0])
		if err != nil {
			return fmt.Errorf("resolving %q: %w", args[0], err)
		}
		b, err := s.Resolve(args[1])
		if err != nil {
			return fmt.Errorf("resolving %q: %w", args[1], err)
		}

		result := differ.Diff(a, b)

		if len(result.Added)+len(result.Removed)+len(result.Changed) == 0 {
			fmt.Println(color.GreenString("✓ No differences found — environments are identical."))
			return nil
		}

		red := color.New(color.FgRed).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()
		bold := color.New(color.Bold).SprintFunc()

		fmt.Printf("\nDiff: %s → %s\n\n", bold(a.ShortID()), bold(b.ShortID()))

		if len(result.Added) > 0 {
			fmt.Printf("%s (%d added)\n", green("+ ADDED"), len(result.Added))
			keys := sortedKeys(result.Added)
			for _, k := range keys {
				fmt.Printf("  %s %s=%s\n", green("+"), bold(k), result.Added[k])
			}
			fmt.Println()
		}

		if len(result.Removed) > 0 {
			fmt.Printf("%s (%d removed)\n", red("- REMOVED"), len(result.Removed))
			keys := sortedKeys(result.Removed)
			for _, k := range keys {
				fmt.Printf("  %s %s=%s\n", red("-"), bold(k), result.Removed[k])
			}
			fmt.Println()
		}

		if len(result.Changed) > 0 {
			fmt.Printf("%s (%d changed)\n", yellow("~ CHANGED"), len(result.Changed))
			changedKeys := make([]string, 0, len(result.Changed))
			for k := range result.Changed {
				changedKeys = append(changedKeys, k)
			}
			sort.Strings(changedKeys)
			for _, k := range changedKeys {
				c := result.Changed[k]
				fmt.Printf("  %s %s\n", yellow("~"), bold(k))
				fmt.Printf("    %s %s\n", red("-"), c.Before)
				fmt.Printf("    %s %s\n", green("+"), c.After)
			}
		}

		return nil
	},
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
