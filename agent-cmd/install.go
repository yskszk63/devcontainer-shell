package agentcmd

import (
	"io"
	"os"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command {
	Use: "install [dest]",
	Args: cobra.MinimumNArgs(1),
	SilenceErrors: true,
	SilenceUsage: true,
	RunE: func(c *cobra.Command, args []string) error{
		dest := args[0]

		me, err := os.Executable()
		if err != nil {
			return err
		}

		fp, err := os.OpenFile(dest, os.O_CREATE | os.O_WRONLY | os.O_TRUNC, 0o755)
		if err != nil {
			return err
		}
		defer fp.Close()

		mefp, err := os.Open(me)
		if err != nil {
			return err
		}
		defer mefp.Close()

		if _, err := io.Copy(fp, mefp); err != nil {
			return err
		}
		return nil
	},
}
