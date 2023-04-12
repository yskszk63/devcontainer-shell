package devcontainershell

import (
	"fmt"
	"os"
)

type DevcontainerShell struct {
	devcontainer devcontainer
	docker       docker
	relativePath string
}

func NewDevcontainerShell() (*DevcontainerShell, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	wf, rel, err := resolveWorkspaceFolder(os.DirFS("/"), cwd)
	if err != nil {
		return nil, err
	}

	devcontainer := devcontainer{
		workspaceFolder: wf,
		spawner:         defaultSpawner,
		execer:          defaultExecer,
	}

	docker := docker{
		spawner: defaultSpawner,
	}

	return &DevcontainerShell{
		devcontainer: devcontainer,
		docker:       docker,
		relativePath: rel,
	}, nil
}

func (d *DevcontainerShell) Exec(removeExistingContainer bool, cmd string, args ...string) error {
	r, err := d.devcontainer.up(removeExistingContainer)
	if err != nil {
		return err
	}

	return d.devcontainer.exec(r.ContainerId, cmd, args...)
}

func (d *DevcontainerShell) Kill() error {
	r, err := d.docker.ps(dockerPsInput{
		noTrunc: true,
		filter:  fmt.Sprintf("label=devcontainer.local_folder=%s", d.devcontainer.workspaceFolder),
	})
	if err != nil {
		return err
	}

	if r == nil {
		return nil
	}

	if err := d.docker.kill(r.ID); err != nil {
		return err
	}

	return nil
}
