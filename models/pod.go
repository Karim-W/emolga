package models

import "github.com/karim-w/emolga/modules/scclient"

type Pod struct {
	PodName string `json:"id"`
	PodIp   string `json:"ip"`
	Client  scclient.Client
}
