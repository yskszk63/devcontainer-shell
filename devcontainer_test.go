package devcontainershell

import (
	"testing"
)

func TestDevcontainerUp(t *testing.T) {
	spawner := func(cmd string, args ...string) ([]byte, error) {
		b := `{}`
		return []byte(b), nil
	}

	d := devcontainer{
		spawner: spawner,
	}

	_, err := d.up(devcontainerUpInput{
		removeExistingContainer: true,
	})

	if err != nil {
		t.Fatal(err)
	}
}

func TestDevcontainerExec(t *testing.T) {
	execer := func(cmd string, args ...string) error {
		return nil
	}

	d := devcontainer{
		execer: execer,
	}

	if err := d.exec(devcontainerExecInput{
		containerId: "1",
	}); err != nil {
		t.Fatal(err)
	}
}
