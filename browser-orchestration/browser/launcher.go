package browser

import "github.com/krnveersharma/browserdeck/schemas"

type BrowserLauncher interface {
	Launch(sessionId int64, instructions []schemas.Instruction, url string) error
}
