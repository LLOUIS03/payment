package handlers

import (
	"net/http"

	"github.com/deuna/payment/domain/domainerrors"
	"github.com/labstack/echo/v4"
)

func HandleError(c echo.Context, err error) error {
	switch v := err.(type) {
	case domainerrors.CredentialsExpiredError:
		return c.JSON(http.StatusInternalServerError, v.Message())
	case domainerrors.NoTransactionError:
		return c.JSON(http.StatusInternalServerError, v.Message())
	case domainerrors.UnauthorizedTranError:
		return c.JSON(http.StatusInternalServerError, v.Message())
	case domainerrors.InternalServerError:
		return c.JSON(http.StatusInternalServerError, v.Message())
	case domainerrors.InvalidCredentialsError:
		return c.JSON(http.StatusUnauthorized, v.Message())
	case *domainerrors.ValidationError:
		return c.JSON(http.StatusBadRequest, v.Error())
	}

	return c.JSON(http.StatusInternalServerError, err.Error())
}
