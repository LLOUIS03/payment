package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/deuna/payment/test/mocks"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateToken_HappyPath(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTM2OTA2MTgsImVtYWlsIjoiamtpcmtuZXNzMEBXZW5keS5jb20iLCJtZXJjaGFudF9pZCI6ImIxYjliMWI5LTFiOWItMWI5Yi0xYjliLTFiOWIxYjliMWI5YiJ9.enEZmSheCIeak1ctaYbq-wsQA6R79nai5Kz-vU2NLwQ"

	ctrl := gomock.NewController(t)

	e := echo.New()
	authSvc := mocks.NewMockAuthorization(ctrl)
	authSvc.EXPECT().CreateToken(gomock.Any(), gomock.Any(), gomock.Any()).Return(&token, nil)

	req := CreateTokenRequest{
		Email:    "Luis@gmail.com",
		Password: "123456",
	}

	jsonBody, err := json.Marshal(req)
	assert.NoError(t, err)

	httpReq := httptest.NewRequest(http.MethodPost, route, bytes.NewReader(jsonBody))
	httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	SetupHandler(e, authSvc)

	// act
	respHttp := httptest.NewRecorder()
	e.ServeHTTP(respHttp, httpReq)

	// assert
	assert.Equal(t, http.StatusOK, respHttp.Code)

	respBody, err := io.ReadAll(respHttp.Body)
	assert.NoError(t, err)

	resp := CreateTokenResponse{}
	err = json.Unmarshal(respBody, &resp)

	assert.NoError(t, err)
	assert.Equal(t, token, resp.Token)
}

func TestCreateToken_BadReq(t *testing.T) {
	ctrl := gomock.NewController(t)

	e := echo.New()
	authSvc := mocks.NewMockAuthorization(ctrl)

	httpReq := httptest.NewRequest(http.MethodPost, route, bytes.NewReader([]byte("bad")))

	SetupHandler(e, authSvc)

	// act
	respHttp := httptest.NewRecorder()
	e.ServeHTTP(respHttp, httpReq)

	// assert
	assert.Equal(t, http.StatusBadRequest, respHttp.Code)
}

func TestCreateToken_InternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	expetedErr := errors.New("test error")

	e := echo.New()
	authSvc := mocks.NewMockAuthorization(ctrl)
	authSvc.EXPECT().CreateToken(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, expetedErr)

	req := CreateTokenRequest{
		Email:    "Luis@gmail.com",
		Password: "123456",
	}

	jsonBody, err := json.Marshal(req)
	assert.NoError(t, err)

	httpReq := httptest.NewRequest(http.MethodPost, route, bytes.NewReader(jsonBody))
	httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	SetupHandler(e, authSvc)

	// act
	respHttp := httptest.NewRecorder()
	e.ServeHTTP(respHttp, httpReq)

	// assert
	assert.Equal(t, http.StatusInternalServerError, respHttp.Code)

}
