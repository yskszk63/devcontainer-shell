package devcontainershell

import (
	"reflect"
	"testing"
)

func TestBuildArgs(t *testing.T) {
	tests := []struct {
		name   string
		target DockerExec
		wants  []string
		err    string
	}{
		{"invalid1", DockerExec{}, nil, "containerId must set."},

		{"invalid2", DockerExec{
			Docker: "a",
		}, nil, "containerId must set."},

		{"invalid2", DockerExec{
			Docker:      "a",
			ContainerId: "b",
		}, nil, "bin must set."},

		{"bin only", DockerExec{
			Docker:      "a",
			ContainerId: "b",
			Bin:         "c",
		}, []string{"exec", "-it", "b", "c"}, ""},

		{"with full", DockerExec{
			Docker:      "a",
			ContainerId: "b",
			Bin:         "c",
			User:        "u",
			Cwd:         "w",
			Args:        []string{"a1", "a2"},
		}, []string{"exec", "-it", "-u", "u", "-w", "w", "b", "c", "a1", "a2"}, ""},

		{"no tty", DockerExec{
			Docker:      "a",
			ContainerId: "b",
			Bin:         "c",
			Notty:       true,
		}, []string{"exec", "b", "c"}, ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.target.buildArgs()
			if err != nil {
				if err.Error() != test.err {
					t.Fatal(err)
				}
				return
			}

			if !reflect.DeepEqual(got, test.wants) {
				t.Errorf("%#v != %#v", got, test.wants)
			}
		})
	}
}
