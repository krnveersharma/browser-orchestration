package browser

import "fmt"

func GetLauncher(browser string) BrowserLauncher {
	switch browser {
	case "chrome":
		launcher, err := NewChromeLauncher()
		if err != nil {
			panic(fmt.Sprintf("failed to create chrome launcher: %v", err))
		}
		return launcher
	default:
		launcher, err := NewChromeLauncher()
		if err != nil {
			panic(fmt.Sprintf("failed to create chrome launcher: %v", err))
		}
		return launcher
	}
}
