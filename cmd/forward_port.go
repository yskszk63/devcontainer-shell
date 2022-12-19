package cmd

import (
	"fmt"
	"os"
	osexec "os/exec"
	"strings"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/yskszk63/devcontainer-shell"
)

var nodaemon bool

func daemonize(container string) error {
	me, err := os.Executable()
	if err != nil {
		return err
	}

	args := []string{
		"forwardport",
		container,
		"-D",
	}
	if zap.L().Level().Enabled(zap.DebugLevel) {
		zap.L().Debug(fmt.Sprintf("%s %s", me, strings.Join(args, " ")))
	}

	proc := osexec.Command(me, args...)
	proc.Stdin = nil
	proc.Stdout = nil
	proc.Stderr = nil
	if err := proc.Start(); err != nil {
		return err
	}

	os.Exit(0)
	return nil
}

func forwardPort(cmd *cobra.Command, args []string) error {
	container := args[0]

	if !nodaemon {
		if err := daemonize(container); err != nil {
			return err
		}
	}

	if err := devcontainershell.ForwardPort(container); err != nil {
		return err
	}

	return nil
}

var forwardPortCmd = &cobra.Command {
	Use: "forwardport [container id]",
	Short: "port forward",
	Args: cobra.MinimumNArgs(1),
	SilenceErrors: true,
	SilenceUsage: true,
	RunE: forwardPort,
}

func init() {
	forwardPortCmd.Flags().BoolVarP(&nodaemon, "no-daemon", "D", false, "disable daemonize")
}
