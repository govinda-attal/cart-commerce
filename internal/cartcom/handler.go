package cartcom

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/govinda-attal/cart-commerce/pkg/core/status"
	"github.com/govinda-attal/cart-commerce/pkg/eshop"
)

func NewHandler(api eshop.CartCommerce) *handlerImpl {
	return &handlerImpl{api}
}

type handlerImpl struct {
	api eshop.CartCommerce
}

func (h *handlerImpl) UpdateCartItem(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	vars := mux.Vars(r)
	cartID := vars["cartId"]

	var item eshop.Item
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		return status.ErrBadRequest
	}

	cartID, err = h.api.UpdateCartItem(ctx, cartID, item.ItemName, item.Quantity)
	if err != nil {
		return err
	}

	w.Header().Set("App-CartID", cartID)
	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&status.Success)

	if err != nil {
		return status.ErrInternal.WithError(err)
	}

	return nil
}

func (h *handlerImpl) UpdateCartState(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	vars := mux.Vars(r)
	cartID, ok := vars["cartId"]
	if !ok {
		return status.ErrBadRequest
	}

	r.ParseForm()
	state := r.Form.Get("state")
	if len(state) == 0 {
		return status.ErrBadRequest
	}

	err := h.api.UpdateCartState(ctx, cartID, state)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (h *handlerImpl) FetchCartItems(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	vars := mux.Vars(r)
	cartID, ok := vars["cartId"]
	if !ok {
		return status.ErrBadRequest
	}

	cart, err := h.api.FetchCartItems(ctx, cartID)
	if err != nil {
		return err
	}
	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(cart)
	if err != nil {
		return status.ErrInternal.WithError(err)
	}
	return nil
}
