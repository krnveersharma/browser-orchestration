package browser

type BrowserLauncher interface {
	Launch(sessionId int64) error
}
