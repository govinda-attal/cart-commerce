package promo

import (
	"encoding/json"
	"net/http"

	"github.com/govinda-attal/cart-commerce/pkg/core/status"
	"github.com/govinda-attal/cart-commerce/pkg/eshop"
)

func NewHandler(api eshop.Promotions) *handlerImpl {
	return &handlerImpl{api}
}

type handlerImpl struct {
	api eshop.Promotions
}

func (h *handlerImpl) Fetch(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	p, err := h.api.Fetch(ctx)
	if err != nil {
		return err
	}
	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(p)
	if err != nil {
		return status.ErrInternal.WithError(err)
	}
	return nil
}
