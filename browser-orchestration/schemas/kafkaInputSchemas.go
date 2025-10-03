package schemas

type SessionMessage struct {
	SessionID    int64         `json:"sessionId"`
	Browser      string        `json:"browser"`
	Instructions []Instruction `json:"instructions"`
	Url          string        `json:"url"`
}
