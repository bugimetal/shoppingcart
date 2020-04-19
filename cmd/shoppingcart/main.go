// Shopping cart service
//
// Documentation of API
//
//     Schemes: http
//     BasePath: /
//     Version: 1
//     Host: localhost:8080
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - basic
//
//    SecurityDefinitions:
//    basic:
//      type: basic
//
// swagger:meta
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/bugimetal/shoppingcart/handler"
	"github.com/bugimetal/shoppingcart/internal/mock/auth"
	"github.com/bugimetal/shoppingcart/service"
	"github.com/bugimetal/shoppingcart/storage/mysql"
)

var (
	bind = flag.String("bind", ":8080", "The socket to bind the HTTP server")
)

func main() {
	flag.Parse()

	config, err := NewConfig()
	if err != nil {
		log.Fatalf("Can't read the config: %v", err)
	}

	storage, err := mysql.New(
		config.Database.User,
		config.Database.Password,
		config.Database.Host,
		config.Database.Database,
		config.Database.Port,
	)
	defer storage.Close()

	if err != nil {
		log.Fatalf("Can't connect to database: %v", err)
	}

	// Service covers the high-level business logic.
	services := service.New(service.Dependencies{
		ShoppingCartStorage: storage,
	})

	// Initializing external authorisation service
	authService := auth.New()

	h := handler.New(handler.Services{
		ShoppingCart: services.ShoppingCart,
		Auth:         authService,
	})

	fmt.Printf("Listening on %s\n", *bind)
	log.Fatal(http.ListenAndServe(*bind, h))
}
