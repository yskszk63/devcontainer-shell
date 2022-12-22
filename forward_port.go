package devcontainershell

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/netip"
	"os"
	"os/signal"

	"go.uber.org/zap"
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

	docker, err := resolveDocker()
	if err != nil {
		return err
	}

	// if already launched
	// `open /data/devcontainer-shell-agent: text file busy` occurres.
	install := dockerRunRm{
		image: "ghcr.io/yskszk63/devcontainer-shell-agent",
		mounts: []string{
			"type=volume,src=devcontainer-shell,dst=/data",
		},
		cmd: []string{
			"install",
			"/data/devcontainer-shell-agent",
		},
	}
	if err := docker.run(install); err != nil {
		return err
	}

	exec := dockerExec{
		containerId: container,
		bin:         "/opt/devcontainer-shell/devcontainer-shell-agent",
		args:        []string{"watch-listens"},
		notty:       true,
		noInput:     false,
	}
	return docker.runWithPipe(exec, stdin, stdout)
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

func detectIPAddress(o *dockerContainerInspectOutput) (string, error) {
	if o.NetworkSettings.IPAddress != "" {
		return o.NetworkSettings.IPAddress, nil
	}

	// try use compose default network
	pname, exists := o.Config.Labels["com.docker.compose.project"]
	if !exists {
		return "", errors.New("Could not detect ip address.")
	}

	nname := fmt.Sprintf("%s_default", pname)
	nw, exists := o.NetworkSettings.Networks[nname]
	if !exists || nw.IPAddress == "" {
		return "", errors.New("Could not detect ip address..")
	}
	return nw.IPAddress, nil
}

func localJob(cx context.Context, stdin *io.PipeWriter, stdout *io.PipeReader, container string) error {
	defer stdin.Close()
	defer stdout.Close()

	cx, cancel := context.WithCancel(cx)
	defer cancel()

	docker, err := resolveDocker()
	if err != nil {
		return err
	}

	var inspect []dockerContainerInspectOutput
	if err := docker.runWithParse(dockerContainerInspect([]string{container}), &inspect); err != nil {
		return err
	}
	if len(inspect) < 1 {
		return errors.New("no result found.")
	}
	addr, err := detectIPAddress(&inspect[0])
	if err != nil {
		return err
	}

	ports := make(map[uint16]struct{})

	listens := make(map[netip.AddrPort]context.CancelFunc)
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

		addrPort, err := event.addrPort()
		if err != nil {
			zap.L().Warn(err.Error())
			continue
		}

		if addrPort.Addr().IsLoopback() {
			continue // TODO
		}

		switch event.Type {
		case "ADD":
			{
				_, exists := ports[addrPort.Port()]
				if exists {
					// 0.0.0.0 vs ::
					continue
				}

				cx, cancel := context.WithCancel(cx)
				listens[addrPort] = cancel
				ports[addrPort.Port()] = struct{}{}
				go func() {
					if err := listen(cx, addr, addrPort.Port()); err != nil {
						log.Fatal(err)
					}
				}()
			}

		case "REMOVE":
			{
				cancel, exists := listens[addrPort]
				if !exists {
					continue
				}

				delete(listens, addrPort)
				delete(ports, addrPort.Port())
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

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
	}()

	return g.Wait()
}
