package models

type Link struct {
	Owner     string      `json:"owner"`
	Original  string      `json:"original"`
	Processed string			`json:"processed"`
	TTL 			int64				`json:"ttl"`
}
