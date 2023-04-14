package devcontainershell

import (
	"os"
	"os/exec"
	"syscall"
)

type spawner func(cmd string, args ...string) ([]byte, error)

type execer func(cmd string, args ...string) error

func defaultSpawner(cmd string, args ...string) ([]byte, error) {
	proc := exec.Command(cmd, args...)
	proc.Stdin = nil
	proc.Stderr = os.Stderr

	raw, err := proc.Output()
	if err != nil {
		return nil, err
	}

	return raw, nil
}

func defaultExecer(cmd string, args ...string) error {
	bin, err := exec.LookPath(cmd)
	if err != nil {
		return err
	}

	a := []string{cmd}
	a = append(a, args...)

	env := os.Environ()

	return syscall.Exec(bin, a, env)
}
