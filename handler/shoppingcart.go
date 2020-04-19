package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bugimetal/shoppingcart"
	"github.com/sirupsen/logrus"

	"github.com/julienschmidt/httprouter"
)

// swagger:operation POST /v1/shoppingcart ShoppingCart createShoppingCart
// ---
// summary: Creates a shopping cart item and persist it in the storage
// description:
// responses:
//   "201":
//     "$ref": "#/responses/ShoppingCart"
//   "401":
//     "$ref": "#/responses/errorResponse"
//   "500":
//     "$ref": "#/responses/errorResponse"
func (handler *Handler) createShoppingCart(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user, err := handler.authUser(r)
	if err != nil {
		handler.Error(w, r, err)
		return
	}

	cart := shoppingcart.ShoppingCart{
		UserID: user.ID,
	}

	if err := handler.shoppingCartService.Create(r.Context(), &cart); err != nil {
		handler.Error(w, r, err)
		logrus.Errorf("Unable to create shopping cart for user %d: %s", user.ID, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(cart)
}

// swagger:operation GET /v1/shoppingcart/{id} ShoppingCart getShoppingCart
// ---
// summary: Retrieves shopping cart from the storage along with shopping cart items
// description:
// parameters:
// - name: id
//   in: path
//   description: shopping cart id
//   required: true
//   type: integer
//   format: int64
// responses:
//   "200":
//     "$ref": "#/responses/ShoppingCart"
//   "401":
//     "$ref": "#/responses/errorResponse"
//   "404":
//     "$ref": "#/responses/errorResponse"
//   "500":
//     "$ref": "#/responses/errorResponse"
func (handler *Handler) getShoppingCart(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ID, err := strconv.ParseInt(ps.ByName("id"), 10, 64)
	if err != nil {
		handler.Error(w, r, err)
		return
	}

	user, err := handler.authUser(r)
	if err != nil {
		handler.Error(w, r, err)
		return
	}

	cart, err := handler.shoppingCartService.Get(r.Context(), ID, user.ID)
	if err != nil {
		handler.Error(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(cart)
}

// swagger:operation DELETE /v1/shoppingcart/{id}/item ShoppingCart emptyCart
// ---
// summary: Removes shopping cart items from storage
// description: If shopping cart has no items to delete, error will be returned
// parameters:
// - name: id
//   in: path
//   description: shopping cart id
//   required: true
//   type: integer
//   format: int64
// responses:
//   "204":
//   "401":
//     "$ref": "#/responses/errorResponse"
//   "404":
//     "$ref": "#/responses/errorResponse"
//   "500":
//     "$ref": "#/responses/errorResponse"
func (handler *Handler) emptyCart(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ID, err := strconv.ParseInt(ps.ByName("id"), 10, 64)
	if err != nil {
		handler.Error(w, r, err)
		return
	}

	user, err := handler.authUser(r)
	if err != nil {
		handler.Error(w, r, err)
		return
	}

	// Checking if shopping cart belong to this user
	cart, err := handler.shoppingCartService.Get(r.Context(), ID, user.ID)
	if err != nil {
		handler.Error(w, r, err)
		return
	}

	if len(cart.Items) == 0 {
		handler.Error(w, r, shoppingcart.ErrCartHasNoItems)
		return
	}

	if err := handler.shoppingCartService.Empty(r.Context(), ID); err != nil {
		handler.Error(w, r, err)
		logrus.Errorf("Unable to add empty shopping cart %d: %s", ID, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// swagger:operation POST /v1/shoppingcart/{id}/item ShoppingCartItem addProduct
// ---
// summary: add product to existing shopping cart
// description: If product already added to specified shopping cart, error will be returned
// parameters:
// - name: id
//   in: path
//   description: shopping cart id
//   required: true
//   type: integer
//   format: int64
// - name: item
//   in: body
//   description: shopping cart item
//   required: true
//   schema:
//     "$ref": "#/definitions/ShoppingCartItem"
// responses:
//   "201":
//     "$ref": "#/responses/ShoppingCartItem"
//   "401":
//     "$ref": "#/responses/errorResponse"
//   "404":
//     "$ref": "#/responses/errorResponse"
//   "500":
//     "$ref": "#/responses/errorResponse"
func (handler *Handler) addProduct(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var cartItem shoppingcart.ShoppingCartItem
	if err := json.NewDecoder(r.Body).Decode(&cartItem); err != nil {
		handler.Error(w, r, err)
		return
	}

	cartID, err := strconv.ParseInt(ps.ByName("id"), 10, 64)
	if err != nil {
		handler.Error(w, r, err)
		return
	}

	user, err := handler.authUser(r)
	if err != nil {
		handler.Error(w, r, err)
		return
	}

	// Checking if shopping cart belong to this user
	cart, err := handler.shoppingCartService.Get(r.Context(), cartID, user.ID)
	if err != nil {
		handler.Error(w, r, err)
		return
	}

	if cart.HasProduct(cartItem.ProductID) {
		// Product already added to the cart, returning error.
		// Product can be also updated in this case, it depends on requirements
		handler.Error(w, r, shoppingcart.ErrCartItemAlreadyExists)
		return
	}

	cartItem.ShoppingCartID = cart.ID

	if err := handler.shoppingCartService.AddProduct(r.Context(), &cartItem); err != nil {
		handler.Error(w, r, err)
		logrus.Errorf("Unable to add shopping cart item %v: %s", cartItem, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(cartItem)

}

// swagger:operation DELETE /v1/shoppingcart/{id}/item/{product_id} ShoppingCartItem removeProduct
// ---
// summary: removes product from existing shopping cart
// description: If product doesn't exists in specified shopping cart, error will be returned
// parameters:
// - name: id
//   in: path
//   description: shopping cart id
//   required: true
//   type: integer
//   format: int64
// - name: product_id
//   in: path
//   description: product id to delete
//   required: true
//   type: integer
//   format: int64
// responses:
//   "204":
//   "401":
//     "$ref": "#/responses/errorResponse"
//   "404":
//     "$ref": "#/responses/errorResponse"
//   "500":
//     "$ref": "#/responses/errorResponse"
func (handler *Handler) removeProduct(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cartID, err := strconv.ParseInt(ps.ByName("id"), 10, 64)
	if err != nil {
		handler.Error(w, r, err)
		return
	}

	productID, err := strconv.ParseInt(ps.ByName("product_id"), 10, 64)
	if err != nil {
		handler.Error(w, r, err)
		return
	}

	user, err := handler.authUser(r)
	if err != nil {
		handler.Error(w, r, err)
		return
	}

	// Checking if shopping cart belong to this user
	cart, err := handler.shoppingCartService.Get(r.Context(), cartID, user.ID)
	if err != nil {
		handler.Error(w, r, err)
		return
	}

	if !cart.HasProduct(productID) {
		handler.Error(w, r, shoppingcart.ErrCartItemNotFound)
		return
	}

	if err := handler.shoppingCartService.RemoveProduct(r.Context(), cartID, productID); err != nil {
		handler.Error(w, r, err)
		logrus.Errorf("Unable to remove shopping cart item (%d:%d): %s", cartID, productID, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
