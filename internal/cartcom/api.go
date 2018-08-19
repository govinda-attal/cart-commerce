package cartcom

import (
	"context"
	"sort"

	"github.com/govinda-attal/cart-commerce/pkg/core/dbsqlx"
	"github.com/govinda-attal/cart-commerce/pkg/core/status"
	"github.com/govinda-attal/cart-commerce/pkg/eshop"
)

func NewApiImpl(db dbsqlx.DB, promos eshop.Promotions) *cartComImpl {
	return &cartComImpl{db, promos}
}

type cartComImpl struct {
	db     dbsqlx.DB
	promos eshop.Promotions
}

func (cc *cartComImpl) UpdateCartItem(ctx context.Context, cartID string, itemName string, quantity int) (string, error) {
	err := cc.db.QueryRow("SELECT updateCart($1,$2,$3)", cartID, itemName, quantity).Scan(&cartID)
	if err != nil {
		return cartID, status.ErrInternal.WithError(err)
	}
	return cartID, nil
}

func (cc *cartComImpl) UpdateCartState(ctx context.Context, cartID string, state eshop.CartState) error {
	var result bool
	err := cc.db.QueryRow("SELECT updateCartState($1,$2)", cartID, state).Scan(&result)
	if err != nil {
		return status.ErrInternal.WithError(err)
	}
	return nil
}

func (cc *cartComImpl) FetchCartItems(ctx context.Context, cartID string) (*eshop.Cart, error) {
	var cart eshop.Cart
	err := cc.db.QueryRow("SELECT C.State FROM Cart C WHERE C.ID = $1", cartID).Scan(&cart.State)

	if err != nil {
		return nil, status.ErrInternal.WithError(err)
	}

	if cart.State != eshop.CartStateInProgress {
		return nil, status.ErrNotFound
	}

	rows, err := cc.db.Query("SELECT CI.ItemName, CI.Quantity, I.Price AS UnitPrice FROM CartItem CI JOIN Inventory I ON CI.ItemName = I.ItemName WHERE CI.cartID = $1", cartID)
	if err != nil {
		return nil, status.ErrInternal.WithError(err)
	}

	cart.ID = cartID
	for rows.Next() {
		var item eshop.Item
		err := rows.Scan(&item.ItemName, &item.Quantity, &item.UnitPrice)
		if err != nil {
			return nil, status.ErrInternal.WithError(err)
		}
		price := eshop.Money(float64(*item.UnitPrice) * float64(item.Quantity))
		item.Price = &price
		cart.Items = append(cart.Items, &item)
	}

	sort.Sort(sort.Reverse(eshop.CartItemsByTotalPrice(cart.Items)))

	promoRules, err := cc.promos.Fetch(ctx)
	if err != nil {
		return nil, err
	}
	for _, item := range cart.Items {
		if itemRules := promoRules.FindApplicablePromos(item.ItemName, item.Quantity); len(itemRules) > 0 {
			for _, rule := range itemRules {
				rule.ApplyPromos(cart.Items)
			}
		}
	}
	return &cart, nil
}
