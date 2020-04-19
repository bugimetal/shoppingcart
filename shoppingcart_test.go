package shoppingcart

import (
	"testing"
	"time"
)

func TestShoppingCartItem_Validate(t *testing.T) {
	type fields struct {
		ID             int64
		ShoppingCartID int64
		ProductID      int64
		Quantity       uint64
		Price          uint64
		CreatedAt      time.Time
		UpdatedAt      time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "no product",
			fields: fields{
				ShoppingCartID: 1,
				Quantity:       1,
				Price:          1000000,
			},
			wantErr: true,
		},
		{
			name: "no quantity",
			fields: fields{
				ShoppingCartID: 1,
				ProductID:      1,
				Price:          1000000,
			},
			wantErr: true,
		},
		{
			name: "no price",
			fields: fields{
				ShoppingCartID: 1,
				ProductID:      1,
				Quantity:       1,
			},
			wantErr: true,
		},
		{
			name: "everything set",
			fields: fields{
				ShoppingCartID: 1,
				ProductID:      1,
				Quantity:       1,
				Price:          1000000,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cartItem := &ShoppingCartItem{
				ShoppingCartID: tt.fields.ShoppingCartID,
				ProductID:      tt.fields.ProductID,
				Quantity:       tt.fields.Quantity,
				Price:          tt.fields.Price,
			}
			if err := cartItem.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
