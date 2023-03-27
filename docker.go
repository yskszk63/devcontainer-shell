package devcontainershell

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"go.uber.org/zap"
)

type buildArgs interface {
	buildArgs() ([]string, error)
}

type docker string

func resolveDocker() (*docker, error) {
	d, err := exec.LookPath("docker")
	if err != nil {
		return nil, err
	}

	ret := docker(d)
	return &ret, nil
}

func (d *docker) run(input buildArgs) error {
	args, err := input.buildArgs()
	if err != nil {
		return err
	}

	if zap.L().Core().Enabled(zap.DebugLevel) {
		zap.L().Debug(fmt.Sprintf("%s %s", *d, strings.Join(args, " ")))
	}
	proc := exec.Command(string(*d), args...)
	proc.Stdin = os.Stdin
	proc.Stdout = os.Stdout
	proc.Stderr = os.Stderr

	return proc.Run()
}

func (d *docker) runWithPipe(input buildArgs, stdin io.Reader, stdout io.Writer) error {
	args, err := input.buildArgs()
	if err != nil {
		return err
	}

	if zap.L().Core().Enabled(zap.DebugLevel) {
		zap.L().Debug(fmt.Sprintf("%s %s", *d, strings.Join(args, " ")))
	}
	proc := exec.Command(string(*d), args...)
	proc.Stdin = stdin
	proc.Stdout = stdout
	proc.Stderr = os.Stderr

	return proc.Run()
}

func (d *docker) runWithParse(input buildArgs, output any) error {
	args, err := input.buildArgs()
	if err != nil {
		return err
	}

	if zap.L().Core().Enabled(zap.DebugLevel) {
		zap.L().Debug(fmt.Sprintf("%s %s", *d, strings.Join(args, " ")))
	}
	proc := exec.Command(string(*d), args...)
	proc.Stdin = os.Stdin
	proc.Stderr = os.Stderr

	raw, err := proc.Output()
	if err != nil {
		return err
	}

	if err := json.Unmarshal(raw, output); err != nil {
		return err
	}

	return nil
}

type dockerExec struct {
	containerId string
	bin         string
	args        []string
	user        string
	cwd         string
	notty       bool
	noInput     bool
}

func (d dockerExec) buildArgs() ([]string, error) {
	if d.containerId == "" {
		return nil, errors.New("containerId must set.")
	}
	if d.bin == "" {
		return nil, errors.New("bin must set.")
	}

	args := []string{
		"exec",
	}

	if !d.notty {
		args = append(args, "-t")
	}

	if !d.noInput {
		args = append(args, "-i")
	}

	if d.user != "" {
		args = append(args, "-u", d.user)
	}

	if d.cwd != "" {
		args = append(args, "-w", d.cwd)
	}

	args = append(args, d.containerId, d.bin)
	if d.args != nil {
		args = append(args, d.args...)
	}

	return args, nil
}

type dockerContainerInspect []string

func (d dockerContainerInspect) buildArgs() ([]string, error) {
	if d == nil {
		return nil, errors.New("Must not nil.")
	}

	args := []string{
		"container",
		"inspect",
	}
	for _, c := range d {
		args = append(args, c)
	}

	return args, nil
}

type dockerContainerInspectOutput struct {
	Id     string
	Config struct {
		Labels map[string]string
	}
	Mounts []struct {
		Type        string
		Name        string
		Destination string
	}
}

type dockerPs struct {
	filter []string
	quiet  bool
}

func (d dockerPs) buildArgs() ([]string, error) {
	args := []string{
		"ps",
	}

	if d.filter != nil {
		for _, f := range d.filter {
			args = append(args, "-f", f)
		}
	}

	if d.quiet {
		args = append(args, "-q")
	}

	return args, nil
}

type dockerKill []string

func (d dockerKill) buildArgs() ([]string, error) {
	args := []string{
		"kill",
	}

	args = append(args, []string(d)...)

	return args, nil
}

type dockerComposeKill struct {
	project *string
}

func (d dockerComposeKill) buildArgs() ([]string, error) {
	args := []string{
		"compose",
	}

	if d.project != nil {
		args = append(args, "-p", *d.project)
	}

	args = append(args, "kill")

	return args, nil
}
