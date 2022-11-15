package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
)

func lookupPaths() (string, string, error) {
	devcontainer, err := exec.LookPath("devcontainer")
	if err != nil {
		return "", "", err
	}
	docker, err := exec.LookPath("docker")
	if err != nil {
		return "", "", err
	}
	return devcontainer, docker, nil
}

type devcontainerUpOutput struct {
	Outcome               string `json:"outcome"`
	ContainerId           string `json:"containerId"`
	RemoteUser            string `json:"remoteUser"`
	RemoteWorkspaceFolder string `json:"remoteWorkspaceFolder"`
}

func devcontainerUp(devcontainer string) (*devcontainerUpOutput, error) {
	proc := exec.Command(devcontainer, "--workspace-folder", ".", "up")
	proc.Stderr = os.Stderr
	raw, err := proc.Output()
	if err != nil {
		return nil, err
	}

	output := new(devcontainerUpOutput)
	if err := json.Unmarshal(raw, output); err != nil {
		return nil, err
	}

	return output, nil
}

func dockerExec(bin, docker string, output *devcontainerUpOutput) error {
	// TODO execve? on xxix
	proc := exec.Command(docker, "exec", "-it", "-u", output.RemoteUser, "-w", output.RemoteWorkspaceFolder, output.ContainerId, bin)
	proc.Stdin = os.Stdin
	proc.Stdout = os.Stdout
	proc.Stderr = os.Stderr

	if err := proc.Run(); err != nil {
		return err
	}
	return nil
}

func run() error {
	devcontainer, docker, err := lookupPaths()
	if err != nil {
		return err
	}

	output, err := devcontainerUp(devcontainer)
	if err != nil {
		return err
	}
	if output.Outcome != "success" {
		return errors.New("failed to `devcontainer up`")
	}
	fmt.Printf("%s %s %s\n", output.ContainerId, output.RemoteUser, output.RemoteWorkspaceFolder)

	if err := dockerExec("bash", docker, output); err != nil {
		return err
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(-1)
	}
}
