package orm

import "time"

type Device struct {
	ID           int        `json:"id"`
	Name         string     `json:"name"`
	Info         *string    `json:"info"`
	LastActivity *time.Time `json:"lastActivity"`
	CreatedAt    time.Time  `json:"createdAt"`
}
