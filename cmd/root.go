package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "driftmap",
	Short: "Capture, diff, and replay environment variable snapshots",
	Long: color.New(color.FgCyan, color.Bold).Sprint(`
 ____       _  __ _   __  __
|  _ \ _ __(_)/ _| |_|  \/  | __ _ _ __
| | | | '__| | |_| __| |\/| |/ _' | '_ \
| |_| | |  | |  _| |_| |  | | (_| | |_) |
|____/|_|  |_|_|  \__|_|  |_|\__,_| .__/
                                    |_|
`) + `
Track configuration drift across environments.
Capture snapshots, diff them, replay them, and find what changed.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(captureCmd)
	rootCmd.AddCommand(diffCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(replayCmd)
	rootCmd.AddCommand(watchCmd)
}
