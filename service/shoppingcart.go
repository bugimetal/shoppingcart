package service

import (
	"context"

	"github.com/bugimetal/shoppingcart"
)

type ShoppingCart struct {
	storage ShoppingCartStorage
}

func NewShoppingCart(deps Dependencies) *ShoppingCart {
	return &ShoppingCart{storage: deps.ShoppingCartStorage}
}

// Create creates a new shopping cart in storage
func (service *ShoppingCart) Create(ctx context.Context, cart *shoppingcart.ShoppingCart) error {
	if err := cart.Validate(); err != nil {
		return err
	}

	return service.storage.Create(ctx, cart)
}

// Get retrieves a shopping cart from the storage
func (service *ShoppingCart) Get(ctx context.Context, ID, userID int64) (shoppingcart.ShoppingCart, error) {
	return service.storage.Get(ctx, ID, userID)
}

// Empty removes items associated with shopping cart
func (service *ShoppingCart) Empty(ctx context.Context, shoppingCartID int64) error {
	return service.storage.Empty(ctx, shoppingCartID)
}

func (service *ShoppingCart) AddProduct(ctx context.Context, cartItem *shoppingcart.ShoppingCartItem) error {
	if err := cartItem.Validate(); err != nil {
		return err
	}

	return service.storage.AddProduct(ctx, cartItem)
}

func (service *ShoppingCart) RemoveProduct(ctx context.Context, shoppingCartID, productID int64) error {
	return service.storage.RemoveProduct(ctx, shoppingCartID, productID)
}
