package auth

import (
	"net/http"

	"github.com/deuna/payment/api/handlers"
	authSvc "github.com/deuna/payment/domain/services/auth"
	"github.com/labstack/echo/v4"
)

type handler struct {
	authService authSvc.Authorization
}

// NewAuthHandler creates a new auth handler
// @Summary Create Token
// @Description Generates a token
// @Tags Authorization
// @Param requestBody body CreateTokenRequest true "Generate token"
// @Success 200 {object} CreateTokenResponse
// @Router /v1/auth/token [post]
func (a handler) CreateToken(c echo.Context) error {
	req := CreateTokenRequest{}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "invalid request")
	}

	ctx := c.Request().Context()

	token, err := a.authService.CreateToken(ctx, req.Email, req.Password)
	if err != nil {
		return handlers.HandleError(c, err)
	}

	return c.JSON(http.StatusOK, CreateTokenResponse{Token: *token})

}
