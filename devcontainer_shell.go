package devcontainershell

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

type DevcontainerShell struct {
	mutex                sync.Mutex
	devcontainerUpOutput *DevcontainerUpOutput
	docker               *docker
	devcontainerPath     string
	containerCwd         string
	inject               bool
}

func (d *DevcontainerShell) ContainerId() string {
	return d.devcontainerUpOutput.ContainerId
}

func (d *DevcontainerShell) ensureDockerResolved() error {
	if d.docker != nil {
		return nil
	}

	docker, err := resolveDocker()
	if err != nil {
		return err
	}

	d.docker = docker
	return nil
}

func (d *DevcontainerShell) ensureDevcontainerResolved() error {
	if d.devcontainerPath != "" {
		return nil
	}

	devcontainer, err := exec.LookPath("devcontainer")
	if err != nil {
		return err
	}

	d.devcontainerPath = devcontainer
	return nil
}

func (d *DevcontainerShell) ensureResolvePaths() error {
	if err := d.ensureDockerResolved(); err != nil {
		return err
	}
	if err := d.ensureDevcontainerResolved(); err != nil {
		return err
	}
	return nil
}

func (d *DevcontainerShell) Inject() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if err := d.ensureResolvePaths(); err != nil {
		return err
	}

	if err := d.docker.run(dockerVolumeCreate("devcontainer-shell")); err != nil {
		return err
	}

	d.inject = true
	return nil
}

func (d *DevcontainerShell) Up() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if err := d.ensureResolvePaths(); err != nil {
		return err
	}

	root := os.DirFS("/")
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	wf, rel, err := ResolveWorkspaceFolder(root, cwd)
	if err != nil {
		return err
	}

	mounts := make([]string, 0)
	if d.inject {
		mounts = append(mounts, fmt.Sprintf("type=volume,source=devcontainer-shell,target=/opt/devcontainer-shell"))
	}
	o, err := DevcontainerUp(DevcontainerUpInput{
		Bin:             d.devcontainerPath,
		WorkspaceFolder: wf,
		Mounts:          mounts,
	})
	if err != nil {
		return err
	}
	if o.Outcome != "success" {
		return errors.New("failed to run `devcontainer up`")
	}

	d.devcontainerUpOutput = o
	d.containerCwd = filepath.Join(o.RemoteWorkspaceFolder, rel)

	return nil
}

func (d *DevcontainerShell) Exec(prog string) error {
	if d.devcontainerUpOutput == nil {
		return errors.New("must call Up() before")
	}

	dockerExec := dockerExec{
		containerId: d.devcontainerUpOutput.ContainerId,
		bin:         prog,
		cwd:         d.containerCwd,
		user:        d.devcontainerUpOutput.RemoteUser,
	}
	return d.docker.run(dockerExec)
}
