package browser

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/krnveersharma/browserdeck/schemas"
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
	caps := selenium.Capabilities{
		"browserName": "chrome",
		"goog:chromeOptions": map[string]any{
			"args": []string{
				"--no-sandbox",
				"--disable-dev-shm-usage",
				"--disable-features=SigninIntercept,SignInProfileCreation",
			},
		},
	}

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

func waitForPageLoad(driver selenium.WebDriver, timeout time.Duration) error {
	start := time.Now()
	for {
		readyState, err := driver.ExecuteScript("return document.readyState", nil)
		if err != nil {
			return fmt.Errorf("error checking page readyState: %v", err)
		}
		if readyState == "complete" {
			return nil
		}
		if time.Since(start) > timeout {
			return fmt.Errorf("timeout waiting for page to load")
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func executeInstruction(driver selenium.WebDriver, instructions []schemas.Instruction, url string) error {
	if err := driver.Get(url); err != nil {
		return fmt.Errorf("failed to open URL: %v", err)
	}
	if err := waitForPageLoad(driver, 15*time.Second); err != nil {
		return err
	}
	for _, instruction := range instructions {
		log.Printf("executing instruction: %v", instruction)
		switch instruction.Action {
		case "goto":
			if err := driver.Get(instruction.Value); err != nil {
				return fmt.Errorf("failed to open URL: %v", err)
			}
		case "click":
			element, err := driver.FindElement(selenium.ByXPATH, instruction.Selector)
			if err != nil {
				return fmt.Errorf("failed to find element: %v", err)
			}
			element.Click()
		case "type":
			element, err := driver.FindElement(selenium.ByXPATH, instruction.Selector)
			if err != nil {
				return fmt.Errorf("failed to find element: %v", err)
			}
			element.SendKeys(instruction.Value)
		case "scroll":
			js := fmt.Sprintf("window.scrollBy(0, %s);", instruction.Value)
			_, err := driver.ExecuteScript(js, nil)
			if err != nil {
				return fmt.Errorf("failed to scroll: %v", err)
			}

		case "scrollToElement":
			element, err := driver.FindElement(selenium.ByXPATH, instruction.Selector)
			if err != nil {
				return fmt.Errorf("failed to find element for scrolling: %v", err)
			}
			_, err = driver.ExecuteScript("arguments[0].scrollIntoView(true);", []interface{}{element})
			if err != nil {
				return fmt.Errorf("failed to scroll to element: %v", err)
			}
		}
		time.Sleep(10 * time.Second)
	}
	return nil
}

func (c *ChromeLauncher) Launch(sessionId int64, instructions []schemas.Instruction, url string) error {
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

	err = executeInstruction(driver, instructions, url)
	if err != nil {
		return err
	}
	// go RecordTest()

	return err
}
