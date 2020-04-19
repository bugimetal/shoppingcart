package shoppingcart

import (
	"errors"
	"time"
)

// These errors can be returned by the service or storage packages when performing
// operating on ShoppingCart
var (
	ErrNoPermission = errors.New("user does not have the permissions")
	ErrUserNotSet   = errors.New("no user set")

	ErrCartNotFound   = errors.New("shopping cart not found")
	ErrCartHasNoItems = errors.New("shopping cart has no items")

	ErrCartItemNotFound      = errors.New("shopping cart item not found")
	ErrCartItemAlreadyExists = errors.New("this product already added to shopping cart")
	ErrCartItemNoProductSet  = errors.New("product is not specified")
	ErrCartItemNoQuantitySet = errors.New("quantity is not specified")
)

// ShoppingCart describes shopping cart
// swagger:response ShoppingCart
type ShoppingCart struct {
	ID        int64              `json:"id" gorm:"primary_key"`
	UserID    int64              `json:"user_id" gorm:"primary_key"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	Items     []ShoppingCartItem `json:"items,omitempty" gorm:"foreignkey:ShoppingCartID;association_foreignkey:ID"`
}

func (ShoppingCart) TableName() string {
	return "shoppingcart"
}

// Validate this shopping cart
func (cart *ShoppingCart) Validate() error {
	if cart.UserID == 0 {
		return ErrUserNotSet
	}
	return nil
}

func (cart *ShoppingCart) HasProduct(productID int64) bool {
	for _, item := range cart.Items {
		if item.ProductID == productID {
			return true
		}
	}
	return false
}

func (cart *ShoppingCart) GetProduct(productID int64) (ShoppingCartItem, error) {
	for _, item := range cart.Items {
		if item.ProductID == productID {
			return item, nil
		}
	}
	return ShoppingCartItem{}, ErrCartItemNotFound
}

// ShoppingCartItem represents shopping cart entity
// swagger:response ShoppingCartItem
type ShoppingCartItem struct {
	ID             int64     `json:"id"`
	ShoppingCartID int64     `json:"shoppingcart_id" gorm:"column:shoppingcart_id"`
	ProductID      int64     `json:"product_id"`
	Quantity       uint64    `json:"quantity"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (ShoppingCartItem) TableName() string {
	return "shoppingcart_item"
}

func (cartItem *ShoppingCartItem) Validate() error {
	switch {
	case cartItem.ProductID == 0:
		return ErrCartItemNoProductSet
	case cartItem.Quantity == 0:
		return ErrCartItemNoQuantitySet
	}

	return nil
}
