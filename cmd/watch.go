package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/DucChau/driftmap/internal/differ"
	"github.com/DucChau/driftmap/internal/snapshot"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var watchInterval int

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch for environment changes in real-time",
	Long:  `Polls the current environment at a fixed interval and prints any changes detected.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		interval := time.Duration(watchInterval) * time.Second
		if interval < time.Second {
			interval = 2 * time.Second
		}

		yellow := color.New(color.FgYellow, color.Bold).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()
		red := color.New(color.FgRed).SprintFunc()
		dim := color.New(color.Faint).SprintFunc()

		fmt.Fprintf(os.Stderr, "%s Watching environment every %s (Ctrl+C to stop)\n",
			yellow("◉"), interval)

		prev := snapshot.New("", os.Environ())

		for {
			time.Sleep(interval)
			curr := snapshot.New("", os.Environ())
			result := differ.Diff(prev, curr)

			if len(result.Added)+len(result.Removed)+len(result.Changed) == 0 {
				fmt.Printf("%s [%s] no changes\n",
					dim("·"), dim(time.Now().Format("15:04:05")))
			} else {
				fmt.Printf("\n%s [%s] %d change(s) detected:\n",
					yellow("▲"), time.Now().Format("15:04:05"),
					len(result.Added)+len(result.Removed)+len(result.Changed))

				for k, v := range result.Added {
					fmt.Printf("  %s %s=%s\n", green("+"), k, v)
				}
				for k, v := range result.Removed {
					fmt.Printf("  %s %s=%s\n", red("-"), k, v)
				}
				for k, c := range result.Changed {
					fmt.Printf("  ~ %s: %s → %s\n", k, red(c.Before), green(c.After))
				}
				prev = curr
			}
		}
	},
}

func init() {
	watchCmd.Flags().IntVarP(&watchInterval, "interval", "i", 2, "Poll interval in seconds")
}
