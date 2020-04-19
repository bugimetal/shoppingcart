package handler

import (
	"context"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/bugimetal/shoppingcart/internal/mock/auth"

	"github.com/bugimetal/shoppingcart"
	"github.com/julienschmidt/httprouter"
)

type key int

var userKey key

// ShoppingCartService provides an interface to the service that deals with operations
// on shopping cart and cart items.
type ShoppingCartService interface {
	Create(context.Context, *shoppingcart.ShoppingCart) error
	Get(ctx context.Context, shoppingCartID int64, userID int64) (shoppingcart.ShoppingCart, error)
	Empty(ctx context.Context, shoppingCartID int64) error

	AddProduct(context.Context, *shoppingcart.ShoppingCartItem) error
	RemoveProduct(ctx context.Context, shoppingCartID, productID int64) error
}

// AuthService provides an interface to the service that deals with user authentication.
type AuthService interface {
	Authenticate(context.Context, auth.User) (auth.User, error)
}

// Services describe the external services that the Handler relies on.
type Services struct {
	ShoppingCart ShoppingCartService
	Auth         AuthService
}

// Handler provides an generic interface for handling HTTP requests.
type Handler struct {
	http                http.Handler
	shoppingCartService ShoppingCartService
	authService         AuthService
}

// New returns a new Handler.
func New(services Services) *Handler {
	handler := &Handler{
		shoppingCartService: services.ShoppingCart,
		authService:         services.Auth,
	}

	// Set up a custom HTTP router and install the routes on it.
	router := httprouter.New()

	router.POST("/v1/shoppingcart", handler.authMiddleware(handler.createShoppingCart))
	router.GET("/v1/shoppingcart/:id", handler.authMiddleware(handler.getShoppingCart))
	router.DELETE("/v1/shoppingcart/:id/item", handler.authMiddleware(handler.emptyCart))

	router.POST("/v1/shoppingcart/:id/item", handler.authMiddleware(handler.addProduct))
	router.DELETE("/v1/shoppingcart/:id/item/:product_id", handler.authMiddleware(handler.removeProduct))

	// Running swagger API documentation
	router.ServeFiles("/swagger/*filepath", http.Dir("./swagger/"))

	handler.http = router

	return handler
}

// ServeHTTP handles every incoming HTTP request and passes the request along
// to the configured HTTP router.
func (handler *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler.http.ServeHTTP(w, r)
}

// authUser returns the authenticated user.
func (handler *Handler) authUser(r *http.Request) (auth.User, error) {
	user, ok := r.Context().Value(userKey).(auth.User)
	if !ok {
		return user, shoppingcart.ErrNoPermission
	}

	return user, nil
}

// authMiddleware authenticates user using auth service
func (handler *Handler) authMiddleware(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Using Basic Auth just to save development time.
		if r.Header.Get("Authorization") == "" {
			handler.Error(w, r, shoppingcart.ErrNoPermission)
			return
		}

		authHeader := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

		if len(authHeader) != 2 || authHeader[0] != "Basic" {
			handler.Error(w, r, shoppingcart.ErrNoPermission)
			return
		}

		payload, _ := base64.StdEncoding.DecodeString(authHeader[1])
		pair := strings.SplitN(string(payload), ":", 2)

		user := auth.User{
			Name:     pair[0],
			Password: []byte(pair[1]),
		}

		user, err := handler.authService.Authenticate(r.Context(), user)
		if err != nil {
			handler.Error(w, r, shoppingcart.ErrNoPermission)
			return
		}

		ctx := context.WithValue(r.Context(), userKey, user)
		next(w, r.WithContext(ctx), ps)
	}
}
