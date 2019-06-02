package models

type Link struct {
	Original  string      `json:"original"`
	Processed string			`json:"processed"`
	TTL 			int64				`json:"ttl"`
}
