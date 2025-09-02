package subscriptions

import (
	"fmt"
	"strings"
	"time"
)

func (p SubscriptionPlan) IsValid() bool {
	_, valid := validPlans[p]
	return valid
}

func (s *CreateSubscription) Validate() map[string]string {
	s.Name = strings.Trim(s.Name, " ")
	errors := make(map[string]string)

	if s.Name == "" {
		errors["name"] = "name cannot be empty"
	} else if len(s.Name) < 3 || len(s.Name) > 255 {
		errors["name"] = "name must be between 3 and 255 characters"
	}

	if s.Price <= 0 {
		errors["price"] = "price must be greater than 0"
	}

	if s.NextPaymentAt.Before(time.Now()) {
		errors["next_payment_at"] = "next payment at must be a future date"
	}

	if !s.Plan.IsValid() {
		errors["plan"] = fmt.Sprintf("plan must be one of %v", validPlans)
	}

	return errors
}
