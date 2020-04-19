package shoppingcart

import (
	"context"

	"github.com/bugimetal/shoppingcart"
)

type MockStorage struct{}

func (db *MockStorage) Create(ctx context.Context, cart *shoppingcart.ShoppingCart) error {
	return nil
}

func (db *MockStorage) Get(ctx context.Context, ID, userID int64) (shoppingcart.ShoppingCart, error) {
	if ID == 1 && userID == 1 {
		return shoppingcart.ShoppingCart{
			ID:     1,
			UserID: 1,
			Items: []shoppingcart.ShoppingCartItem{
				{ProductID: 1, Quantity: 1},
				{ProductID: 2, Quantity: 10},
			},
		}, nil
	}

	return shoppingcart.ShoppingCart{}, shoppingcart.ErrCartNotFound
}

func (db *MockStorage) Empty(ctx context.Context, shoppingCartID int64) error {
	return nil
}

func (db *MockStorage) AddProduct(ctx context.Context, cartItem *shoppingcart.ShoppingCartItem) error {
	return nil
}

func (db *MockStorage) UpdateProduct(ctx context.Context, cartItem *shoppingcart.ShoppingCartItem) error {
	return nil
}

func (db *MockStorage) RemoveProduct(ctx context.Context, shoppingCartID, productID int64) error {
	return nil
}
