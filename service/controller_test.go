package service_test

import (
	"os"
	"testing"

	"github.com/TeaMeow/KitSvc/config"
	"github.com/TeaMeow/KitSvc/database"
	"github.com/TeaMeow/KitSvc/instrumenting"
	"github.com/TeaMeow/KitSvc/logging"
	"github.com/TeaMeow/KitSvc/messaging"
	"github.com/TeaMeow/KitSvc/service"
	"github.com/go-kit/kit/log"
)

var (
	svc service.Service
)

func init() {
	// Load the configurations.
	conf := config.Load(".././")

	// The logger.
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.NewContext(logger).With("listen", ":8080").With("caller", log.DefaultCaller)

	// Connect to the database.
	db := database.Connect(conf.Database)

	// Create the model with the database connection.
	model := service.Model{DB: db}

	// Register to the message system.
	msg := messaging.Create(conf, logger)

	// Service and the middlewares.
	svc = service.Concrete{Message: msg.Producer, Model: model}
	svc = logging.CreateMiddleware(logger)(svc)
	svc = instrumenting.CreateMiddleware(conf)(svc)
}

func TestUppercase(t *testing.T) {
	svc.Uppercase("test")
	svc.Uppercase("")
}

func TestCount(t *testing.T) {
	svc.Count("test")
}
