package promo

import (
	"context"
	"io/ioutil"

	"github.com/govinda-attal/cart-commerce/pkg/core/status"
	"github.com/govinda-attal/cart-commerce/pkg/eshop"
	"gopkg.in/yaml.v2"
)

func NewApiImpl() *promoImpl {
	return &promoImpl{}
}

type promoImpl struct{}

func (pi *promoImpl) Fetch(ctx context.Context) (*eshop.PromotionRules, error) {
	data, err := ioutil.ReadFile("rules/rules.yml")
	if err != nil {
		return nil, status.ErrInternal.WithError(err)
	}

	var rules eshop.PromotionRules
	err = yaml.Unmarshal(data, &rules)
	if err != nil {
		return nil, status.ErrInternal.WithError(err)
	}
	return &rules, nil
}
