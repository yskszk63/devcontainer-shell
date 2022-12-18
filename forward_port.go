package devcontainershell

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"

	"golang.org/x/sync/errgroup"
)

func remoteJob(cx context.Context, stdin *io.PipeReader, stdout *io.PipeWriter, container string) error {
	defer stdin.Close()
	defer stdout.Close()

	go func() {
		defer stdin.Close()
		defer stdout.Close()
		<-cx.Done()
	}()

	docker, err := exec.LookPath("docker")
	if err != nil {
		return err
	}

	// if already launched
	// `open /data/devcontainer-shell-agent: text file busy` occurres.
	install := DockerRunRm{
		Docker: docker,
		Image:  "ghcr.io/yskszk63/devcontainer-shell-agent",
		Mounts: []string{
			"type=volume,src=devcontainer-shell,dst=/data",
		},
		Cmd: []string{
			"install",
			"/data/devcontainer-shell-agent",
		},
	}
	if err := install.Run(); err != nil {
		return err
	}

	exec := DockerExec{
		Docker:      docker,
		ContainerId: container,
		Bin:         "/opt/devcontainer-shell/devcontainer-shell-agent",
		Args:        []string{"watch-listens"},
		Notty:       true,
		NoInput:     false,
	}
	return exec.ExecWithPipe(cx, stdin, stdout)
}

func forward(cx context.Context, addr string, port uint16, sock net.Conn) error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		return err
	}

	g, cx := errgroup.WithContext(cx)

	g.Go(func() error {
		if _, err := io.Copy(conn, sock); err != nil {
			return err
		}
		return nil
	})
	g.Go(func() error {
		if _, err := io.Copy(sock, conn); err != nil {
			return err
		}
		return nil
	})

	return g.Wait()
}

func listen(cx context.Context, addr string, port uint16) error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	defer l.Close()

	go func() {
		<-cx.Done()
		l.Close()
	}()

	for {
		sock, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return nil
		}

		go func() {
			defer sock.Close()

			if err := forward(cx, addr, port, sock); err != nil {
				log.Println(err)
			}
		}()
	}
}

func localJob(cx context.Context, stdin *io.PipeWriter, stdout *io.PipeReader, container string) error {
	defer stdin.Close()
	defer stdout.Close()

	cx, cancel := context.WithCancel(cx)
	defer cancel()

	docker, err := exec.LookPath("docker")
	if err != nil {
		return err
	}

	inspect, err := DockerContainerInspect(docker, container)
	if err != nil {
		return err
	}
	addr := inspect.NetworkSettings.IPAddress

	listens := make(map[uint16]context.CancelFunc)
	dec := json.NewDecoder(stdout)
	for {
		select {
		case <-cx.Done():
			return nil
		default:
		}

		var event ListenEvent
		err := dec.Decode(&event)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		if event.IP != "0.0.0.0" {
			continue // TODO
		}
		port := event.Port

		switch event.Type {
		case "ADD":
			{
				cx, cancel := context.WithCancel(cx)
				listens[port] = cancel
				go func() {
					if err := listen(cx, addr, port); err != nil {
						log.Fatal(err)
					}
				}()
			}

		case "REMOVE":
			{
				cancel, exists := listens[port]
				if !exists {
					continue
				}

				delete(listens, port)
				cancel()
			}
		}
		fmt.Printf("%#v\n", event)
	}
}

func ForwardPort(container string) error {
	cx, cancel := context.WithCancel(context.Background())
	defer cancel()

	g, cx := errgroup.WithContext(cx)
	r, w := io.Pipe()
	r2, w2 := io.Pipe()

	g.Go(func() error {
		defer cancel()
		return remoteJob(cx, r2, w, container)
	})
	g.Go(func() error {
		defer cancel()
		return localJob(cx, w2, r, container)
	})

	go func() {
		defer cancel()

		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt)
		<-c
	}()

	return g.Wait()
}
