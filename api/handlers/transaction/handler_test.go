package transaction

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	authSvc "github.com/deuna/payment/domain/services/auth"
	"github.com/deuna/payment/infraestructure/db/repos"
	"github.com/deuna/payment/test/mocks"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestPlace_HappyPath(t *testing.T) {
	// arrange
	tokenStr := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTM2OTA2MTgsImVtYWlsIjoiamtpcmtuZXNzMEBXZW5keS5jb20iLCJtZXJjaGFudF9pZCI6ImIxYjliMWI5LTFiOWItMWI5Yi0xYjliLTFiOWIxYjliMWI5YiJ9.enEZmSheCIeak1ctaYbq-wsQA6R79nai5Kz-vU2NLwQ"
	tranUUID := uuid.New()
	e := echo.New()

	ctrl := gomock.NewController(t)

	tranSvc := mocks.NewMockTransaction(ctrl)
	tranSvc.EXPECT().Place(gomock.Any(), gomock.Any(), gomock.Any()).Return(&tranUUID, nil)

	req := PlaceRequest{
		Amount: 100,
	}

	jsonBody, err := json.Marshal(req)
	assert.NoError(t, err)

	token, err := jwt.ParseWithClaims(tokenStr, &authSvc.JwtCustomClaims{}, func(token *jwt.Token) (any, error) {
		return authSvc.Secret, nil
	})
	assert.NoError(t, err)

	httpReq := httptest.NewRequest(http.MethodPost, placeroute, bytes.NewReader(jsonBody))
	httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	respHttp := httptest.NewRecorder()
	c := e.NewContext(httpReq, respHttp)
	c.Set(userIDKey, token)

	hdl := handler{transactionService: tranSvc}

	// act
	err = hdl.Place(c)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, respHttp.Code)

	respBody, err := io.ReadAll(respHttp.Body)
	assert.NoError(t, err)
	assert.Contains(t, string(respBody), tranUUID.String())
}

func TestPlace_BadRequest(t *testing.T) {
	// arrange
	tokenStr := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTM2OTA2MTgsImVtYWlsIjoiamtpcmtuZXNzMEBXZW5keS5jb20iLCJtZXJjaGFudF9pZCI6ImIxYjliMWI5LTFiOWItMWI5Yi0xYjliLTFiOWIxYjliMWI5YiJ9.enEZmSheCIeak1ctaYbq-wsQA6R79nai5Kz-vU2NLwQ"
	e := echo.New()

	ctrl := gomock.NewController(t)

	tranSvc := mocks.NewMockTransaction(ctrl)

	token, err := jwt.ParseWithClaims(tokenStr, &authSvc.JwtCustomClaims{}, func(token *jwt.Token) (any, error) {
		return authSvc.Secret, nil
	})
	assert.NoError(t, err)

	httpReq := httptest.NewRequest(http.MethodPost, placeroute, bytes.NewReader([]byte("test")))
	httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	respHttp := httptest.NewRecorder()
	c := e.NewContext(httpReq, respHttp)
	c.Set(userIDKey, token)

	hdl := handler{transactionService: tranSvc}

	// act
	err = hdl.Place(c)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, respHttp.Code)
}

func TestPlace_InternalServerError(t *testing.T) {
	// arrange
	tokenStr := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTM2OTA2MTgsImVtYWlsIjoiamtpcmtuZXNzMEBXZW5keS5jb20iLCJtZXJjaGFudF9pZCI6ImIxYjliMWI5LTFiOWItMWI5Yi0xYjliLTFiOWIxYjliMWI5YiJ9.enEZmSheCIeak1ctaYbq-wsQA6R79nai5Kz-vU2NLwQ"
	e := echo.New()

	expectedErr := errors.New("test error")

	ctrl := gomock.NewController(t)

	tranSvc := mocks.NewMockTransaction(ctrl)
	tranSvc.EXPECT().Place(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, expectedErr)

	req := PlaceRequest{
		Amount: 100,
	}

	jsonBody, err := json.Marshal(req)
	assert.NoError(t, err)

	token, err := jwt.ParseWithClaims(tokenStr, &authSvc.JwtCustomClaims{}, func(token *jwt.Token) (any, error) {
		return authSvc.Secret, nil
	})
	assert.NoError(t, err)

	httpReq := httptest.NewRequest(http.MethodPost, placeroute, bytes.NewReader(jsonBody))
	httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	respHttp := httptest.NewRecorder()
	c := e.NewContext(httpReq, respHttp)
	c.Set(userIDKey, token)

	hdl := handler{transactionService: tranSvc}

	// act
	err = hdl.Place(c)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, respHttp.Code)
}

