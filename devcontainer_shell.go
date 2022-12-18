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
	dockerPath           string
	devcontainerPath     string
	containerCwd         string
	injectBin            string
}

func (d *DevcontainerShell) ContainerId() string {
	return d.devcontainerUpOutput.ContainerId
}

func (d *DevcontainerShell) ensureDockerResolved() error {
	if d.dockerPath != "" {
		return nil
	}

	docker, err := exec.LookPath("docker")
	if err != nil {
		return err
	}

	d.dockerPath = docker
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

func (d *DevcontainerShell) Inject(self string) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if err := d.ensureResolvePaths(); err != nil {
		return err
	}

	if err := DockerVolumeCreate(d.dockerPath, "devcontainer-shell"); err != nil {
		return err
	}

	d.injectBin = self
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
	if d.injectBin != "" {
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

	dockerExec := DockerExec{
		Docker:      d.dockerPath,
		ContainerId: d.devcontainerUpOutput.ContainerId,
		Bin:         prog,
		Cwd:         d.containerCwd,
	}
	return dockerExec.Exec()
}
