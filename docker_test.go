package devcontainershell

import (
	"testing"
)

func TestDockerPs(t *testing.T) {
	spawner := func(cmd string, args ...string) ([]byte, error) {
		b := `{ "id": "ok" }`
		return []byte(b), nil
	}
	d := docker{
		spawner: spawner,
	}

	o, err := d.ps(dockerPsInput{})
	if err != nil {
		t.Fatal(err)
	}
	if o.ID != "ok" {
		t.Fail()
	}
}

func TestDockerPsOutputError(t *testing.T) {
	spawner := func(cmd string, args ...string) ([]byte, error) {
		b := `{`
		return []byte(b), nil
	}
	d := docker{
		spawner: spawner,
	}

	if _, err := d.ps(dockerPsInput{}); err == nil {
		t.Fail()
	}
}

func TestDockerKill(t *testing.T) {
	spawner := func(cmd string, args ...string) ([]byte, error) {
		b := `{}`
		return []byte(b), nil
	}
	d := docker{
		spawner: spawner,
	}

	if err := d.kill("id"); err != nil {
		t.Fatal(err)
	}
}
