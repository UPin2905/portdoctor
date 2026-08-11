package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "0.1.1"

var rootCmd = &cobra.Command{
	Use:   "portdoctor [port]",
	Short: "Diagnose local port conflicts without the detective work",
	Long: `🩺 PortDoctor — Find out what's using your port, and why.

A fast, cross-platform CLI that tells you what is using a local port,
where the process came from, and what you can safely do about it.`,
	Version: version,
	Args:    cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 {
			return runInspect(args[0])
		}
		return cmd.Help()
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.AddCommand(inspectCmd)
	rootCmd.AddCommand(scanCmd)
	rootCmd.AddCommand(killCmd)
	rootCmd.AddCommand(findCmd)
}

// Execute là entry point chính của CLI.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
