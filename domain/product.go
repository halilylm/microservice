package domain

import "time"

type Product struct {
	ID        int64     `json:"-"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Price     int       `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
