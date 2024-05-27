package models

type Type int

const (
	Unknown Type = iota
	Message
)

type Meta struct {
	ChatID   int
	UserName string
}

type Event struct {
	Type Type
	Text string
	Meta Meta
}

type UpdatesResponce struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result,omitempty"`
}

type Update struct {
	ID      int              `json:"update_id"`
	Message *IncomingMessage `json:"message,omitempty"`
}

type IncomingMessage struct {
	Text string `json:"text"`
	From From   `json:"from"`
	Chat Chat   `json:"chat"`
}

type From struct {
	UserName string `json:"username"`
}

type Chat struct {
	ID int `json:"id"`
}
