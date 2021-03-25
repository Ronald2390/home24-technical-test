package model

import "time"

// User represents the user
// swagger:model
type User struct {
	ID        int        `json:"-" db:"id"`
	Name      string     `json:"name" db:"name"`
	Email     string     `json:"email" db:"email"`
	Address   string     `json:"address" db:"address"`
	Password  string     `json:"-" db:"password"`
	CreatedBy int        `json:"-" db:"createdBy"`
	CreatedAt time.Time  `json:"-" db:"createdAt"`
	UpdatedAt time.Time  `json:"-" db:"updatedAt"`
	UpdatedBy int        `json:"-" db:"updatedBy"`
	DeletedAt *time.Time `json:"-" db:"deletedAt"`
	DeletedBy *int       `json:"-" db:"deletedBy"`
}
