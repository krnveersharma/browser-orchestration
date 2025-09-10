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

type ChromeLauncher struct {
	cli *mobyclient.Client
}

func NewChromeLauncher() (*ChromeLauncher, error) {
	cli, err := mobyclient.NewClientWithOpts(mobyclient.FromEnv, mobyclient.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &ChromeLauncher{cli: cli}, nil
}

func (c *ChromeLauncher) startContainer(sessionId int64) (string, string, error) {
	port, _ := nat.NewPort("tcp", "4444")

	resp, err := c.cli.ContainerCreate(
		context.Background(),
		&mobycontainer.Config{
			Image: "selenium/standalone-chrome",
			Env:   []string{fmt.Sprintf("SESSION_ID=%d", sessionId)},
		},
		&mobycontainer.HostConfig{
			PortBindings: nat.PortMap{
				port:       []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: ""}},     // random host port
				"5900/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "5900"}}, // optional VNC debug
			},
		},
		nil, nil, "",
	)
	if err != nil {
		return "", "", err
	}

	if err := c.cli.ContainerStart(context.Background(), resp.ID, mobyclient.ContainerStartOptions{}); err != nil {
		return "", "", err
	}
	log.Println("Started Chrome container with ID:", resp.ID)

	inspect, err := c.cli.ContainerInspect(context.Background(), resp.ID)
	if err != nil {
		return "", "", err
	}
	hostPort := inspect.NetworkSettings.Ports[port][0].HostPort
	return resp.ID, hostPort, nil
}

func (c *ChromeLauncher) stopContainer(containerID string) {
	if err := c.cli.ContainerStop(context.Background(), containerID, mobyclient.ContainerStopOptions{}); err != nil {
		log.Printf("Failed to stop container %s: %v", containerID, err)
	}
	if err := c.cli.ContainerRemove(context.Background(), containerID, mobyclient.ContainerRemoveOptions{Force: true}); err != nil {
		log.Printf("Failed to remove container %s: %v", containerID, err)
	} else {
		log.Println("Removed container:", containerID)
	}
}

func connectWebDriver(port string) (selenium.WebDriver, error) {
	caps := selenium.Capabilities{"browserName": "chrome"}
	url := fmt.Sprintf("http://localhost:%s/wd/hub", port)

	for range 10 {
		driver, err := selenium.NewRemote(caps, url)
		if err == nil {
			return driver, nil
		}
		time.Sleep(2 * time.Second)
	}
	return nil, fmt.Errorf("selenium not ready on port %s", port)
}

func executeInstruction(driver selenium.WebDriver, instruction, url string) error {
	if err := driver.Get(url); err != nil {
		return fmt.Errorf("failed to open URL: %v", err)
	}

	switch instruction {
	case "scrollAndLike":
		driver.ExecuteScript("window.scrollBy(0,800);", nil)
		time.Sleep(3 * time.Second)

		elem, err := driver.FindElement(selenium.ByXPATH, `(//a[@id="thumbnail"])[1]`)
		if err == nil {
			elem.Click()
		}
		time.Sleep(5 * time.Second)

		likeBtn, err := driver.FindElement(selenium.ByXPATH, `//ytd-toggle-button-renderer[@id="like-button"]//button`)
		if err == nil {
			likeBtn.Click()
		}
	default:
		log.Printf("Instruction %s not recognized", instruction)
	}
	return nil
}

func (c *ChromeLauncher) Launch(sessionId int64, instruction, url string) error {
	containerID, port, err := c.startContainer(sessionId)
	if err != nil {
		return err
	}
	defer c.stopContainer(containerID)

	driver, err := connectWebDriver(port)
	if err != nil {
		return err
	}
	defer driver.Quit()

	return executeInstruction(driver, instruction, url)
}
