package models

type Creative struct {
	ID               int      `db:"id" json:"id"`
	Link             *string  `db:"link" json:"link,omitempty"`
	URL              string   `db:"url" json:"url"`
	TTL              int64    `db:"ttl" json:"ttl"`
	Width            float32  `db:"width" json:"width"`
	Height           float32  `db:"height" json:"height"`
	Type             string   `db:"type" json:"type"`
	Extension        string   `db:"extension" json:"extension"`
	Duration         *int     `db:"duration" json:"duration"`
	SkipOffset       *int     `db:"skip_offset" json:"skip_offset"`
	EndCardURL       *string  `db:"end_card_url" json:"end_card_url,omitempty"`
	EndCardWidth     *float32 `db:"end_card_width" json:"end_card_width,omitempty"`
	EndCardHeight    *float32 `db:"end_card_height" json:"end_card_height,omitempty"`
	EndCardExtension *string  `db:"end_card_extension" json:"end_card_extension,omitempty"`
	EndCardLink      *string  `db:"end_card_link" json:"end_card_link,omitempty"`
}
