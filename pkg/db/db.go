package db

import (
	"time"
)

type DbTime struct {
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type IdType struct {
	Id uint64 `json:"id,omitempty"`
}

type Program struct {
	IdType
	Name    string  `json:"name"`
	Content string  `json:"content"`
	Logs    *string `json:"logs"`
	Status  uint64  `json:"status"`
	DbTime
}
