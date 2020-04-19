package handler

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"

	"github.com/bugimetal/shoppingcart/internal/mock/auth"

	"github.com/bugimetal/shoppingcart/internal/mock/shoppingcart"
)

func TestHandler_createShoppingCart(t *testing.T) {
	shoppingCartService := &shoppingcart.MockShoppingCartService{}

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
	shoppingCartService := &shoppingcart.MockShoppingCartService{}

	handler := &Handler{
		shoppingCartService: shoppingCartService,
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

	t.Run("get shopping cart which don't belong to you", func(t *testing.T) {
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
