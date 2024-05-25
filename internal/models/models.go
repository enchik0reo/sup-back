package models

import "time"

type SupInfo struct {
	ID    int64  `json:"id"`
	Name  string `json:"model_name"`
	Price int64  `json:"price"`
}

type Sup struct {
	SupInfo
	ReservedDays []time.Time `json:"reserved_days"`
}

type ApproveSup struct {
	ID   int64     `json:"id"`
	Name string    `json:"model_name"`
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

type Approve struct {
	ID           int64        `json:"id"`
	ClientNumber string       `json:"client_phone"`
	ClientName   string       `json:"client_name"`
	SupsInfo     []ApproveSup `json:"sups_info"`
	FullPrice    int64        `json:"price"`
}

type ApproveData struct {
	ID           int64  `json:"id"`
	ClientNumber string `json:"client_phone"`
}

type Reserved struct {
	Day        time.Time
	ModelID    int64
	ModelName  string
	ModelPrice int64
	ApproveID  int64
}

// Telegram types ...

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
