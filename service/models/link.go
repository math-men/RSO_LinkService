package models

type Link struct {
  ID        int         `json: "id"`
	Owner     string      `json:"owner"`
	Original  string      `json:"original"`
	Processed string			`json:"processed"`
	Cost      int			    `json:"cost"`
}
