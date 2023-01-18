package devcontainershell

import (
	"errors"
	"fmt"
	"os/exec"

	"go.uber.org/zap"
)

func tryRunForwardServer(docker *docker, o *devcontainerUpOutput) error {
	var outputs []dockerContainerInspectOutput
	if err := docker.runWithParse(dockerContainerInspect{o.ContainerId}, &outputs); err != nil {
		return err
	}
	if len(outputs) != 1 {
		return errors.New("len(outputs) != 1")
	}

	output := outputs[0]

	conname := "devcontainer-shell-portforward"
	for key, val := range output.Config.Labels {
		if key != "com.docker.compose.project" {
			continue
		}
		conname = fmt.Sprintf("devcontainer-shell-portforward-%s", val)
	}

	var volname string
	for _, vol := range output.Mounts {
		if vol.Type != "volume" {
			continue
		}
		if vol.Destination != "/run/devcontainer-portforward" {
			continue
		}
		volname = vol.Name
	}

	err := docker.runSilently(dockerRunRm{
		image:  "ghcr.io/yskszk63/devcontainer-portforward-server",
		name:   conname,
		mounts: []string{fmt.Sprintf("type=volume,source=%s,target=/data", volname)},
		net:    "host",
		detach: true,
	}, !zap.L().Level().Enabled(zap.DebugLevel))
	if err != nil {
		_, known := err.(*exec.ExitError)
		if !known {
			return err
		}
	}

	return nil
}
