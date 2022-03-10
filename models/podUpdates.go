package models

type PodUpdates struct {
	PodName string `json:"podName"`
	PodIp   string `json:"podIp"`
	State   string `json:"state"`
}
