package cmd

import (
	"fmt"
	"time"

	"github.com/DucChau/driftmap/internal/store"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all captured snapshots",
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := store.Open()
		if err != nil {
			return err
		}

		snaps, err := s.List()
		if err != nil {
			return err
		}

		if len(snaps) == 0 {
			fmt.Println(color.YellowString("No snapshots found. Run 'driftmap capture' first."))
			return nil
		}

		cyan := color.New(color.FgCyan).SprintFunc()
		bold := color.New(color.Bold).SprintFunc()
		dim := color.New(color.Faint).SprintFunc()

		fmt.Printf("%-12s  %-20s  %-8s  %s\n",
			bold("ID"), bold("Captured At"), bold("Vars"), bold("Label"))
		fmt.Println("─────────────────────────────────────────────────────────")

		for _, snap := range snaps {
			label := snap.Label
			if label == "" {
				label = dim("(unlabeled)")
			}
			fmt.Printf("%-12s  %-20s  %-8d  %s\n",
				cyan(snap.ShortID()),
				snap.CapturedAt.Format(time.RFC3339),
				len(snap.Vars),
				label,
			)
		}

		return nil
	},
}
