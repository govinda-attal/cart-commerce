package httputil

import (
	"encoding/json"
	"net/http"

	"github.com/govinda-attal/cart-commerce/pkg/core/status"
)

// WrapperHandler is wrapper function to wrap API handlers and retuns as http.HandlerFunc.
// API Handlers may return error, and this wrapper simplifies error handling for API Handlers.
func WrapperHandler(f func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			w.Header().Add("Content-Type", "application/json")
			if errSvc, ok := err.(status.ErrServiceStatus); ok {
				w.WriteHeader(errSvc.Code)
				errSvc.Code = 0 // Zero value to avoid value within JSON encoded message.
				json.NewEncoder(w).Encode(&errSvc)
				return
			}
			errSvc := status.ErrInternal.WithMessage(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(&errSvc)
		}
	}
}
