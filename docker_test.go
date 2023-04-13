package devcontainershell

import (
	"testing"
)

func TestDockerPs(t *testing.T) {
	spawner := func(cmd string, args ...string) ([]byte, error) {
		b := `{}`
		return []byte(b), nil
	}
	d := docker{
		spawner: spawner,
	}

	d.ps(dockerPsInput{})
}

func TestDockerKill(t *testing.T) {
	spawner := func(cmd string, args ...string) ([]byte, error) {
		b := `{}`
		return []byte(b), nil
	}
	d := docker{
		spawner: spawner,
	}

	d.kill("id")
}
