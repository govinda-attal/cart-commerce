package eshop

import "context"

type CartCommerce interface {
	UpdateCartItem(ctx context.Context, cartID string, itemName string, quantity int) (string, error)
	FetchCartItems(ctx context.Context, cartID string) (*Cart, error)
	UpdateCartState(ctx context.Context, cartID string, state CartState) error
}

type Promotions interface {
	Fetch(ctx context.Context) (*PromotionRules, error)
}
