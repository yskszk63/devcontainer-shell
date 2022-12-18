package devcontainershell

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"
)

type DockerRunRm struct {
	Docker string
	Image  string
	Mounts []string
	Cmd    []string
}

func (d *DockerRunRm) buildArgs() ([]string, error) {
	if d.Image == "" {
		return nil, errors.New("image must set.")
	}

	args := []string{
		"run",
		"--rm",
	}
	if d.Mounts != nil {
		for _, m := range d.Mounts {
			args = append(args, "--mount", m)
		}
	}
	args = append(args, d.Image)
	if d.Cmd != nil {
		args = append(args, d.Cmd...)
	}
	return args, nil
}

func (d *DockerRunRm) Run() error {
	if d.Docker == "" {
		return errors.New("docker must set.")
	}
	args, err := d.buildArgs()
	if err != nil {
		return err
	}

	proc := exec.Command(d.Docker, args...)
	proc.Stdin = os.Stdin
	proc.Stdout = os.Stdout
	proc.Stderr = os.Stderr

	return proc.Run()
}

type DockerExec struct {
	Docker      string
	ContainerId string
	Bin         string
	Args        []string
	User        string
	Cwd         string
	Notty       bool
	NoInput     bool
}

func (d *DockerExec) buildArgs() ([]string, error) {
	if d.ContainerId == "" {
		return nil, errors.New("containerId must set.")
	}
	if d.Bin == "" {
		return nil, errors.New("bin must set.")
	}

	args := []string{
		"exec",
	}

	if !d.Notty {
		args = append(args, "-t")
	}

	if !d.NoInput {
		args = append(args, "-i")
	}

	if d.User != "" {
		args = append(args, "-u", d.User)
	}

	if d.Cwd != "" {
		args = append(args, "-w", d.Cwd)
	}

	args = append(args, d.ContainerId, d.Bin)
	if d.Args != nil {
		args = append(args, d.Args...)
	}

	return args, nil
}

func (d *DockerExec) Exec() error {
	if d.Docker == "" {
		return errors.New("docker must set.")
	}
	args, err := d.buildArgs()
	if err != nil {
		return err
	}

	proc := exec.Command(d.Docker, args...)
	proc.Stdin = os.Stdin
	proc.Stdout = os.Stdout
	proc.Stderr = os.Stderr

	return proc.Run()
}

func (d *DockerExec) ExecWithPipe(cx context.Context, stdin io.Reader, stdout io.Writer) error {
	if d.Docker == "" {
		return errors.New("docker must set.")
	}
	args, err := d.buildArgs()
	if err != nil {
		return err
	}

	proc := exec.Command(d.Docker, args...)
	proc.Stdin = stdin
	proc.Stdout = stdout
	proc.Stderr = os.Stderr

	return proc.Run()
}

func (d *DockerExec) SyscallExec() error {
	if d.Docker == "" {
		return errors.New("docker must set.")
	}
	args, err := d.buildArgs()
	if err != nil {
		return err
	}
	argv := make([]string, 0, len(args)+1)
	argv = append(argv, d.Docker)
	argv = append(argv, args...)

	return syscall.Exec(d.Docker, argv, os.Environ())
}

func DockerVolumeCreate(docker, name string) error {
	proc := exec.Command(docker, "volume", "create", name)
	proc.Stdin = nil
	proc.Stdout = nil
	proc.Stderr = os.Stderr

	return proc.Run()
}

func DockerCp(docker, srcPath, container, destPath string) error {
	proc := exec.Command(docker, "cp", srcPath, fmt.Sprintf("%s:%s", container, destPath))
	proc.Stdin = nil
	proc.Stdout = nil
	proc.Stderr = os.Stderr

	return proc.Run()
}

type DockerContainerInspectOutput struct {
	NetworkSettings struct {
		IPAddress string
	}
}

func DockerContainerInspect(docker, container string) (*DockerContainerInspectOutput, error) {
	proc := exec.Command(docker, "container", "inspect", container)
	proc.Stdin = nil
	proc.Stderr = os.Stderr

	raw, err := proc.Output()
	if err != nil {
		return nil, err
	}

	var o []DockerContainerInspectOutput
	if err := json.Unmarshal(raw, &o); err != nil {
		return nil, err
	}
	if len(o) < 1 {
		return nil, errors.New("empty result")
	}
	return &o[0], nil
}
