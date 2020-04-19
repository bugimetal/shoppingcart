package shoppingcart

import (
	"context"

	"github.com/bugimetal/shoppingcart"
)

//
const hackerUserID = 2

type MockShoppingCartService struct{}

func (service *MockShoppingCartService) Create(ctx context.Context, cart *shoppingcart.ShoppingCart) error {
	if err := cart.Validate(); err != nil {
		return err
	}

	cart.ID = 1
	return nil
}

func (service *MockShoppingCartService) Get(ctx context.Context, ID, userID int64) (shoppingcart.ShoppingCart, error) {
	if userID == hackerUserID {
		return shoppingcart.ShoppingCart{}, shoppingcart.ErrCartNotFound
	}
	return shoppingcart.ShoppingCart{}, nil
}

func (service *MockShoppingCartService) Empty(ctx context.Context, shoppingCartID, userID int64) error {
	return nil
}

func (service *MockShoppingCartService) AddProduct(ctx context.Context, cartItem *shoppingcart.ShoppingCartItem) error {
	if err := cartItem.Validate(); err != nil {
		return err
	}

	return nil
}

func (service *MockShoppingCartService) RemoveProduct(ctx context.Context, shoppingCartID, productID, userID int64) error {
	return nil
}
