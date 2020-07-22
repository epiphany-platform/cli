/*
 * Copyright © 2020 Mateusz Kyc
 */

package docker

import (
	"context"
	"github.com/docker/docker/pkg/stdcopy"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type Job struct {
	Image                string
	Command              string
	Args                 []string
	WorkDirectory        string
	EnvironmentVariables map[string]string
}

func (j Job) Run() error {
	return run(j)
}

func run(job Job) error {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}
	reader, err := cli.ImagePull(ctx, job.Image, types.ImagePullOptions{}) //TODO format output
	if err != nil {
		return err
	}
	_, _ = io.Copy(os.Stdout, reader)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:      job.Image,
		Cmd:        []string{"init"}, //TODO add arguments
		WorkingDir: job.WorkDirectory,
		Volumes:    nil, //TODO create volumes
		Env:        nil, //TODO pass ENVs
		Tty:        false,
	}, nil, nil, "")
	if err != nil {
		return err
	}
	defer cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{}) //TODO warn if error while removing

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}
	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true, Follow: true})
	if err != nil {
		return err
	}

	_, _ = stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	return nil
}
