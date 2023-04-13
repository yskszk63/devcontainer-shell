package devcontainershell

import "testing"

func TestDevaultSpawner(t *testing.T) {
	_, err := defaultSpawner("uname", "-a")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDefaultExecerFail(t *testing.T) {
	err := defaultExecer("/dev/null")
	if err == nil {
		t.Fail()
	}
}
