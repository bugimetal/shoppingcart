package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

// newRequest returns a new http request. If a body has been specified, it will
// be converted to a JSON data structure.
func newRequest(method, path string, body interface{}) *http.Request {
	var r *http.Request
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			panic(err)
		}

		r = httptest.NewRequest(method, path, bytes.NewReader(data))
	} else {
		r = httptest.NewRequest(method, path, nil)
	}

	return r
}
