package devcontainershell

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

type DockerExec struct {
	Docker      string
	ContainerId string
	Bin         string
	Args        []string
	User        string
	Cwd         string
	Notty       bool
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
		args = append(args, "-it")
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
