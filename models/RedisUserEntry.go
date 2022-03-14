package models

type RedisUserEntry struct {
	Type           string   `json:"type"`
	Id             string   `json:"id"`
	Sessions       []string `json:"sessions"`
	Hearings       []string `json:"hearings"`
	State          string   `json:"state"`
	UserEmail      string   `json:"userEmail"`
	ServerInstance string   `json:"serverInstance"`
	SocketId       string   `json:"socketId"`
	FirstName      string   `json:"firstName"`
	LastName       string   `json:"lastName"`
}
