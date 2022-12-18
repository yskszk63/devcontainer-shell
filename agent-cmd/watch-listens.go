package agentcmd

import (
	"context"
	"log"
	"os"

	"github.com/yskszk63/devcontainer-shell"

	"github.com/spf13/cobra"
)

var watchListensCmd = &cobra.Command {
	Use: "watch-listens",
	Run: func(c *cobra.Command, args []string) {
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
			log.Fatal(err)
		}
	},
}
