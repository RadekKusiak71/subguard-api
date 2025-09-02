package subscriptions

import (
	"net/http"

	"github.com/RadekKusiak71/subguard-api/internal/users"
	"github.com/RadekKusiak71/subguard-api/internal/utils"
)

type Handler struct {
	store SubscriptionStore
}

func NewHandler(subStore SubscriptionStore) *Handler {
	return &Handler{
		store: subStore,
	}
}

func (h *Handler) ListSubscriptions(w http.ResponseWriter, r *http.Request) error {
	userID := r.Context().Value(users.UserContextKey).(int)

	subscriptions, err := h.store.List(userID)

	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, subscriptions)
}

func (h *Handler) CreateSubscription(w http.ResponseWriter, r *http.Request) error {
	userID := r.Context().Value(users.UserContextKey).(int)
	var subscriptionPayload CreateSubscription

	if err := utils.ParseJSON(r, &subscriptionPayload); err != nil {
		return utils.InvalidJSON()
	}

	if validationErrors := subscriptionPayload.Validate(); len(validationErrors) > 0 {
		return utils.InvalidRequest(validationErrors)
	}

	_, err := h.store.GetByName(userID, subscriptionPayload.Name)

	if err == nil {
		return SubscriptionExists()
	}

	newSubscription := Subscription{
		UserID:        userID,
		Name:          subscriptionPayload.Name,
		Plan:          subscriptionPayload.Plan,
		Price:         subscriptionPayload.Price,
		NextPaymentAt: subscriptionPayload.NextPaymentAt,
	}

	if err := h.store.Create(&newSubscription); err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusCreated, newSubscription)
}
