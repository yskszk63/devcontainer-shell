package devcontainershell

import (
	"encoding/json"
)

type docker struct{
	spawner spawner
}

type dockerPsInput struct{
	noTrunc bool
	filter string
}

type dockerPsOutput struct{
	Command string
	CreatedAt string
	ID string
	Image string
	Labels string
	LocalVolumes string
	Mounts string
	Names string
	Networks string
	Ports string
	RunningFor string
	Size string
	State string
	Status string
}

func (d *docker) ps(input dockerPsInput) (*dockerPsOutput, error) {
	args := []string{
		"ps",
		"--format",
		"json",
		"--latest",
	}

	if input.noTrunc {
		args = append(args, "--no-trunc")
	}

	if input.filter != "" {
		args = append(args, "--filter", input.filter)
	}

	b, err := d.spawner("docker", args...)
	if err != nil {
		return nil, err
	}

	if len(b) == 0 {
		return nil, nil
	}

	r := new(dockerPsOutput)
	if err := json.Unmarshal(b, r); err != nil {
		return nil, err
	}
	return r, nil
}

func (d *docker) kill(id string) error {
	_, err := d.spawner("docker", "kill", id)
	if err != nil {
		return err
	}

	return nil
}
