package commands

type AdminCommand struct {
	Command      string                 `json:"command"`
	Data         map[string]interface{} `json:"data"`
	Audience     []string               `json:"audience"`
	AudienceType string                 `json:"audienceType"`
}
