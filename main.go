package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

func resolveWorkspaceFolder() (string, string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", "", err
	}

	wf := cwd
	rel := ""
	for {
		_, err := os.Stat(filepath.Join(wf, ".devcontainer"))
		if !os.IsNotExist(err) {
			return wf, rel, nil
		}

		if wf == "/" {
			return "", "", errors.New("working-folder not found.")
		}

		rel = filepath.Join(filepath.Base(wf), rel)
		wf = filepath.Dir(wf)
	}
}

func devcontainerUp(devcontainer, workingFolder string) (*devcontainerUpOutput, error) {
	proc := exec.Command(devcontainer, "--workspace-folder", workingFolder, "up")
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

func dockerExec(bin, docker string, output *devcontainerUpOutput, rel string) error {
	// TODO execve? on xxix
	proc := exec.Command(docker, "exec", "-it", "-u", output.RemoteUser, "-w", filepath.Join(output.RemoteWorkspaceFolder, rel), output.ContainerId, bin)
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

	wf, rel, err := resolveWorkspaceFolder()
	if err != nil {
		return err
	}

	output, err := devcontainerUp(devcontainer, wf)
	if err != nil {
		return err
	}
	if output.Outcome != "success" {
		return errors.New("failed to `devcontainer up`")
	}

	if err := dockerExec("bash", docker, output, rel); err != nil {
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
