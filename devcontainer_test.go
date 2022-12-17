package devcontainershell

import (
	"os"
	"testing"
)

func TestResolveWorkspaceFolder(t *testing.T) {
	tests := []struct {
		name     string
		root     string
		dir      string
		err      string
		wantsWf  string
		wantsRel string
	}{
		{
			name:     "exists",
			root:     "./testdata/exists/",
			dir:      "/a",
			wantsWf:  "/a",
			wantsRel: "",
		},
		{
			name:     "exists2",
			root:     "./testdata/exists/",
			dir:      "/a/b",
			wantsWf:  "/a",
			wantsRel: "b",
		},
		{
			name:     "exists3",
			root:     "./testdata/exists/",
			dir:      "/a/b/c",
			wantsWf:  "/a",
			wantsRel: "b/c",
		},
		{
			name: "notexists",
			root: "./testdata/notexists/",
			dir:  "/a/b/c",
			err:  "workspace-folder not found.",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fsys := os.DirFS(test.root)
			wf, rel, err := ResolveWorkspaceFolder(fsys, test.dir)
			if err != nil {
				if err.Error() != test.err {
					t.Fatal(err)
				}
				return
			}

			if wf != test.wantsWf {
				t.Errorf("%s != %s", wf, test.wantsWf)
			}
			if rel != test.wantsRel {
				t.Errorf("%s != %s", rel, test.wantsRel)
			}
		})
	}
}
