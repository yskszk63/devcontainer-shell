package devcontainershell

import "testing"

func TestDevaultSpawner(t *testing.T) {
	_, err := defaultSpawner("uname", "-a")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDevaultSpawnerError(t *testing.T) {
	_, err := defaultSpawner("false")
	if err == nil {
		t.Fail()
	}
}

func TestDefaultExecerFail(t *testing.T) {
	err := defaultExecer("/dev/null")
	if err == nil {
		t.Fail()
	}
}
