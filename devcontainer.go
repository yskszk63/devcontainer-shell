package devcontainershell

import (
	"encoding/json"
)

type devcontainer struct {
	workspaceFolder string
	spawner         spawner
	execer          execer
}

func newDevcontainer(workspaceFolder string) (*devcontainer, error) {
	return &devcontainer{
		workspaceFolder: workspaceFolder,
		spawner:         defaultSpawner,
		execer:          defaultExecer,
	}, nil
}

type devcontainerUpOutput struct {
	Outcome               string `json:"outcome"`
	ContainerId           string `json:"containerId"`
	RemoteUser            string `json:"remoteUser"`
	RemoteWorkspaceFolder string `json:"remoteWorkspaceFolder"`
}

func (d *devcontainer) up(removeExistingContainer bool) (*devcontainerUpOutput, error) {
	args := []string{
		"up",
		"--workspace-folder",
		d.workspaceFolder,
	}
	if removeExistingContainer {
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

func (d *devcontainer) exec(containerId, cmd string, args ...string) error {
	a := []string{
		"exec",
		"--workspace-folder",
		d.workspaceFolder,
		cmd,
	}
	a = append(a, args...)

	return d.execer("devcontainer", a...)
}
