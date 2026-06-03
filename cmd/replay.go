package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/DucChau/driftmap/internal/store"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var replayCmd = &cobra.Command{
	Use:   "replay <snapshot-id> -- <command> [args...]",
	Short: "Run a command with a snapshot's environment injected",
	Long: `Replay a captured snapshot by injecting its environment variables
into a subprocess. The current environment is replaced entirely.

Example:
  driftmap replay abc123 -- env | grep PATH
  driftmap replay prod-snap -- ./start.sh`,
	DisableFlagParsing: false,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("usage: driftmap replay <snapshot-id> -- <command> [args...]")
		}

		snapID := args[0]
		subArgs := args[1:]

		s, err := store.Open()
		if err != nil {
			return err
		}

		snap, err := s.Resolve(snapID)
		if err != nil {
			return err
		}

		envSlice := make([]string, 0, len(snap.Vars))
		keys := make([]string, 0, len(snap.Vars))
		for k := range snap.Vars {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			envSlice = append(envSlice, k+"="+snap.Vars[k])
		}

		cyan := color.New(color.FgCyan).SprintFunc()
		fmt.Fprintf(os.Stderr, "↪ Replaying snapshot %s with %d vars → %s\n",
			cyan(snap.ShortID()), len(snap.Vars), strings.Join(subArgs, " "))

		subCmd := exec.Command(subArgs[0], subArgs[1:]...)
		subCmd.Env = envSlice
		subCmd.Stdin = os.Stdin
		subCmd.Stdout = os.Stdout
		subCmd.Stderr = os.Stderr

		return subCmd.Run()
	},
}
