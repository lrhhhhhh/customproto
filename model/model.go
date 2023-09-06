package model

type Message struct {
	Id      int    `json:"Id,omitempty"`
	Content string `json:"Content,omitempty"`
}
