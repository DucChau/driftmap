package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/DucChau/driftmap/internal/snapshot"
	"github.com/DucChau/driftmap/internal/store"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var captureCmd = &cobra.Command{
	Use:   "capture [label]",
	Short: "Capture the current environment as a named snapshot",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		label := ""
		if len(args) > 0 {
			label = args[0]
		}

		env := os.Environ()
		snap := snapshot.New(label, env)

		s, err := store.Open()
		if err != nil {
			return fmt.Errorf("opening store: %w", err)
		}

		if err := s.Save(snap); err != nil {
			return fmt.Errorf("saving snapshot: %w", err)
		}

		green := color.New(color.FgGreen, color.Bold).SprintFunc()
		cyan := color.New(color.FgCyan).SprintFunc()

		fmt.Printf("%s Snapshot captured\n", green("✓"))
		fmt.Printf("  ID:    %s\n", cyan(snap.ID))
		if label != "" {
			fmt.Printf("  Label: %s\n", cyan(label))
		}
		fmt.Printf("  Time:  %s\n", snap.CapturedAt.Format(time.RFC3339))
		fmt.Printf("  Vars:  %d\n", len(snap.Vars))

		return nil
	},
}
