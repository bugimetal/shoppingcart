package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/bugimetal/shoppingcart"

	"github.com/sirupsen/logrus"
)

// ErrorStatusCodes maps commonly returned errors to HTTP status codes
var ErrorStatusCodes = map[error]int{
	// Shopping cart
	shoppingcart.ErrCartNotFound:   http.StatusNotFound,
	shoppingcart.ErrCartHasNoItems: http.StatusBadRequest,
	shoppingcart.ErrNoPermission:   http.StatusUnauthorized,

	// Shopping cart item
	shoppingcart.ErrCartItemNoProductSet:  http.StatusBadRequest,
	shoppingcart.ErrCartItemNoQuantitySet: http.StatusBadRequest,
	shoppingcart.ErrCartItemNotFound:      http.StatusNotFound,
	shoppingcart.ErrCartItemAlreadyExists: http.StatusBadRequest,
}

// errorResponse represents error response structure
// swagger:response errorResponse
type errorResponse struct {
	Resource *errorResource `json:"error"`
}

type errorResource struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func newErrorResponse(err error) *errorResponse {
	return &errorResponse{Resource: &errorResource{Code: statusCode(err), Message: err.Error()}}
}

func (er *errorResource) Error() string {
	s := strings.SplitAfter(er.Message, ": ")
	msg := s[len(s)-1]

	return msg
}

// statusCode returns the HTTP status code that is appropriate for the specified error.
func statusCode(err error) int {
	if statusCode, ok := ErrorStatusCodes[err]; ok {
		return statusCode
	}

	return http.StatusInternalServerError
}

func (handler *Handler) Error(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(statusCode(err))

	if err := json.NewEncoder(w).Encode(newErrorResponse(err)); err != nil {
		logrus.Errorf("unable to decode struct to json: %s", err)
	}
}
