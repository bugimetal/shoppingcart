package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bugimetal/shoppingcart"
	"github.com/bugimetal/shoppingcart/internal/mock/auth"
	shoppingcart_mock "github.com/bugimetal/shoppingcart/internal/mock/shoppingcart"
	"github.com/bugimetal/shoppingcart/service"

	"github.com/julienschmidt/httprouter"
)

func TestHandler_createShoppingCart(t *testing.T) {
	shoppingCartService := &shoppingcart_mock.MockShoppingCartService{}

	handler := &Handler{
		shoppingCartService: shoppingCartService,
		authService:         auth.New(),
	}

	t.Run("create shopping cart without auth", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := newRequest(http.MethodPost, "/shoppingcart", nil)

		handler.authMiddleware(handler.createShoppingCart)(w, r, nil)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("Expected HTTP status code %d, but got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("create shopping cart successful", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := newRequest(http.MethodPost, "/shoppingcart", nil)

		r.Header.Set("Authorization", "Basic dGVzdDp0ZXN0Cg==")

		handler.authMiddleware(handler.createShoppingCart)(w, r, nil)

		if w.Code != http.StatusCreated {
			t.Fatalf("Expected HTTP status code %d, but got %d", http.StatusCreated, w.Code)
		}
	})
}

func TestHandler_getShoppingCart(t *testing.T) {
	storageMock := &shoppingcart_mock.MockStorage{}

	services := service.New(service.Dependencies{
		ShoppingCartStorage: storageMock,
	})

	handler := &Handler{
		shoppingCartService: services.ShoppingCart,
		authService:         auth.New(),
	}

	t.Run("get shopping cart successful", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := newRequest(http.MethodGet, "/shoppingcart/1", nil)

		creds := base64.StdEncoding.EncodeToString([]byte(`test:test`))
		r.Header.Set("Authorization", fmt.Sprintf("Basic %s", creds))

		handler.authMiddleware(handler.getShoppingCart)(w, r, []httprouter.Param{{Key: "id", Value: "1"}})

		if w.Code != http.StatusOK {
			t.Fatalf("Expected HTTP status code %d, but got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("get shopping cart which doesn't belong to this user", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := newRequest(http.MethodGet, "/shoppingcart/1", nil)

		creds := base64.StdEncoding.EncodeToString([]byte(`hacker:password`))
		r.Header.Set("Authorization", fmt.Sprintf("Basic %s", creds))

		handler.authMiddleware(handler.getShoppingCart)(w, r, []httprouter.Param{{Key: "id", Value: "1"}})

		if w.Code != http.StatusNotFound {
			t.Fatalf("Expected HTTP status code %d, but got %d", http.StatusNotFound, w.Code)
		}
	})
}

func TestHandler_addProduct(t *testing.T) {
	storageMock := &shoppingcart_mock.MockStorage{}

	services := service.New(service.Dependencies{
		ShoppingCartStorage: storageMock,
	})

	handler := &Handler{
		shoppingCartService: services.ShoppingCart,
		authService:         auth.New(),
	}

	newProductCreate := shoppingcart.ShoppingCartItem{ProductID: 5, Quantity: 1}
	existingProductUpdate := shoppingcart.ShoppingCartItem{ProductID: 2, Quantity: 1}

	creds := base64.StdEncoding.EncodeToString([]byte(`test:test`))

	t.Run("add new product", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := newRequest(http.MethodPost, "/shoppingcart/1/item", newProductCreate)

		r.Header.Set("Authorization", fmt.Sprintf("Basic %s", creds))

		handler.authMiddleware(handler.addProduct)(w, r, []httprouter.Param{{Key: "id", Value: "1"}})

		if w.Code != http.StatusCreated {
			t.Fatalf("Expected HTTP status code %d, but got %d", http.StatusCreated, w.Code)
		}

		var createdProduct shoppingcart.ShoppingCartItem
		if err := json.NewDecoder(w.Body).Decode(&createdProduct); err != nil {
			t.Fatalf("Can't decode response: %v", err)
		}

		if createdProduct.ProductID != newProductCreate.ProductID {
			t.Fatalf("Expected %d product id, got %d", newProductCreate.ProductID, createdProduct.ProductID)
		}

		if createdProduct.Quantity != newProductCreate.Quantity {
			t.Fatalf("Expected %d product quantity, got %d", newProductCreate.Quantity, createdProduct.Quantity)
		}
	})

	t.Run("add existing product", func(t *testing.T) {
		// First we need to pull existing product to compare it later
		w := httptest.NewRecorder()
		r := newRequest(http.MethodGet, "/shoppingcart/1", nil)

		r.Header.Set("Authorization", fmt.Sprintf("Basic %s", creds))

		handler.authMiddleware(handler.getShoppingCart)(w, r, []httprouter.Param{{Key: "id", Value: "1"}})

		if w.Code != http.StatusOK {
			t.Fatalf("Expected HTTP status code %d, but got %d", http.StatusOK, w.Code)
		}

		var existingCart shoppingcart.ShoppingCart
		if err := json.NewDecoder(w.Body).Decode(&existingCart); err != nil {
			t.Fatalf("Can't decode response: %v", err)
		}

		existingProduct, err := existingCart.GetProduct(existingProductUpdate.ProductID)
		if err != nil {
			t.Fatalf("Product with ID %d not found in %v", existingProductUpdate.ProductID, existingCart.Items)
		}

		// Now as existing product found, we can perform update
		w = httptest.NewRecorder()
		r = newRequest(http.MethodPost, "/shoppingcart/1/item", existingProductUpdate)
		r.Header.Set("Authorization", fmt.Sprintf("Basic %s", creds))

		handler.authMiddleware(handler.addProduct)(w, r, []httprouter.Param{{Key: "id", Value: "1"}})

		if w.Code != http.StatusCreated {
			t.Fatalf("Expected HTTP status code %d, but got %d", http.StatusCreated, w.Code)
		}

		var createdProduct shoppingcart.ShoppingCartItem
		if err := json.NewDecoder(w.Body).Decode(&createdProduct); err != nil {
			t.Fatalf("Can't decode response: %v", err)
		}

		if createdProduct.ProductID != existingProductUpdate.ProductID {
			t.Fatalf("Expected %d product id, got %d", existingProductUpdate.ProductID, createdProduct.ProductID)
		}

		expectedQuantity := existingProduct.Quantity + existingProductUpdate.Quantity
		if createdProduct.Quantity != expectedQuantity {
			t.Fatalf("Expected %d product quantity, got %d", expectedQuantity, createdProduct.Quantity)
		}
	})
}
