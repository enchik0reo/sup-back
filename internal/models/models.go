package models

import "time"

type Sup struct {
	ID           int64       `json:"id"`
	Name         string      `json:"model_name"`
	Price        int64       `json:"price"`
	ReservedDays []time.Time `json:"reserved_days"`
}

type ApproveSup struct {
	ID   int64     `json:"id"`
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

type Approve struct {
	ClientNumber string       `json:"client_phone"`
	ClientName   string       `json:"client_name"`
	SupsInfo     []ApproveSup `json:"sups_info"`
	FullPrice    int64        `json:"price"`
}

type Reserved struct {
	Day        time.Time
	ModelID    int64
	ModelName  string
	ModelPrice int64
}
