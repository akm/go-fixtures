package models

import "time"

type User struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Email     string
}
