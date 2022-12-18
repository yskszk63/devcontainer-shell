package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/yskszk63/devcontainer-shell"
)

var shell string
var noInject bool

func exec(cmd *cobra.Command, args []string) {
	ds := new(devcontainershell.DevcontainerShell)

	if !noInject {
		self, err := os.Executable()
		if err != nil {
			log.Fatal(err)
		}
		if err := ds.Inject(self); err != nil {
			log.Fatal(err)
		}
	}

	if err := ds.Up(); err != nil {
		log.Fatal(err)
	}

	if err := ds.Exec(shell); err != nil {
		log.Fatal()
	}
}

var execCmd = &cobra.Command {
	Use: "exec",
	Short: "execute shell on devcontainer",
	Run: exec,
}

func setupExecCmd(c *cobra.Command) {
	c.Flags().StringVarP(&shell, "shell", "s", "bash", "using shell program")
	c.Flags().BoolVarP(&noInject, "no inject", "I", false, "no inject agent")
}

func init() {
	setupExecCmd(execCmd)
}
