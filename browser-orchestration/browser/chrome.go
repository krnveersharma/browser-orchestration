package browser

import (
	"context"
	"fmt"
	"log"

	mobycontainer "github.com/moby/moby/api/types/container"
	mobyclient "github.com/moby/moby/client"
)

type ChromeLauncher struct{}

func (c *ChromeLauncher) Launch(sessionId int64) error {
	cli, err := mobyclient.NewClientWithOpts(mobyclient.FromEnv, mobyclient.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	resp, err := cli.ContainerCreate(
		context.Background(),
		&mobycontainer.Config{
			Image: "selenium/standalone-chrome",
			Env:   []string{fmt.Sprintf("SESSION_ID=%d", sessionId)},
		},
		nil, nil, nil, "",
	)
	if err != nil {
		return err
	}

	if err := cli.ContainerStart(context.Background(), resp.ID, mobyclient.ContainerStartOptions{}); err != nil {
		return err
	}

	log.Println("Started Chrome container with ID:", resp.ID)
	return nil
}
