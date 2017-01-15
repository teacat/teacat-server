package main

import (
	"net/http"
	"os"

	"github.com/TeaMeow/KitSvc/router"
	"github.com/TeaMeow/KitSvc/router/middleware"
)

func main() {
	handler := router.Load(middleware.Store())

	http.ListenAndServe(os.Getenv("KITSVC_ADDR"), handler)
}
