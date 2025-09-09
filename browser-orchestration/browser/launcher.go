package browser

type BrowserLauncher interface {
	Launch(sessionId int64, instructions, url string) error
}
