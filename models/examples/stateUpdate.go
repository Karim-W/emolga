package examples

type StateUpdateExample struct {
	Command      string    `json:"command"`
	Data         StateData `json:"data"`
	Audience     []string  `json:"audience"`
	AudienceType string    `json:"audienceType"`
}

type StateData struct {
	State string `json:"state"`
}
