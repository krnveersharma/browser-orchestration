package browser

func GetLauncher(browser string) BrowserLauncher {
	return &ChromeLauncher{}
}
