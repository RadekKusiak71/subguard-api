package users

import "time"

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type UserStore interface {
	Get(userID int) (*User, error)
	GetByEmail(email string) (*User, error)
	Create(user *User) error
}

type userContextKey string

const UserContextKey userContextKey = "userID"
