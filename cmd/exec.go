package cmd

import (
	"errors"
	"fmt"
	"os"
	osexec "os/exec"
	"strings"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/yskszk63/devcontainer-shell"
)

var shell string
var noInject bool
var noForwardport bool
var rebuild bool

func exec(cmd *cobra.Command, args []string) error {
	ds := new(devcontainershell.DevcontainerShell)
	if rebuild {
		ds.Rebuild = true
	}

	if !noInject {
		if err := ds.Inject(); err != nil {
			return err
		}
	}

	if err := ds.Up(); err != nil {
		return err
	}

	if !noForwardport {
		if noInject {
			return errors.New("err...") // TODO
		}

		self, err := os.Executable() // TODO REMOVE
		if err != nil {
			return err
		}

		args := []string{
			"forwardport",
			ds.ContainerId(),
		}
		if debug {
			// no daemon
			args = append(args, "-D", "-d")
		}

		if zap.L().Level().Enabled(zap.DebugLevel) {
			zap.L().Debug(fmt.Sprintf("%s %s", self, strings.Join(args, " ")))
		}

		proc := osexec.Command(self, args...)
		proc.Stdin = nil
		proc.Stdout = os.Stdout
		proc.Stderr = os.Stderr
		if err := proc.Start(); err != nil {
			return err
		}

		if debug {
			defer func() {
				if proc.Process == nil {
					return
				}
				proc.Process.Kill()
			}()
		}
	}

	if err := ds.Exec(shell); err != nil {
		return err
	}

	return nil
}

var execCmd = &cobra.Command {
	Use: "exec",
	Short: "execute shell on devcontainer",
	SilenceErrors: true,
	SilenceUsage: true,
	RunE: exec,
}

func setupExecCmd(c *cobra.Command) {
	c.Flags().StringVarP(&shell, "shell", "s", "bash", "using shell program")
	c.Flags().BoolVarP(&noInject, "no-inject", "I", false, "no inject agent")
	c.Flags().BoolVarP(&noForwardport, "no-forward", "F", false, "no foward port")
	c.Flags().BoolVarP(&rebuild, "rebuild", "b", false, "remove existing container")
}

func init() {
	setupExecCmd(execCmd)
}
