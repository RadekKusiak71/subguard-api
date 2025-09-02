package subscriptions

import "time"

type SubscriptionPlan string

const (
	Monthly SubscriptionPlan = "monthly"
	Yearly  SubscriptionPlan = "yearly"
)

var validPlans = map[SubscriptionPlan]string{
	Monthly: "monthly",
	Yearly:  "yearly",
}

func (p SubscriptionPlan) String() string {
	if name, ok := validPlans[p]; ok {
		return name
	}
	return "Unknown Plan"
}

type SubscriptionStore interface {
	List(userID int) ([]Subscription, error)
	GetByName(userID int, name string) (*Subscription, error)
	Get(userID, subscriptionID int) (*Subscription, error)
	Update(subscription *Subscription) error
	Create(subscription *Subscription) error
	Delete(userID, subscriptionID int) error
	GetExpiringSoon() ([]Subscription, error)
	UpdateNextPaymentBatch(subs []Subscription) error
}

type Subscription struct {
	ID            int              `json:"id"`
	UserID        int              `json:"user_id"`
	Name          string           `json:"name"`
	Price         float64          `json:"price"`
	Plan          SubscriptionPlan `json:"plan"`
	NextPaymentAt time.Time        `json:"next_payment_at"`
	CreatedAt     time.Time        `json:"created_at"`
}

type CreateSubscription struct {
	Name          string           `json:"name"`
	Price         float64          `json:"price"`
	Plan          SubscriptionPlan `json:"plan"`
	NextPaymentAt time.Time        `json:"next_payment_at"`
}
