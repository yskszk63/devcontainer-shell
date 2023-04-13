package devcontainershell

import (
	"encoding/json"
)

type devcontainer struct {
	workspaceFolder string
	spawner         spawner
	execer          execer
}

type devcontainerUpInput struct {
	removeExistingContainer bool
}

type devcontainerUpOutput struct {
	Outcome               string `json:"outcome"`
	ContainerId           string `json:"containerId"`
	RemoteUser            string `json:"remoteUser"`
	RemoteWorkspaceFolder string `json:"remoteWorkspaceFolder"`
}

func (d *devcontainer) up(input devcontainerUpInput) (*devcontainerUpOutput, error) {
	args := []string{
		"up",
		"--workspace-folder",
		d.workspaceFolder,
	}
	if input.removeExistingContainer {
		args = append(args, "--remove-existing-container")
	}
	o, err := d.spawner("devcontainer", args...)
	if err != nil {
		return nil, err
	}

	ret := new(devcontainerUpOutput)
	if err := json.Unmarshal(o, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

type devcontainerExecInput struct {
	containerId string
	cmd         string
	args        []string
}

func (d *devcontainer) exec(input devcontainerExecInput) error {
	a := []string{
		"exec",
		"--workspace-folder",
		d.workspaceFolder,
		input.cmd,
	}
	a = append(a, input.args...)

	return d.execer("devcontainer", a...)
}
