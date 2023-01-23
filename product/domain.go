package product

import "time"

type Product struct {
	ID        int64     `json:"-"`
	Name      string    `json:"name" validate:"required"`
	Slug      string    `json:"slug"`
	Price     int       `json:"price" validate:"required,number"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
