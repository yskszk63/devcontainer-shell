package devcontainershell

import (
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
