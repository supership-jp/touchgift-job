package models

import "strconv"

type Coupon struct {
	ID       int    `db:"coupon_id" json:"id"`
	Name     string `db:"coupon_name" json:"name"`
	Code     string `db:"coupon_code" json:"code"`
	ImageURL string `db:"coupon_image_url" json:"image_url"`
	Rate     string `db:"coupon_rate" json:"rate"`
}

func (c *Coupon) CreateDeliveryCouponData() *DeliveryCouponData {
	rate, _ := strconv.Atoi(c.Rate)
	return &DeliveryCouponData{
		ID:       c.ID,
		Name:     c.Name,
		Code:     c.Code,
		ImageURL: c.ImageURL,
		Rate:     rate,
	}
}

type Gimmick struct {
	URL  string `db:"gimmick_url" json:"gimmick_url,omitempty"`
	Code string `db:"gimmick_code" json:"gimmick_code,omitempty"`
}

type Content struct {
	CampaignID string    `db:"campaign_id" json:"campaign_id"`
	Coupons    []Coupon  `json:"coupons"`
	Gimmicks   []Gimmick `json:"gimmicks"`
}
