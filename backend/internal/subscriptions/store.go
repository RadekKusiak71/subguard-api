package subscriptions

import (
	"database/sql"
	"errors"
)

var (
	ErrSubscriptionNotFound = errors.New("subscription not found")
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) List(userID int) ([]Subscription, error) {
	rows, err := s.db.Query(
		`SELECT id, user_id, name, price, plan, next_payment_at, created_at 
		FROM subscriptions
		WHERE user_id = $1`,
		userID,
	)

	if err != nil {
		return nil, err
	}

	var subscriptions []Subscription

	for rows.Next() {
		var subscription Subscription

		if err := rows.Scan(
			&subscription.ID,
			&subscription.UserID,
			&subscription.Name,
			&subscription.Price,
			&subscription.Plan,
			&subscription.NextPaymentAt,
			&subscription.CreatedAt,
		); err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, subscription)
	}

	if len(subscriptions) == 0 {
		return []Subscription{}, nil
	}

	return subscriptions, nil
}

func (s *Store) GetByName(userID int, name string) (*Subscription, error) {
	row := s.db.QueryRow(
		`SELECT id, user_id, name, price, plan, next_payment_at, created_at 
		FROM subscriptions
		WHERE user_id = $1 AND name = $2 LIMIT 1`,
		userID,
		name,
	)

	subscription := new(Subscription)

	if err := row.Scan(
		&subscription.ID,
		&subscription.UserID,
		&subscription.Name,
		&subscription.Price,
		&subscription.Plan,
		&subscription.NextPaymentAt,
		&subscription.CreatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSubscriptionNotFound
		}
		return nil, err
	}

	return subscription, nil
}

func (s *Store) Create(subscription *Subscription) error {

	row := s.db.QueryRow(
		`INSERT INTO subscriptions (user_id, name, price, plan, next_payment_at) 
		VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`,
		subscription.UserID,
		subscription.Name,
		subscription.Price,
		subscription.Plan,
		subscription.NextPaymentAt,
	)

	if err := row.Scan(
		&subscription.ID,
		&subscription.CreatedAt,
	); err != nil {
		return err
	}

	return nil
}
