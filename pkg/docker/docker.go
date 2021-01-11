package docker

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/pkg/stdcopy"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type Image struct {
	Name string
}

func (image *Image) Pull() (string, error) { //TODO remove splitting log streams here, but use zerolog multiwriter
	debug("will try to pull")
	ctx, cli, err := clientAndContext()
	if err != nil {
		return "", err
	}
	reader, err := cli.ImagePull(ctx, image.Name, types.ImagePullOptions{}) //TODO format output
	if err != nil {
		return "", err
	}
	logR, logW := io.Pipe()
	stdoutR, stdoutW := io.Pipe()

	done := make(chan bool)
	defer close(done)

	var result string
	go func() {
		s := bufio.NewScanner(stdoutR)
		for s.Scan() {
			txt := s.Text()
			debugJson([]byte(txt), "pulling")
		}
		done <- true
	}()

	go func() {
		buf := new(strings.Builder)
		_, _ = io.Copy(buf, logR)
		result = buf.String()
		done <- true
	}()

	go func() {
		defer logW.Close()
		defer stdoutW.Close()

		// build the MultiWriter for all the pipes
		mw := io.MultiWriter(logW, stdoutW)

		// copy the data into the MultiWriter
		_, _ = io.Copy(mw, reader)
	}()

	for c := 0; c < 2; c++ {
		<-done
	}

	err = reader.Close()
	if err != nil {
		return result, err
	}
	return result, nil
}

func (image *Image) IsPulled() (bool, error) {
	ctx, cli, err := clientAndContext()
	if err != nil {
		return false, err
	}
	summaries, err := cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		return false, err
	}
	for _, s := range summaries {
		for _, rt := range s.RepoTags {
			logger.Debug().Msgf("repo tag: %s", rt)
			if strings.HasSuffix(image.Name, rt) {
				return true, nil
			}
		}
	}
	return false, nil
}

type Job struct {
	Image                string
	Command              string
	Args                 []string
	WorkDirectory        string
	Mounts               map[string]string
	EnvironmentVariables map[string]string
}

func (job Job) Run() error {
	return run(job)
}

func run(job Job) error {
	ctx, cli, err := clientAndContext()
	if err != nil {
		return err
	}
	var envs []string
	for k, v := range job.EnvironmentVariables {
		envs = append(envs, fmt.Sprintf("%s=%s", k, v))
	}
	commandAndArgs := append([]string{job.Command}, job.Args...)
	var mounts []mount.Mount
	for k, v := range job.Mounts {
		mounts = append(
			mounts,
			mount.Mount{
				Type:   mount.TypeBind,
				Source: v,
				Target: k,
			})
	}
	logger.Debug().Msgf("Job run mounts: %#v", mounts)

	resp, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image:      job.Image,
			Cmd:        commandAndArgs,
			WorkingDir: job.WorkDirectory,
			Env:        envs,
			Tty:        false,
		}, &container.HostConfig{
			Mounts: mounts,
		},
		nil,
		"",
	)
	if err != nil {
		return err
	}
	defer removeFinishedContainer(cli, ctx, resp.ID)

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}
	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true, Follow: true})
	if err != nil {
		return err
	}

	_, _ = stdcopy.StdCopy(os.Stdout, os.Stderr, out) //TODO write logs to file as well

	return nil
}

func clientAndContext() (context.Context, *client.Client, error) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, nil, err
	}
	return ctx, cli, nil
}

func removeFinishedContainer(cli *client.Client, ctx context.Context, containerID string) {
	//TODO probably add check if container is running with retry because of:
	//Error response from daemon: You cannot remove a running container XXX. Stop the container before attempting removal or force remove

	err := cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{})
	if err != nil {
		warnRemovingContainer(err)
	}
}
