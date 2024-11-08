package models

import "time"

type Article struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time
	AuthorID  int
	Title     string
	Body      string
}
