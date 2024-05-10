package common

import (
	"time"
)

type Event struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	State     string    `json:"state"`
}
