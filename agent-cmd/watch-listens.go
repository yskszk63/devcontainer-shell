package agentcmd

import (
	"context"
	"os"

	"github.com/yskszk63/devcontainer-shell"

	"github.com/spf13/cobra"
)

var watchListensCmd = &cobra.Command {
	Use: "watch-listens",
	SilenceErrors: true,
	SilenceUsage: true,
	RunE: func(c *cobra.Command, args []string) error {
		cx, cancel := context.WithCancel(context.Background())

		go func() {
			defer cancel()

			buf := make([]byte, 1)
			for {
				_, err := os.Stdin.Read(buf)
				if err != nil {
					return
				}
			}
		}()

		err := devcontainershell.WatchListens(cx)
		if err != nil {
			return err
		}
		return nil
	},
}
