package devcontainershell

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
)

func ResolveWorkspaceFolder(fsys fs.FS, cwd string) (string, string, error) {
	wf := cwd
	rel := ""

	for {
		_, err := fs.Stat(fsys, filepath.Join(wf, ".devcontainer/devcontainer.json")[1:])
		if err != nil && !os.IsNotExist(err) {
			return "", "", err
		}
		if err == nil {
			return wf, rel, nil
		}

		if wf == "/" {
			return "", "", errors.New("workspace-folder not found.")
		}

		rel = filepath.Join(filepath.Base(wf), rel)
		wf = filepath.Dir(wf)
	}
}

type DevcontainerUpOutput struct {
	Outcome               string `json:"outcome"`
	ContainerId           string `json:"containerId"`
	RemoteUser            string `json:"remoteUser"`
	RemoteWorkspaceFolder string `json:"remoteWorkspaceFolder"`
}

func DevcontainerUp(bin, workspaceFolder string) (*DevcontainerUpOutput, error) {
	proc := exec.Command(bin, "--workspace-folder", workspaceFolder, "up")
	proc.Stdin = nil
	proc.Stderr = os.Stderr

	raw, err := proc.Output()
	if err != nil {
		return nil, err
	}

	var o DevcontainerUpOutput
	if err := json.Unmarshal(raw, &o); err != nil {
		return nil, err
	}
	return &o, nil
}