func TestRefund_HappyPath(t *testing.T) {
	// arrange
	tokenStr := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTM2OTA2MTgsImVtYWlsIjoiamtpcmtuZXNzMEBXZW5keS5jb20iLCJtZXJjaGFudF9pZCI6ImIxYjliMWI5LTFiOWItMWI5Yi0xYjliLTFiOWIxYjliMWI5YiJ9.enEZmSheCIeak1ctaYbq-wsQA6R79nai5Kz-vU2NLwQ"
	tranUUID := uuid.New()
	e := echo.New()

	ctrl := gomock.NewController(t)

	tranSvc := mocks.NewMockTransaction(ctrl)
	tranSvc.EXPECT().Refund(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	req := RefundRequest{
		ID: tranUUID.String(),
	}

	jsonBody, err := json.Marshal(req)
	assert.NoError(t, err)

	token, err := jwt.ParseWithClaims(tokenStr, &authSvc.JwtCustomClaims{}, func(token *jwt.Token) (any, error) {
		return authSvc.Secret, nil
	})
	assert.NoError(t, err)

	httpReq := httptest.NewRequest(http.MethodPost, placeroute, bytes.NewReader(jsonBody))
	httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	respHttp := httptest.NewRecorder()
	c := e.NewContext(httpReq, respHttp)
	c.Set(userIDKey, token)

	hdl := handler{transactionService: tranSvc}

	// act
	err = hdl.Refund(c)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, respHttp.Code)
}

func TestRefund_BadRequest(t *testing.T) {
	// arrange
	tokenStr := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTM2OTA2MTgsImVtYWlsIjoiamtpcmtuZXNzMEBXZW5keS5jb20iLCJtZXJjaGFudF9pZCI6ImIxYjliMWI5LTFiOWItMWI5Yi0xYjliLTFiOWIxYjliMWI5YiJ9.enEZmSheCIeak1ctaYbq-wsQA6R79nai5Kz-vU2NLwQ"
	e := echo.New()

	ctrl := gomock.NewController(t)

	tranSvc := mocks.NewMockTransaction(ctrl)

	token, err := jwt.ParseWithClaims(tokenStr, &authSvc.JwtCustomClaims{}, func(token *jwt.Token) (any, error) {
		return authSvc.Secret, nil
	})
	assert.NoError(t, err)

	httpReq := httptest.NewRequest(http.MethodPost, placeroute, bytes.NewReader([]byte("badreq")))
	httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	respHttp := httptest.NewRecorder()
	c := e.NewContext(httpReq, respHttp)
	c.Set(userIDKey, token)

	hdl := handler{transactionService: tranSvc}

	// act
	err = hdl.Refund(c)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, respHttp.Code)
}

