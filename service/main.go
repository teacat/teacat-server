// +build !test

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
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
		resetDB    = flag.Bool("reinitialize-database", false, "Set true to reinitialize the database, it's useful with the unit testing.")
	)
	flag.Parse()

	r, _ := strconv.ParseBool(fmt.Sprint(reflect.ValueOf(resetDB).Elem()))

	// Create the logger with the specified listen port.
	logger := createLogger(listenPort)
	// Create the database connection.
	db := createDatabase(r)
	// Create the model with the database connection.
	model := createModel(db)
	//
	es := createEventStore()

	// Create the main service with what it needs.
	createService(logger, model, es)
	// Register the service to the service registry.
	registerService(logger)

	// Log the ports.
	logger.Log("msg", "HTTP", "addr", *listenPort)
	// Start the service and listening to the requests.
	logger.Log("err", http.ListenAndServe(*listenPort, nil))
}
