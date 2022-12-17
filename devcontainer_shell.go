package devcontainershell

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
)

func DevcontainerExec(prog string) error {
	docker, err := exec.LookPath("docker")
	if err != nil {
		return err
	}

	devcontainer, err := exec.LookPath("devcontainer")
	if err != nil {
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

	o, err := DevcontainerUp(devcontainer, wf)
	if err != nil {
		return err
	}
	if o.Outcome != "success" {
		return errors.New("failed to run `devcontainer up`")
	}

	dockerExec := DockerExec{
		Docker:      docker,
		ContainerId: o.ContainerId,
		Bin:         prog,
		Cwd:         filepath.Join(o.RemoteWorkspaceFolder, rel),
	}
	return dockerExec.Exec()
}
