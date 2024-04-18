package transaction

import (
	"github.com/deuna/payment/domain/services/transaction"
	"github.com/labstack/echo/v4"
)

const (
	placeroute  = "/place"
	refundroute = "/refund"
	listtxroute = "/listtx"
)

func SetupHandler(e *echo.Group, svc transaction.Transaction) {
	handler := handler{transactionService: svc}
	e.POST(placeroute, handler.Place)
	e.POST(refundroute, handler.Refund)
	e.GET(listtxroute, handler.ListTx)
}
