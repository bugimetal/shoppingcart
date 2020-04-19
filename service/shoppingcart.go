package service

import (
	"context"

	"github.com/bugimetal/shoppingcart"
)

// ShoppingCart service responsible for shopping cart operations
type ShoppingCart struct {
	storage ShoppingCartStorage
}

// NewShoppingCart returns a new Shopping cart service
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
func (service *ShoppingCart) Empty(ctx context.Context, shoppingCartID, userID int64) error {
	// Checking if shopping cart belong to this user
	cart, err := service.Get(ctx, shoppingCartID, userID)
	if err != nil {
		return err
	}

	if len(cart.Items) == 0 {
		return nil
	}

	return service.storage.Empty(ctx, shoppingCartID)
}

// AddProduct adds new product to existing shopping cart
// If product exists, quantity will be updated
func (service *ShoppingCart) AddProduct(ctx context.Context, cartItem *shoppingcart.ShoppingCartItem, userID int64) error {
	if err := cartItem.Validate(); err != nil {
		return err
	}

	cart, err := service.Get(ctx, cartItem.ShoppingCartID, userID)
	if err != nil {
		return err
	}

	if cart.HasProduct(cartItem.ProductID) {
		existingItem, err := cart.GetProduct(cartItem.ProductID)
		if err != nil {
			return err
		}
		existingItem.Quantity += cartItem.Quantity
		*cartItem = existingItem
		return service.storage.UpdateProduct(ctx, cartItem)
	}

	return service.storage.AddProduct(ctx, cartItem)
}

// RemoveProduct removes a product from existing shopping cart
func (service *ShoppingCart) RemoveProduct(ctx context.Context, shoppingCartID, productID, userID int64) error {
	// Checking if shopping cart belong to this user
	cart, err := service.Get(ctx, shoppingCartID, userID)
	if err != nil {
		return err
	}

	if !cart.HasProduct(productID) {
		return nil
	}

	return service.storage.RemoveProduct(ctx, shoppingCartID, productID)
}
