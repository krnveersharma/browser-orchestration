package browser

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/docker/go-connections/nat"
	mobycontainer "github.com/moby/moby/api/types/container"
	mobyclient "github.com/moby/moby/client"
	"github.com/tebeka/selenium"
)

type ChromeLauncher struct{}

func (c *ChromeLauncher) Launch(sessionId int64, instruction, url string) error {
	hostPort := "0"
	port, _ := nat.NewPort("tcp", "4444")
	cli, err := mobyclient.NewClientWithOpts(mobyclient.FromEnv, mobyclient.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	// remove this 5900(it is for testing), eerytime new container is made will forward to 5900
	resp, err := cli.ContainerCreate(
		context.Background(),
		&mobycontainer.Config{
			Image: "selenium/standalone-chrome",
			Env:   []string{fmt.Sprintf("SESSION_ID=%d", sessionId)},
		},
		&mobycontainer.HostConfig{
			PortBindings: nat.PortMap{
				port:       []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: hostPort}},
				"5900/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "5900"}},
			},
		},
		nil, nil, "",
	)
	if err != nil {
		return err
	}

	if err := cli.ContainerStart(context.Background(), resp.ID, mobyclient.ContainerStartOptions{}); err != nil {
		return err
	}
	log.Println("Started Chrome container with ID:", resp.ID)

	inspect, err := cli.ContainerInspect(context.Background(), resp.ID)
	hostPort = inspect.NetworkSettings.Ports[port][0].HostPort
	seleniumURL := fmt.Sprintf("http://localhost:%s/wd/hub", hostPort)
	caps := selenium.Capabilities{"browserName": "chrome"}

	time.Sleep(5 * time.Second) // waiting for  selenium to be set up
	webDriver, err := selenium.NewRemote(caps, seleniumURL)
	if err != nil {
		return fmt.Errorf("failed to connect to Selenium: %v", err)
	}
	defer webDriver.Quit()

	if err := webDriver.Get(url); err != nil {
		return fmt.Errorf("failed to open URL: %v", err)
	}

	defer func() {
		if err := cli.ContainerStop(context.Background(), resp.ID, mobyclient.ContainerStopOptions{}); err != nil {
			log.Printf("Failed to stop container %s: %v", resp.ID, err)
		}

		if err := cli.ContainerRemove(context.Background(), resp.ID, mobyclient.ContainerRemoveOptions{Force: true}); err != nil {
			log.Printf("Failed to remove container %s: %v", resp.ID, err)
		} else {
			log.Println("Removed container:", resp.ID)
		}
		log.Println("removed container")
	}()

	// this is testing code, i will add main instructions by user
	webDriver.ExecuteScript("window.scrollBy(0,800);", nil)
	time.Sleep(10 * time.Second)

	elem, _ := webDriver.FindElement(selenium.ByXPATH, `(//a[@id="thumbnail"])[1]`)
	elem.Click()
	time.Sleep(20 * time.Second)

	likeBtn, _ := webDriver.FindElement(selenium.ByXPATH, `//ytd-toggle-button-renderer[@id="like-button"]//button`)
	likeBtn.Click()
	time.Sleep(5 * time.Second)

	return nil
}
