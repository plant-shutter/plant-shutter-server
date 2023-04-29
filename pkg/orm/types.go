package orm

import "time"

type Device struct {
	ID           int        `json:"id"`
	Name         string     `json:"name"`
	Info         string     `json:"info"`
	UserID       int        `json:"userID"`
	LastActivity *time.Time `json:"lastActivity"`
	CreatedAt    time.Time  `json:"createdAt"`
}

type Image struct {
	ID        int       `json:"ID"`
	ProjectID int       `json:"projectID"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

type Project struct {
	ID        int       `json:"ID"`
	UserID    int       `json:"userID"`
	DeviceID  int       `json:"deviceID"`
	Name      string    `json:"name"`
	Info      string    `json:"info"`
	CreatedAt time.Time `json:"createdAt"`
}
