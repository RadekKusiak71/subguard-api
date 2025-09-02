package subscriptions

import (
	"database/sql"
	"errors"
	"time"
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

func (s *Store) Get(userID, subscriptionID int) (*Subscription, error) {
	row := s.db.QueryRow(
		"SELECT * FROM subscriptions WHERE user_id = $1 AND id = $2",
		userID,
		subscriptionID,
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

func (s *Store) Update(subscription *Subscription) error {
	result, err := s.db.Exec(
		`
        UPDATE subscriptions
        SET name = $1, price = $2, plan = $3, next_payment_at = $4
        WHERE user_id = $5 AND id = $6
        `,
		subscription.Name,
		subscription.Price,
		subscription.Plan,
		subscription.NextPaymentAt,
		subscription.UserID,
		subscription.ID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrSubscriptionNotFound
	}

	return nil
}

func (s *Store) UpdateNextPaymentBatch(subs []Subscription) error {
	ids := make([]int, len(subs))
	for i, sub := range subs {
		ids[i] = sub.ID
	}

	query := `
        UPDATE subscriptions
        SET next_payment_at = CASE plan
            WHEN 'monthly' THEN next_payment_at + INTERVAL '1 month'
            WHEN 'yearly'  THEN next_payment_at + INTERVAL '1 year'
            ELSE next_payment_at
        END
        WHERE id = ANY($1)
    `
	_, err := s.db.Exec(query, ids)
	return err
}

func (s *Store) GetExpiringSoon() ([]Subscription, error) {
	now := time.Now()
	tomorrow := now.Add(24 * time.Hour)

	startOfTomorrow := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, tomorrow.Location())
	endOfTomorrow := startOfTomorrow.Add(24*time.Hour - time.Nanosecond)

	rows, err := s.db.Query(
		`SELECT id, user_id, name, price, plan, next_payment_at, created_at
		FROM subscriptions
		WHERE next_payment_at BETWEEN $1 AND $2`,
		startOfTomorrow,
		endOfTomorrow,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

	return subscriptions, nil
}

func (s *Store) Delete(userID, subscriptionID int) error {
	result, err := s.db.Exec(
		`DELETE FROM subscriptions WHERE user_id = $1 AND id = $2`,
		userID, subscriptionID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrSubscriptionNotFound
	}

	return nil
}
