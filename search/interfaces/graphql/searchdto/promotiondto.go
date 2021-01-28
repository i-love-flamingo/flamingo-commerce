package searchdto

import (
	"flamingo.me/flamingo-commerce/v3/search/domain"
)

type (
	// PromotionDTO contains promotion data exposed via graphql
	PromotionDTO struct {
		promotion *domain.Promotion
	}
)

// WrapPromotion of search domain with PromotionDTO
func WrapPromotion(promotion *domain.Promotion) *PromotionDTO {
	return &PromotionDTO{promotion: promotion}

}

// Title of the promotion
func (p *PromotionDTO) Title() string {
	return p.promotion.Title
}

// Content of the promotion
func (p *PromotionDTO) Content() string {
	return p.promotion.Content

}

// URL of the promotion
func (p *PromotionDTO) URL() string {
	return p.promotion.URL

}

// Media of the promotion
func (p *PromotionDTO) Media() *domain.Media {
	if len(p.promotion.Media) > 0 {
		return &p.promotion.Media[0]
	}
	return nil
}
