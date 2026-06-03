package cmd

import (
	"fmt"
	"sort"
	"time"

	"github.com/DucChau/driftmap/internal/store"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show <snapshot-id>",
	Short: "Show all variables in a snapshot",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := store.Open()
		if err != nil {
			return err
		}

		snap, err := s.Resolve(args[0])
		if err != nil {
			return err
		}

		bold := color.New(color.Bold).SprintFunc()
		cyan := color.New(color.FgCyan).SprintFunc()
		dim := color.New(color.Faint).SprintFunc()

		fmt.Printf("%s\n", bold("Snapshot: ")+cyan(snap.ShortID()))
		if snap.Label != "" {
			fmt.Printf("%s %s\n", bold("Label:"), snap.Label)
		}
		fmt.Printf("%s %s\n", bold("Captured:"), snap.CapturedAt.Format(time.RFC3339))
		fmt.Printf("%s %d\n\n", bold("Variables:"), len(snap.Vars))

		keys := make([]string, 0, len(snap.Vars))
		for k := range snap.Vars {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			v := snap.Vars[k]
			if len(v) > 80 {
				v = v[:77] + dim("...")
			}
			fmt.Printf("  %s=%s\n", bold(k), v)
		}

		return nil
	},
}
