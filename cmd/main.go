package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"emperror.dev/errors"
	"github.com/deuna/payment/api"
	"github.com/deuna/payment/config"
	"github.com/deuna/payment/domain/clients"
	"github.com/deuna/payment/domain/services/auth"
	"github.com/deuna/payment/domain/services/transaction"
	"github.com/deuna/payment/infraestructure/db/repos"
	"github.com/deuna/payment/infraestructure/db/setupdb"
	"github.com/neko-neko/echo-logrus/v2/log"
)

// @title Payment API
// @version 1.0
// @description Payment API
func main() {
	// Parse the environment flag
	envF := flag.String("env", "local", "Environment")
	flag.Parse()

	// Read the config file
	cfg, err := config.ReadConfig(*envF, ".")
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Set up the database connection
	db, err := setupDB(ctx, cfg)
	if err != nil {
		log.Fatalf("error setting up database: %v", err)
	}

	// Set up services
	services := setupServices(db)

	// Set up the API
	httpAPI := api.NewAPI(services)

	// Create quit channel and start shutdown goroutine
	quit := make(chan os.Signal, 1)
	go shutdown(ctx, httpAPI, cancel, db, quit)

	// Start the API
	httpAPI.Start(ctx, quit)

}

// shutdown shuts down the application
func shutdown(ctx context.Context, httpAPI *api.API,
	cancel context.CancelFunc, db setupdb.Database,
	quit chan os.Signal) {

	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Info("shutting down the http server...")
	if err := httpAPI.Stop(ctx); err != nil {
		log.Errorf("error stopping the http server: %v", err)
	}

	log.Info("closing the database connection...")
	db.CloseFunc()

	log.Info("canceling the main context...")
	cancel()

	log.Info("exiting the application...")
	os.Exit(1)
}

// setupServices sets up the services
func setupServices(db setupdb.Database) api.Services {
	querier := repos.New(db.DB())
	bankClient := clients.NewBank()
	txSvc := transaction.NewService(querier, bankClient)
	authSvc := auth.NewService(querier)

	return api.Services{
		TxSvc:   txSvc,
		AuthSvc: authSvc,
	}
}

// setupDB sets up the database connection
func setupDB(ctx context.Context, cfg config.Config,
) (setupdb.Database, error) {
	log.Infof("connection string: %s", cfg.Database.ConnectionString)

	migrations := setupdb.NewGooseMigrations(filepath.Join(
		cfg.Database.MigrationsDir, "migrations"),
		log.Logger())

	DBBuilder := setupdb.NewPostgresDB(migrations,
		cfg.Database.ConnectionString,
		cfg.Database.MaxIdleConns,
		cfg.Database.MaxOpenConns)

	db, err := DBBuilder.Setup(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "setting up database")
	}

	return db, nil
}
