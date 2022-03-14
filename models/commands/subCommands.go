package commands

type SubCommands struct {
	Command string      `json:"command"`
	Payload interface{} `json:"payload"`
}
