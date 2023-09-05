package model

type File struct {
	filename string
	filesize int
	path     string
}

type Message struct {
	Id      int    `json:"Id,omitempty"`
	Content string `json:"Content,omitempty"`
}
