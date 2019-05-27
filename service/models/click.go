package models

type Click struct {
	Link      string      `json:"processed"`
	Timestamp  string     `json:"timestamp"`
	Owner      string     `json:"owner"`
}
