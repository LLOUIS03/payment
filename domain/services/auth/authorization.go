package auth

import (
	"context"
	"time"

	"github.com/deuna/payment/domain/domainerrors"
	"github.com/deuna/payment/infraestructure/db/repos"
	"github.com/golang-jwt/jwt/v5"
)

// Authenticator is the interface that wraps the authenticator methods
type Authorization interface {
	CreateToken(context.Context, string, string) (*string, error)
}

type authenticator struct {
	queries repos.Querier
}

// NewService creates a new authenticator
func NewService(queries repos.Querier) Authorization {
	return &authenticator{
		queries: queries,
	}
}

// CreateToken creates a token for the user
func (a *authenticator) CreateToken(ctx context.Context, email, password string) (*string, error) {
	userRow, err := a.queries.GetUser(ctx, email)
	if err != nil {
		return nil, domainerrors.NewInvalidCredentialsError(err)
	}

	if userRow.Password != password {
		return nil, domainerrors.NewInvalidCredentialsError(nil)
	}

	claims := JwtCustomClaims{
		Email: email,
		ID:    userRow.ID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(Secret)
	if err != nil {
		return nil, domainerrors.NewInvalidCredentialsError(err)
	}

	return &tokenString, nil
}
