package models

type Coupon struct {
	ID       int    `db:"coupon_id" json:"id"`
	Name     string `db:"coupon_name" json:"name"`
	Code     string `db:"coupon_code" json:"code"`
	ImageURL string `db:"coupon_image_url" json:"image_url"`
	Rate     string `db:"coupon_rate" json:"rate"`
}

type ContentsData struct {
	CampaignID int      `db:"campaign_id" json:"campaign_id"`
	Coupons    []Coupon `json:"coupons"`
	GimmickURL string   `db:"gimmick_url" json:"gimmick_url"`
}
