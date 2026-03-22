package models

type InfoResponse struct {
	Version string `json:"version"`
	Service string `json:"service"`
	Author  string `json:"author"`
}
