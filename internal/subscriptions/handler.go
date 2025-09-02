package subscriptions

import (
	"errors"
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

func (h *Handler) GetSubscription(w http.ResponseWriter, r *http.Request) error {
	userID := r.Context().Value(users.UserContextKey).(int)
	subscriptionID, err := utils.ReadParamFromPathAsInt(r, "subscriptionID")
	if err != nil {
		return err
	}

	subscription, err := h.store.Get(userID, subscriptionID)

	if err != nil {
		if errors.Is(err, ErrSubscriptionNotFound) {
			return SubscriptionNotFound()
		}
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, subscription)
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

func (h *Handler) UpdateSubscription(w http.ResponseWriter, r *http.Request) error {
	userID := r.Context().Value(users.UserContextKey).(int)
	subscriptionID, err := utils.ReadParamFromPathAsInt(r, "subscriptionID")

	if err != nil {
		return err
	}

	var subscriptionPayload CreateSubscription
	if err := utils.ParseJSON(r, &subscriptionPayload); err != nil {
		return utils.InvalidJSON()
	}

	if validationErrors := subscriptionPayload.Validate(); len(validationErrors) > 0 {
		return utils.InvalidRequest(validationErrors)
	}

	updatedSubscription := Subscription{
		ID:            subscriptionID,
		UserID:        userID,
		Name:          subscriptionPayload.Name,
		Price:         subscriptionPayload.Price,
		Plan:          subscriptionPayload.Plan,
		NextPaymentAt: subscriptionPayload.NextPaymentAt,
	}

	if err := h.store.Update(&updatedSubscription); err != nil {
		if errors.Is(err, ErrSubscriptionNotFound) {
			return SubscriptionNotFound()
		}
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, updatedSubscription)
}

func (h *Handler) DeleteSubscription(w http.ResponseWriter, r *http.Request) error {
	userID := r.Context().Value(users.UserContextKey).(int)
	subscriptionID, err := utils.ReadParamFromPathAsInt(r, "subscriptionID")

	if err != nil {
		return err
	}

	if err := h.store.Delete(userID, subscriptionID); err != nil {
		if errors.Is(err, ErrSubscriptionNotFound) {
			return SubscriptionNotFound()
		}
		return err
	}

	return utils.WriteJSON(w, http.StatusNoContent, nil)
}
