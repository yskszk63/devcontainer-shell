package devcontainershell

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"os/exec"
)

type buildArgs interface {
	buildArgs() ([]string, error)
}

type docker string

func resolveDocker() (*docker, error) {
	d, err := exec.LookPath("docker")
	if err != nil {
		return nil, err
	}

	ret := docker(d)
	return &ret, nil
}

func (d *docker) run(input buildArgs) error {
	args, err := input.buildArgs()
	if err != nil {
		return err
	}

	proc := exec.Command(string(*d), args...)
	proc.Stdin = os.Stdin
	proc.Stdout = os.Stdout
	proc.Stderr = os.Stderr

	return proc.Run()
}

func (d *docker) runWithPipe(input buildArgs, stdin io.Reader, stdout io.Writer) error {
	args, err := input.buildArgs()
	if err != nil {
		return err
	}

	proc := exec.Command(string(*d), args...)
	proc.Stdin = stdin
	proc.Stdout = stdout
	proc.Stderr = os.Stderr

	return proc.Run()
}

func (d *docker) runWithParse(input buildArgs, output any) error {
	args, err := input.buildArgs()
	if err != nil {
		return err
	}

	proc := exec.Command(string(*d), args...)
	proc.Stdin = os.Stdin
	proc.Stderr = os.Stderr

	raw, err := proc.Output()
	if err != nil {
		return err
	}

	if err := json.Unmarshal(raw, output); err != nil {
		return err
	}

	return nil
}

type dockerRunRm struct {
	image  string
	mounts []string
	cmd    []string
}

func (d dockerRunRm) buildArgs() ([]string, error) {
	if d.image == "" {
		return nil, errors.New("image must set.")
	}

	args := []string{
		"run",
		"--rm",
	}
	if d.mounts != nil {
		for _, m := range d.mounts {
			args = append(args, "--mount", m)
		}
	}
	args = append(args, d.image)
	if d.cmd != nil {
		args = append(args, d.cmd...)
	}
	return args, nil
}

type dockerExec struct {
	containerId string
	bin         string
	args        []string
	user        string
	cwd         string
	notty       bool
	noInput     bool
}

func (d dockerExec) buildArgs() ([]string, error) {
	if d.containerId == "" {
		return nil, errors.New("containerId must set.")
	}
	if d.bin == "" {
		return nil, errors.New("bin must set.")
	}

	args := []string{
		"exec",
	}

	if !d.notty {
		args = append(args, "-t")
	}

	if !d.noInput {
		args = append(args, "-i")
	}

	if d.user != "" {
		args = append(args, "-u", d.user)
	}

	if d.cwd != "" {
		args = append(args, "-w", d.cwd)
	}

	args = append(args, d.containerId, d.bin)
	if d.args != nil {
		args = append(args, d.args...)
	}

	return args, nil
}

type dockerVolumeCreate string

func (d dockerVolumeCreate) buildArgs() ([]string, error) {
	return []string{
		"volume",
		"create",
		string(d),
	}, nil
}

type dockerContainerInspect string

func (d dockerContainerInspect) buildArgs() ([]string, error) {
	return []string{
		"container",
		"inspect",
		string(d),
	}, nil
}

type dockerContainerInspectOutput struct {
	NetworkSettings struct {
		IPAddress string
	}
}
