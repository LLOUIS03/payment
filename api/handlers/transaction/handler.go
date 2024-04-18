package transaction

import (
	"net/http"

	"github.com/deuna/payment/api/handlers"
	authsvc "github.com/deuna/payment/domain/services/auth"
	transactionSvc "github.com/deuna/payment/domain/services/transaction"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

const (
	userIDKey = "user"
)

type handler struct {
	transactionService transactionSvc.Transaction
}

// @Summary Place
// @Description Place a transaction
// @Tags Transaction
// @Param Authorization header string false "The JWT token. Example: Bearer {token}"
// @Param requestBody body PlaceRequest true "Generate token"
// @Success 200 {string} string	"transaction id"
// @Router /v1/api/place [post].
func (h *handler) Place(c echo.Context) error {
	req := PlaceRequest{}
	ctx := c.Request().Context()

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "invalid request")
	}

	user := c.Get(userIDKey).(*jwt.Token)
	claims := user.Claims.(*authsvc.JwtCustomClaims)

	id, err := h.transactionService.Place(ctx, req.Amount, claims.ID)
	if err != nil {
		return handlers.HandleError(c, err)
	}

	return c.JSON(http.StatusOK, id)
}

// @Summary Refund
// @Description Refund a transaction
// @Tags Transaction
// @Param Authorization header string false "The JWT token. Example: Bearer {token}"
// @Param requestBody body RefundRequest true "Generate token"
// @Success 200 {object} nil
// @Router /v1/api/refund [post].
func (h *handler) Refund(c echo.Context) error {
	req := RefundRequest{}
	ctx := c.Request().Context()

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "invalid request")
	}

	user := c.Get(userIDKey).(*jwt.Token)
	claims := user.Claims.(*authsvc.JwtCustomClaims)

	err := h.transactionService.Refund(ctx, req.ID, claims.ID)
	if err != nil {
		return handlers.HandleError(c, err)
	}

	return c.JSON(http.StatusOK, nil)
}

// @Summary ListTx
// @Description List all the transactions
// @Tags Transaction
// @Param Authorization header string false "The JWT token. Example: Bearer {token}"
// @Success 200 {object} ListTxResponse
// @Router /v1/api/listtx [get].
func (h *handler) ListTx(c echo.Context) error {
	ctx := c.Request().Context()

	user := c.Get(userIDKey).(*jwt.Token)
	claims := user.Claims.(*authsvc.JwtCustomClaims)

	txs, err := h.transactionService.ListTx(ctx, claims.ID)
	if err != nil {
		return handlers.HandleError(c, err)
	}

	return c.JSON(http.StatusOK, ListTxResponse{Txs: txs})
}
