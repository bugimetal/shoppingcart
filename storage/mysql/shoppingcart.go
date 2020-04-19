package mysql

import (
	"context"
	"time"

	"github.com/bugimetal/shoppingcart"
)

// Create creates shopping cart in the storate
func (db *DB) Create(ctx context.Context, cart *shoppingcart.ShoppingCart) error {
	cart.CreatedAt = time.Now()
	cart.UpdatedAt = cart.CreatedAt

	db.client.Create(cart)

	return nil
}

// Get retrieves shopping cart from the storage along with items
func (db *DB) Get(ctx context.Context, ID, userID int64) (shoppingcart.ShoppingCart, error) {
	var cart shoppingcart.ShoppingCart
	db.client.
		Preload("Items").
		Where("shoppingcart.id = ? AND user_id = ?", ID, userID).
		First(&cart)

	if cart.ID == 0 {
		return cart, shoppingcart.ErrCartNotFound
	}

	return cart, nil
}

// Empty removes items associated with shopping cart
func (db *DB) Empty(ctx context.Context, shoppingCartID int64) error {
	db.client.Where("shoppingcart_id = ?", shoppingCartID).Delete(shoppingcart.ShoppingCartItem{})
	return nil
}

// AddProduct adds product to the shopping cart
func (db *DB) AddProduct(ctx context.Context, cartItem *shoppingcart.ShoppingCartItem) error {
	cartItem.CreatedAt = time.Now()
	cartItem.UpdatedAt = cartItem.CreatedAt

	db.client.Create(cartItem)

	return nil
}

// UpdateProduct updates product in the shopping cart
func (db *DB) UpdateProduct(ctx context.Context, cartItem *shoppingcart.ShoppingCartItem) error {
	cartItem.UpdatedAt = time.Now()

	db.client.Save(cartItem)

	return nil
}

// RemoveProduct removes product from the shopping cart
func (db *DB) RemoveProduct(ctx context.Context, shoppingCartID, productID int64) error {
	db.client.
		Where("shoppingcart_id = ? AND product_id = ?", shoppingCartID, productID).
		Delete(shoppingcart.ShoppingCartItem{})

	return nil
}
