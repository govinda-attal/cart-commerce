package eshop

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	PromoDiscount PromoType = "DISCOUNT"
	PromoPrice    PromoType = "PRICE"
)

const (
	CartStateInProgress CartState = "INPROGRESS"
	CartStateCancelled  CartState = "CANCELLED"
	CartStateSettled    CartState = "SETTLED"
)

type PromoType = string
type CartState = string

type Item struct {
	ItemName   string `json:"itemName"`
	Quantity   int    `json:"quantity"`
	UnitPrice  *Money `json:"unitPrice,omitempty"`
	Price      *Money `json:"price,omitempty"`
	PromoPrice *Money `json:"promoPrice,omitempty"`
}

type Cart struct {
	ID          string    `json:"id"`
	State       CartState `json:"state"`
	Items       []*Item   `json:"items,omitempty"`	
	TotalAmount *Money    `json:"totalAmount,omitempty"`
}

type CartItemsByTotalPrice []*Item

func (s CartItemsByTotalPrice) Len() int           { return len(s) }
func (s CartItemsByTotalPrice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s CartItemsByTotalPrice) Less(i, j int) bool { return *s[i].Price < *s[j].Price }

type Money float64

func (m *Money) MarshalJSON() ([]byte, error) {
	p := fmt.Sprintf("\"%.2f\"", float64(*m))
	return []byte(p), nil
}

func (m *Money) UnmarshalJSON(b []byte) error {
	v, err := strconv.ParseFloat(string(b[1:len(b)-1]), 64)
	if err != nil {
		return err
	}
	*m = Money(v)
	return nil
}

type PromotionItemRule struct {
	ItemName       string           `yaml:"itemName" json:"itemName"`
	Buy            int              `yaml:"buy" json:"buy"`
	PromotionItems []*PromotionItem `yaml:"promos" json:"promos"`
}

type PromotionItem struct {
	ItemName  string    `yaml:"itemName" json:"itemName"`
	PromoType PromoType `yaml:"promoType" json:"promoType"`
	PromoVal  int       `yaml:"promo" json:"promo"`
}

type PromotionRules struct {
	ItemRules []*PromotionItemRule `yaml:"items" json:"items"`
}

func (pr *PromotionRules) FindApplicablePromos(itemName string, quantity int) []*PromotionItemRule {
	var itemRules []*PromotionItemRule
	for _, rule := range pr.ItemRules {
		if strings.EqualFold(itemName, rule.ItemName) && quantity >= rule.Buy {
			itemRules = append(itemRules, rule)
		}
	}
	return itemRules
}

func (pr *PromotionItemRule) ApplyPromos(items []*Item) {
	for _, item := range items {
		if item.PromoPrice != nil {
			continue
		}
		for _, promo := range pr.PromotionItems {
			if promo.ItemName != item.ItemName {
				continue
			}
			switch promo.PromoType {
			case PromoDiscount:
				discPrice := *item.Price - Money(float64(*item.Price)*float64(promo.PromoVal)/100)
				item.PromoPrice = &discPrice
			case PromoPrice:
				if item.Quantity > pr.Buy {
					unitPrice := float64(*item.UnitPrice)
					promoPrice := Money(float64(pr.Buy)*unitPrice + float64((item.Quantity-pr.Buy)*promo.PromoVal))
					item.PromoPrice = &promoPrice
				}
			}
		}
	}
}
