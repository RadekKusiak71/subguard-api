package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/RadekKusiak71/subguard-api/internal/authentication"
	"github.com/RadekKusiak71/subguard-api/internal/users"
	"github.com/RadekKusiak71/subguard-api/internal/utils"
	"github.com/golang-jwt/jwt/v5"
)

func extractToken(r *http.Request) (*jwt.Token, error) {
	authString := r.Header.Get("Authorization")
	if authString == "" {
		return nil, authentication.MissingToken()
	}

	parts := strings.SplitN(authString, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, authentication.InvalidAuthorizationHeader()
	}

	tokenString := strings.TrimSpace(parts[1])
	if tokenString == "" {
		return nil, authentication.InvalidAuthorizationHeader()
	}

	token, err := authentication.ValidateJWT(tokenString)
	if err != nil {
		return nil, authentication.InvalidToken()
	}

	return token, nil
}

func AuthMiddleware(next utils.APIHandler, userStore users.UserStore) utils.APIHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		token, err := extractToken(r)
		if err != nil {
			return err
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return authentication.InvalidToken()
		}

		userIDStr, ok := claims["userID"].(string)
		if !ok {
			return authentication.InvalidToken()
		}

		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			return authentication.InvalidToken()
		}

		user, err := userStore.Get(userID)
		if err != nil {
			if errors.Is(err, users.ErrUserNotFound) {
				return authentication.InvalidToken()
			}
			return err
		}

		if !user.IsActive {
			return users.AccountNotVerified()
		}

		ctx := context.WithValue(
			r.Context(),
			users.UserContextKey,
			user.ID,
		)

		r = r.WithContext(ctx)

		return next(w, r)
	}
}
