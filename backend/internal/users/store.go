package users

import (
	"database/sql"
	"errors"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) Get(userID int) (*User, error) {
	row := s.db.QueryRow("SELECT id, email, password, created_at FROM users WHERE id = $1 LIMIT 1", userID)

	user := new(User)

	if err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (s *Store) GetByEmail(email string) (*User, error) {
	row := s.db.QueryRow(
		"SELECT id, email, password, created_at FROM users WHERE email = $1 LIMIT 1",
		email,
	)

	user := new(User)

	if err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil

}
func (s *Store) Create(user *User) error {
	row := s.db.QueryRow(
		"INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id, created_at",
		user.Email,
		user.Password,
	)

	if err := row.Scan(
		&user.ID,
		&user.CreatedAt,
	); err != nil {
		return err
	}

	return nil
}
