package api

import (
	"github.com/deuna/payment/domain/services/auth"
	"github.com/deuna/payment/domain/services/transaction"
	"github.com/labstack/echo/v4"
)

type API struct {
	server   *echo.Echo
	services Services
}

// Services is a struct that contains all the services
type Services struct {
	TxSvc   transaction.Transaction
	AuthSvc auth.Authorization
}
