package devcontainershell

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

func TestDockerExec(t *testing.T) {
	tests := []struct {
		name   string
		target dockerExec
		wants  []string
		err    string
	}{
		{"invalid1", dockerExec{}, nil, "containerId must set."},

		{"invalid2", dockerExec{
			containerId: "b",
		}, nil, "bin must set."},

		{"bin only", dockerExec{
			containerId: "b",
			bin:         "c",
		}, []string{"exec", "-t", "-i", "b", "c"}, ""},

		{"with full", dockerExec{
			containerId: "b",
			bin:         "c",
			user:        "u",
			cwd:         "w",
			args:        []string{"a1", "a2"},
		}, []string{"exec", "-t", "-i", "-u", "u", "-w", "w", "b", "c", "a1", "a2"}, ""},

		{"no tty", dockerExec{
			containerId: "b",
			bin:         "c",
			notty:       true,
		}, []string{"exec", "-i", "b", "c"}, ""},
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

func TestDockerContainerInspectOutput(t *testing.T) {
	tests := []string{
		"./testdata/docker-inspect/simple.json",
		"./testdata/docker-inspect/compose.json",
	}

	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			var o []dockerContainerInspectOutput
			b, err := os.ReadFile(test)
			if err != nil {
				t.Fatal(err)
			}
			if err := json.Unmarshal(b, &o); err != nil {
				t.Fatal(err)
			}
		})
	}
}
