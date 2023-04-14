package devcontainershell

import (
	"errors"
	"testing"
)

func TestDevcontainerShellExec(t *testing.T) {
	spawner := func(cmd string, args ...string) ([]byte, error) {
		b := `{}`
		return []byte(b), nil
	}

	execer := func(cmd string, args ...string) error {
		return nil
	}

	d := DevcontainerShell{
		devcontainer: devcontainer{
			spawner: spawner,
			execer:  execer,
		},
	}

	err := d.Exec(DevcontainerShellExecInput{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestDevcontainerShellExecError(t *testing.T) {
	spawner := func(cmd string, args ...string) ([]byte, error) {
		return nil, errors.New("ERR")
	}

	execer := func(cmd string, args ...string) error {
		return nil
	}

	d := DevcontainerShell{
		devcontainer: devcontainer{
			spawner: spawner,
			execer:  execer,
		},
	}

	err := d.Exec(DevcontainerShellExecInput{})
	if err == nil {
		t.Fail()
	}
}

func TestDevcontainerShellKill(t *testing.T) {
	spawner := func(cmd string, args ...string) ([]byte, error) {
		b := `{}`
		return []byte(b), nil
	}

	d := DevcontainerShell{
		docker: docker{
			spawner: spawner,
		},
	}

	err := d.Kill()
	if err != nil {
		t.Fatal(err)
	}
}

func TestDevcontainerShellKillError(t *testing.T) {
	spawner := func(cmd string, args ...string) ([]byte, error) {
		return nil, errors.New("ERR")
	}

	d := DevcontainerShell{
		docker: docker{
			spawner: spawner,
		},
	}

	err := d.Kill()
	if err == nil {
		t.Fail()
	}
}

func TestDevcontainerShellNotUpYet(t *testing.T) {
	spawner := func(cmd string, args ...string) ([]byte, error) {
		if args[0] != "ps" {
			return nil, errors.New("ERR")
		}

		b := ""
		return []byte(b), nil
	}

	d := DevcontainerShell{
		docker: docker{
			spawner: spawner,
		},
	}

	err := d.Kill()
	if err != nil {
		t.Fatal(err)
	}
}

func TestDevcontainerShellKillError2(t *testing.T) {
	spawner := func(cmd string, args ...string) ([]byte, error) {
		if args[0] == "ps" {
			b := "{}"
			return []byte(b), nil
		}

		return nil, errors.New("ERR")
	}

	d := DevcontainerShell{
		docker: docker{
			spawner: spawner,
		},
	}

	err := d.Kill()
	if err == nil {
		t.Fail()
	}
}
