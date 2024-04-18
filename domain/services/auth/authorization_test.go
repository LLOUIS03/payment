package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/deuna/payment/domain/domainerrors"
	"github.com/deuna/payment/infraestructure/db/repos"
	"github.com/deuna/payment/test/mocks"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateToken_HappyPath(t *testing.T) {
	// Arrange
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	queries := mocks.NewMockQuerier(ctrl)

	userRow := repos.GetUserRow{
		ID:       uuid.New(),
		Password: "password",
	}

	queries.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(userRow, nil)

	// Act
	auth := NewService(queries)
	tokenStr, err := auth.CreateToken(ctx, "luis.louis@gmail.com", "password")

	token, errSec := jwt.Parse(*tokenStr, func(token *jwt.Token) (any, error) {
		return Secret, nil
	})

	//Assert
	assert.NoError(t, err)
	assert.NotNil(t, token)

	assert.NoError(t, errSec)

}

func TestCreateToken_NoUser(t *testing.T) {
	// Arrange
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	queries := mocks.NewMockQuerier(ctrl)

	userRow := repos.GetUserRow{
		ID:       uuid.New(),
		Password: "password",
	}

	queries.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(userRow, errors.New("no user"))

	// Act
	auth := NewService(queries)
	toke, err := auth.CreateToken(ctx, "luis.louis@gmail.com", "password")

	//Assert
	assert.Nil(t, toke)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "no user")
	assert.IsType(t, domainerrors.InvalidCredentialsError{}, err)
}

func TestCreateToken_DiffPass(t *testing.T) {
	// Arrange
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	queries := mocks.NewMockQuerier(ctrl)

	userRow := repos.GetUserRow{
		ID:       uuid.New(),
		Password: "password",
	}

	queries.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(userRow, nil)

	// Act
	auth := NewService(queries)
	toke, err := auth.CreateToken(ctx, "luis.louis@gmail.com", "1234")

	//Assert
	assert.Nil(t, toke)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "invalid credentials")
	assert.IsType(t, domainerrors.InvalidCredentialsError{}, err)
}
