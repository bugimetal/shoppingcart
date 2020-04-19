package service

import "github.com/bugimetal/shoppingcart/storage"

// ShoppingCartStorage describes the interface to store and retrieve shopping carts.
type ShoppingCartStorage interface {
	storage.ShoppingCart
}

// Dependencies list the interfaces that individual services rely on.
type Dependencies struct {
	ShoppingCartStorage
}

// Services contains all the services that this package has to offer.
type Services struct {
	*ShoppingCart
}

// New returns Services.
func New(deps Dependencies) *Services {
	return &Services{
		ShoppingCart: NewShoppingCart(deps),
	}
}
