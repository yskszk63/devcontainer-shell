package cmd

import (
	"log"
	"os"
	osexec "os/exec"

	"github.com/spf13/cobra"

	"github.com/yskszk63/devcontainer-shell"
)

var shell string
var noInject bool
var noForwardport bool

func exec(cmd *cobra.Command, args []string) {
	ds := new(devcontainershell.DevcontainerShell)

	if !noInject {
		self, err := os.Executable() // TODO REMOVE
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

	if !noForwardport {
		if noInject {
			log.Fatal("err...") // TODO
		}

		self, err := os.Executable() // TODO REMOVE
		if err != nil {
			log.Fatal(err)
		}
		proc := osexec.Command(self, "forwardport", ds.ContainerId())
		proc.Stdin = nil
		proc.Start()
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
	c.Flags().BoolVarP(&noInject, "no-inject", "I", false, "no inject agent")
	c.Flags().BoolVarP(&noForwardport, "no-forward", "F", false, "no foward port")
}

func init() {
	setupExecCmd(execCmd)
}