func TestRefund_InternalError(t *testing.T) {
	// arrange
	tokenStr := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTM2OTA2MTgsImVtYWlsIjoiamtpcmtuZXNzMEBXZW5keS5jb20iLCJtZXJjaGFudF9pZCI6ImIxYjliMWI5LTFiOWItMWI5Yi0xYjliLTFiOWIxYjliMWI5YiJ9.enEZmSheCIeak1ctaYbq-wsQA6R79nai5Kz-vU2NLwQ"
	tranUUID := uuid.New()
	e := echo.New()

	expectedErr := errors.New("test error")

	ctrl := gomock.NewController(t)

	tranSvc := mocks.NewMockTransaction(ctrl)
	tranSvc.EXPECT().Refund(gomock.Any(), gomock.Any(), gomock.Any()).Return(expectedErr)

	req := RefundRequest{
		ID: tranUUID.String(),
	}

	jsonBody, err := json.Marshal(req)
	assert.NoError(t, err)

	token, err := jwt.ParseWithClaims(tokenStr, &authSvc.JwtCustomClaims{}, func(token *jwt.Token) (any, error) {
		return authSvc.Secret, nil
	})
	assert.NoError(t, err)

	httpReq := httptest.NewRequest(http.MethodPost, placeroute, bytes.NewReader(jsonBody))
	httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	respHttp := httptest.NewRecorder()
	c := e.NewContext(httpReq, respHttp)
	c.Set(userIDKey, token)

	hdl := handler{transactionService: tranSvc}

	// act
	err = hdl.Refund(c)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, respHttp.Code)
}

func TestListTx_HappyPath(t *testing.T) {
	// arrange
	tokenStr := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTM2OTA2MTgsImVtYWlsIjoiamtpcmtuZXNzMEBXZW5keS5jb20iLCJtZXJjaGFudF9pZCI6ImIxYjliMWI5LTFiOWItMWI5Yi0xYjliLTFiOWIxYjliMWI5YiJ9.enEZmSheCIeak1ctaYbq-wsQA6R79nai5Kz-vU2NLwQ"
	e := echo.New()

	ctrl := gomock.NewController(t)

	transactions := []repos.Tx{
		{
			ID:         uuid.New(),
			MerchantID: uuid.New(),
			Amount:     100.00,
		}, {
			ID:         uuid.New(),
			MerchantID: uuid.New(),
			Amount:     100.00,
		},
	}

	tranSvc := mocks.NewMockTransaction(ctrl)
	tranSvc.EXPECT().ListTx(gomock.Any(), gomock.Any()).Return(transactions, nil)

	token, err := jwt.ParseWithClaims(tokenStr, &authSvc.JwtCustomClaims{}, func(token *jwt.Token) (any, error) {
		return authSvc.Secret, nil
	})
	assert.NoError(t, err)

	httpReq := httptest.NewRequest(http.MethodPost, placeroute, nil)
	httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	respHttp := httptest.NewRecorder()
	c := e.NewContext(httpReq, respHttp)
	c.Set(userIDKey, token)

	hdl := handler{transactionService: tranSvc}

	// act
	err = hdl.ListTx(c)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, respHttp.Code)

	respBody, err := io.ReadAll(respHttp.Body)
	assert.NoError(t, err)

	resp := ListTxResponse{}
	err = json.Unmarshal(respBody, &resp)
	assert.NoError(t, err)
	assert.Len(t, resp.Txs, 2)
}

func TestListTx_InternalServer(t *testing.T) {
	// arrange
	tokenStr := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTM2OTA2MTgsImVtYWlsIjoiamtpcmtuZXNzMEBXZW5keS5jb20iLCJtZXJjaGFudF9pZCI6ImIxYjliMWI5LTFiOWItMWI5Yi0xYjliLTFiOWIxYjliMWI5YiJ9.enEZmSheCIeak1ctaYbq-wsQA6R79nai5Kz-vU2NLwQ"
	e := echo.New()

	ctrl := gomock.NewController(t)

	expectedErr := errors.New("test error")

	tranSvc := mocks.NewMockTransaction(ctrl)
	tranSvc.EXPECT().ListTx(gomock.Any(), gomock.Any()).Return(nil, expectedErr)

	token, err := jwt.ParseWithClaims(tokenStr, &authSvc.JwtCustomClaims{}, func(token *jwt.Token) (any, error) {
		return authSvc.Secret, nil
	})
	assert.NoError(t, err)

	httpReq := httptest.NewRequest(http.MethodPost, placeroute, nil)
	httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	respHttp := httptest.NewRecorder()
	c := e.NewContext(httpReq, respHttp)
	c.Set(userIDKey, token)

	hdl := handler{transactionService: tranSvc}

	// act
	err = hdl.ListTx(c)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, respHttp.Code)

}
