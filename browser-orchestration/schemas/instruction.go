package schemas

type Instruction struct {
	Action   string `json:"action"`
	Value    string `json:"value"`
	Selector string `json:"selector"`
}
