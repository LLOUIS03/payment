package auth

import (
	authSvc "github.com/deuna/payment/domain/services/auth"
	"github.com/labstack/echo/v4"
)

const route = "/v1/auth/token"

func SetupHandler(e *echo.Echo, a authSvc.Authorization) {
	handler := handler{authService: a}
	e.POST(route, handler.CreateToken)
}
