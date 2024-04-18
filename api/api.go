package api

import (
	"context"
	"net/http"
	"os"
	"syscall"

	_ "github.com/deuna/payment/api/docs"
	"github.com/deuna/payment/api/handlers/auth"
	"github.com/deuna/payment/api/handlers/transaction"
	"github.com/deuna/payment/config"
	authsvc "github.com/deuna/payment/domain/services/auth"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// NewAPI creates a new API
func NewAPI(svc Services) *API {
	return &API{
		server:   echo.New(),
		services: svc,
	}
}

// Start starts the API
func (a *API) Start(ctx context.Context, quit chan os.Signal) {
	cfg := config.Instance()
	a.setupAuthRouter()

	group := a.setupMiddleware()

	a.setupRoutes(group)

	defer func() {
		<-ctx.Done()
		quit <- syscall.SIGTERM
	}()

	err := a.server.Start(cfg.Server.Port)
	if err != nil {
		quit <- syscall.SIGTERM
	}
}

// setupMiddleware sets up the middleware
func (a *API) setupMiddleware() *echo.Group {
	cfg := config.Instance()

	a.server.Use(middleware.Logger())
	a.server.Use(middleware.Recover())
	a.server.Use(middleware.CORS())
	a.server.Use(middleware.Secure())
	a.server.Use(middleware.RequestID())

	a.server.Server.ReadTimeout = cfg.Server.Timeout
	a.server.Logger.SetLevel(lvl(cfg.Server.LogLevel))

	group := a.server.Group(cfg.Server.BaseAddr)
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(authsvc.JwtCustomClaims)
		},
		SigningKey: authsvc.Secret,
	}

	group.Use(echojwt.WithConfig(config))

	return group
}

// Stop stops the API
func (a *API) Stop(ctx context.Context) error {
	if a.server != nil {
		log.Info("shutting down the server...")
		if err := a.server.Shutdown(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *API) setupAuthRouter() {
	auth.SetupHandler(a.server, a.services.AuthSvc)
	a.server.GET("/swagger/*", echoSwagger.WrapHandler)
	a.server.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusPermanentRedirect, "/swagger/index.html")
	})
}

func (a *API) setupRoutes(group *echo.Group) {
	transaction.SetupHandler(group, a.services.TxSvc)
}
