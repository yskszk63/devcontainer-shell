package devcontainershell

import (
	"encoding/json"
	"os"
	"testing"
)

func TestDetectIPAddress(t *testing.T) {
	tests := []struct {
		input string
		wants string
	}{
		{"./testdata/docker-inspect/simple.json", "172.17.0.2"},
		{"./testdata/docker-inspect/compose.json", "172.18.0.2"},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			var o []dockerContainerInspectOutput
			b, err := os.ReadFile(test.input)
			if err != nil {
				t.Fatal(err)
			}
			if err := json.Unmarshal(b, &o); err != nil {
				t.Fatal(err)
			}
			if len(o) < 1 {
				t.Fatal("len")
			}

			actual, err := detectIPAddress(&o[0])
			if err != nil {
				t.Fatal(err)
			}
			if test.wants != actual {
				t.Fatalf("%s != %s", test.wants, actual)
			}
		})
	}
}
