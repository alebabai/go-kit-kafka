package domain

import (
	"time"
)

type Event struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Expired   bool      `json:"expired"`
}
