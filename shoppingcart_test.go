package shoppingcart

import (
	"reflect"
	"testing"
	"time"
)

func TestShoppingCartItem_Validate(t *testing.T) {
	type fields struct {
		ID             int64
		ShoppingCartID int64
		ProductID      int64
		Quantity       uint64
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
			},
			wantErr: true,
		},
		{
			name: "no quantity",
			fields: fields{
				ShoppingCartID: 1,
				ProductID:      1,
			},
			wantErr: true,
		},
		{
			name: "everything set",
			fields: fields{
				ShoppingCartID: 1,
				ProductID:      1,
				Quantity:       1,
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
			}
			if err := cartItem.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestShoppingCart_HasProduct(t *testing.T) {
	type fields struct {
		Items []ShoppingCartItem
	}
	type args struct {
		productID int64
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "product found",
			fields: fields{[]ShoppingCartItem{
				{ProductID: 1},
				{ProductID: 2},
				{ProductID: 3},
			}},
			args: args{productID: 3},
			want: true,
		},
		{
			name: "product not found",
			fields: fields{[]ShoppingCartItem{
				{ProductID: 1},
				{ProductID: 2},
				{ProductID: 3},
			}},
			args: args{productID: 5},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cart := &ShoppingCart{

				Items: tt.fields.Items,
			}
			if got := cart.HasProduct(tt.args.productID); got != tt.want {
				t.Errorf("HasProduct() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShoppingCart_GetProduct(t *testing.T) {
	type fields struct {
		Items []ShoppingCartItem
	}
	type args struct {
		productID int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    ShoppingCartItem
		wantErr bool
	}{
		{
			name: "product found",
			fields: fields{[]ShoppingCartItem{
				{ProductID: 1},
				{ProductID: 2},
				{ProductID: 3},
			}},
			args:    args{productID: 3},
			want:    ShoppingCartItem{ProductID: 3},
			wantErr: false,
		},
		{
			name: "product not found",
			fields: fields{[]ShoppingCartItem{
				{ProductID: 1},
				{ProductID: 2},
				{ProductID: 3},
			}},
			args:    args{productID: 5},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cart := &ShoppingCart{
				Items: tt.fields.Items,
			}
			got, err := cart.GetProduct(tt.args.productID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetProduct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetProduct() got = %v, want %v", got, tt.want)
			}
		})
	}
}
