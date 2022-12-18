package cmd

import (
	"log"
	"os"
	osexec "os/exec"

	"github.com/spf13/cobra"

	"github.com/yskszk63/devcontainer-shell"
)

var nodaemon bool

func daemonize(container string) error {
	me, err := os.Executable()
	if err != nil {
		return err
	}

	proc := osexec.Command(me, "forwardport", container, "-D")
	proc.Stdin = nil
	proc.Stdout = nil
	proc.Stderr = nil
	if err := proc.Start(); err != nil {
		return err
	}

	os.Exit(0)
	return nil
}

func forwardPort(cmd *cobra.Command, args []string) {
	container := args[0]

	if !nodaemon {
		if err := daemonize(container); err != nil {
			log.Fatal(err)
		}
	}

	if err := devcontainershell.ForwardPort(container); err != nil {
		log.Fatal(err)
	}
}

var forwardPortCmd = &cobra.Command {
	Use: "forwardport [container id]",
	Short: "port forward",
	Args: cobra.MinimumNArgs(1),
	Run: forwardPort,
}

func init() {
	forwardPortCmd.Flags().BoolVarP(&nodaemon, "no-daemon", "D", false, "disable daemonize")
}
