//+build !test

package main

import (
	"flag"
	"net/http"
	"os"
)

// The functions, structs down below are the core methods,
// you shouldn't edit them until you know what you're doing,
// or you understand how KitSvc works.
//
// Or if you are brave enough ;)

func main() {

	// Command line flags.
	var (
		listenPort = flag.String("listen", ":"+os.Getenv("KITSVC_PORT"), "HTTP listen address")
		resetDB    = flag.Bool("reinitialize-database", false, "Set true to reinitialize the database, it's useful with the unit testing. The database will backed up before the database was reinitialized.")
	)
	flag.Parse()

	// Create the logger with the specified listen port.
	logger := createLogger(listenPort)
	// Create the database connection.
	db := createDatabase(resetDB)
	// Create the model with the database connection.
	model := createModel(db)
	// Create the messaging service with the logger.
	msg := createMessage()

	// Create the main service with what it needs.
	createService(logger, msg, model)
	// Register the service to the service registry.
	registerService(logger)

	// Log the ports.
	logger.Log("msg", "HTTP", "addr", *listenPort)
	// Start the service and listening to the requests.
	logger.Log("err", http.ListenAndServe(*listenPort, nil))
}
