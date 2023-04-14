package devcontainershell

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

/**
 * returns (workspacefolder, relative path, error)
 */
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
