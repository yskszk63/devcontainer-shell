package agentcmd

import (
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command {
	Use: "install [dest]",
	Args: cobra.MinimumNArgs(1),
	Run: func(c *cobra.Command, args []string) {
		dest := args[0]

		me, err := os.Executable()
		if err != nil {
			log.Fatal(err)
		}

		fp, err := os.OpenFile(dest, os.O_CREATE | os.O_WRONLY | os.O_TRUNC, 0o755)
		if err != nil {
			log.Fatal(err)
		}
		defer fp.Close()

		mefp, err := os.Open(me)
		if err != nil {
			log.Fatal(err)
		}
		defer mefp.Close()

		if _, err := io.Copy(fp, mefp); err != nil {
			log.Fatal(err)
		}
	},
}
