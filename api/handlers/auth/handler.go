package auth

import (
	"net/http"

	authSvc "github.com/deuna/payment/domain/services/auth"
	"github.com/labstack/echo/v4"
)

type auth struct {
	authService authSvc.Authorization
}

// NewAuthHandler creates a new auth handler
// @Summary Create Token
// @Description Generates a token
// @Tags Authorization
// @Param requestBody body CreateTokenRequest true "Generate token"
// @Success 200 {object} CreateTokenResponse
// @Router /v1/auth/token [post]
func (a auth) CreateToken(c echo.Context) error {
	req := CreateTokenRequest{}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "invalid request")
	}

	ctx := c.Request().Context()

	token, err := a.authService.CreateToken(ctx, req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, CreateTokenResponse{Token: *token})

}
