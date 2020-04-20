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
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	httpServer := &http.Server{
		Addr:    *bind,
		Handler: h,
	}

	// Start the HTTP server.
	httpServerErrorChan := make(chan error)
	go func() {
		fmt.Printf("HTTP server listening on %s\n", *bind)
		httpServerErrorChan <- httpServer.ListenAndServe()
	}()

	// Set up the signal channel.
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	// If the HTTP server returned an error, exit here.
	case err := <-httpServerErrorChan:
		log.Printf("HTTP server error: %s", err)
	// If a termination signal was received, shutdown the server.
	case sig := <-signalChan:
		log.Printf("Signal received: %s", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP Server graceful shutdown failed with an error: %s\n", err)
	}
}
