package storage

import (
	"context"

	"github.com/bugimetal/shoppingcart"
)

// ShoppingCart describes an interface to store shopping carts and manipulate with products inside the cart
type ShoppingCart interface {
	Create(context.Context, *shoppingcart.ShoppingCart) error
	Get(ctx context.Context, shoppingCartID int64, userID int64) (shoppingcart.ShoppingCart, error)
	Empty(ctx context.Context, shoppingCartID int64) error

	AddProduct(context.Context, *shoppingcart.ShoppingCartItem) error
	UpdateProduct(context.Context, *shoppingcart.ShoppingCartItem) error
	RemoveProduct(ctx context.Context, shoppingCartID, productID int64) error
}
