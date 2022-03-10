package models

import "github.com/sacOO7/socketcluster-client-go/scclient"

type Pod struct {
	PodName string `json:"id"`
	PodIp   string `json:"ip"`
	Client  scclient.Client
}
