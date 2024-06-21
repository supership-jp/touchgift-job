package models

type Creative struct {
	ID               int     `db:"id" json:"id"`
	Link             *string `db:"link" json:"link,omitempty"`
	URL              string  `db:"url" json:"url"`
	TTL              string  `db:"ttl" json:"ttl"`
	Width            int     `db:"width" json:"width"`
	Height           int     `db:"height" json:"height"`
	Type             string  `db:"type" json:"type"`
	EndCardURL       string  `db:"end_card_url" json:"end_card_url"`
	Duration         string  `db:"duration" json:"duration"`
	SkipOffset       int     `db:"skip_offset" json:"skip_offset"`
	Extension        string  `db:"extension" json:"extension"`
	EndCardWidth     int     `db:"end_card_width" json:"end_card_width"`
	EndCardHeight    int     `db:"end_card_height" json:"end_card_height"`
	EndCardExtension string  `db:"end_card_extension" json:"end_card_extension"`
	EndCardLink      string  `db:"end_card_link" json:"end_card_link"`
}
