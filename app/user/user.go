package user

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// User model
type User struct {
	ID        int       `db:"id" json:"-"`
	UUID      uuid.UUID `db:"uuid" json:"uuid"` // UUID field used to avoid exposing auto increment PKs
	FirstName string    `db:"first_name" json:"firstName"`
	LastName  string    `db:"last_name" json:"lastName"`
	Email     string    `db:"email" json:"email"`
	IsActive  bool      `db:"is_active" json:"isActive"`
	Created   time.Time `db:"created" json:"created"`
	Modified  time.Time `db:"modified" json:"modified"`
}
