package authentication

import (
	"errors"
	"net/http"

	"github.com/RadekKusiak71/subguard-api/internal/users"
	"github.com/RadekKusiak71/subguard-api/internal/utils"
)

type Handler struct {
	userStore users.UserStore
}

func NewHandler(userStore users.UserStore) *Handler {
	return &Handler{
		userStore: userStore,
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) error {
	var registerPayload RegisterUser

	if err := utils.ParseJSON(r, &registerPayload); err != nil {
		return utils.InvalidJSON()
	}

	if validationErrors := registerPayload.Validate(); len(validationErrors) > 0 {
		return utils.InvalidRequest(validationErrors)
	}

	_, err := h.userStore.GetByEmail(registerPayload.Email)

	if err == nil {
		return users.UserAlreadyExist()
	}

	hashedPassword, err := HashPassword(registerPayload.Password)

	if err != nil {
		return err
	}

	user := users.User{
		Email:    registerPayload.Email,
		Password: hashedPassword,
	}

	if err := h.userStore.Create(&user); err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, user)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) error {
	var loginPayload LoginPayload

	if err := utils.ParseJSON(r, &loginPayload); err != nil {
		return utils.InvalidJSON()
	}

	user, err := h.userStore.GetByEmail(loginPayload.Email)

	if err != nil {
		if errors.Is(err, users.ErrUserNotFound) {
			return users.UserDoesNotExist()
		}
		return err
	}

	if !ComparePasswords(loginPayload.Password, user.Password) {
		return InvalidCredentials()
	}

	token, err := CreateToken(user.ID)

	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, TokenResponse{
		AccessToken: token,
	})
}
