package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var debug bool

var rootCmd = &cobra.Command {
	Use: "devcontainer-shell",
	Short: "devcontainer shell helper",
	SilenceErrors: true,
	SilenceUsage: true,
	PersistentPreRun: func(c *cobra.Command, args []string) {
		if debug {
			logger := zap.Must(zap.NewDevelopmentConfig().Build())
			zap.ReplaceGlobals(logger)
		}
	},
	RunE: exec,
}

func init() {
	setupExecCmd(rootCmd)

	rootCmd.AddCommand(execCmd)
	rootCmd.AddCommand(forwardPortCmd)

	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Debug")
}

func Execute() error {
	return rootCmd.Execute()
}
