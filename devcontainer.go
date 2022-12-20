package devcontainershell

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
)

func resolveWorkspaceFolder(fsys fs.FS, cwd string) (string, string, error) {
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

type devcontainerUpInput struct {
	bin             string
	workspaceFolder string
	mounts          []string
	rebuild         bool
}

func (d *devcontainerUpInput) buildArgs() ([]string, error) {
	if d.workspaceFolder == "" {
		return nil, errors.New("WorkspaceFolder must set.")
	}

	ret := []string{
		"up",
		"--workspace-folder",
		d.workspaceFolder,
	}

	for _, mount := range d.mounts {
		ret = append(ret, "--mount", mount)
	}

	if d.rebuild {
		ret = append(ret, "--remove-existing-container")
	}

	return ret, nil
}

type devcontainerUpOutput struct {
	Outcome               string `json:"outcome"`
	ContainerId           string `json:"containerId"`
	RemoteUser            string `json:"remoteUser"`
	RemoteWorkspaceFolder string `json:"remoteWorkspaceFolder"`
}

func devcontainerUp(input devcontainerUpInput) (*devcontainerUpOutput, error) {
	args, err := input.buildArgs()
	if err != nil {
		return nil, err
	}

	if zap.L().Level().Enabled(zap.DebugLevel) {
		zap.L().Debug(fmt.Sprintf("%s %s", input.bin, strings.Join(args, " ")))
	}

	proc := exec.Command(input.bin, args...)
	proc.Stdin = nil
	proc.Stderr = os.Stderr

	raw, err := proc.Output()
	if err != nil {
		return nil, err
	}

	var o devcontainerUpOutput
	if err := json.Unmarshal(raw, &o); err != nil {
		return nil, err
	}
	return &o, nil
}
