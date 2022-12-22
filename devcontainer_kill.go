package devcontainershell

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"golang.org/x/sync/errgroup"
)

func Kill() error {
	root := os.DirFS("/")
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	wf, _, err := resolveWorkspaceFolder(root, cwd)
	if err != nil {
		return err
	}

	docker, err := resolveDocker()
	if err != nil {
		return err
	}

	r, w := io.Pipe()
	g := errgroup.Group{}

	c := make(chan string, 1)

	g.Go(func() error {
		defer w.Close()

		ps := dockerPs{
			filter: []string{
				fmt.Sprintf("label=devcontainer.local_folder=%s", wf),
			},
			quiet: true,
		}

		return docker.runWithPipe(ps, nil, w)
	})

	g.Go(func() error {
		defer r.Close()
		defer close(c)

		s := bufio.NewScanner(r)
		for s.Scan() {
			c <- s.Text()
		}
		return s.Err()
	})

	ids := make([]string, 0)
	for line := range c {
		ids = append(ids, line)
	}

	if err := g.Wait(); err != nil {
		return err
	}

	var inspect []dockerContainerInspectOutput
	if err := docker.runWithParse(dockerContainerInspect(ids), &inspect); err != nil {
		return err
	}

	for _, info := range inspect {
		pj, exists := info.Config.Labels["com.docker.compose.project"]
		if exists {
			docker.run(dockerComposeKill{project: &pj})
			continue
		}
		docker.run(dockerKill{info.Id})
	}

	return nil
}
