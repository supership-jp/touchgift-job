package models

import "strconv"

type Creative struct {
	ID               int      `db:"id" json:"id"`
	Link             string   `db:"link" json:"link,omitempty"`
	URL              string   `db:"url" json:"url"`
	TTL              int64    `db:"ttl" json:"ttl"`
	Width            float32  `db:"width" json:"width"`
	Height           float32  `db:"height" json:"height"`
	Type             string   `db:"type" json:"type"`
	Extension        string   `db:"extension" json:"extension,omitempty"`
	Duration         *int     `db:"duration" json:"duration"`
	SkipOffset       *int     `db:"skip_offset" json:"skip_offset"`
	EndCardURL       *string  `db:"end_card_url" json:"end_card_url,omitempty"`
	EndCardWidth     *float32 `db:"end_card_width" json:"end_card_width,omitempty"`
	EndCardHeight    *float32 `db:"end_card_height" json:"end_card_height,omitempty"`
	EndCardExtension *string  `db:"end_card_extension" json:"end_card_extension,omitempty"`
	EndCardLink      *string  `db:"end_card_link" json:"end_card_link,omitempty"`
}

func (c *Creative) CreateDeliveryDataCreative() *DeliveryDataCreative {
	return &DeliveryDataCreative{
		ID:               strconv.Itoa(c.ID),
		Link:             c.Link,
		URL:              c.URL,
		TTL:              c.TTL,
		Width:            c.Width,
		Height:           c.Height,
		Type:             c.Type,
		Extension:        c.Extension,
		Duration:         c.Duration,
		SkipOffset:       c.SkipOffset,
		EndCardURL:       c.EndCardURL,
		EndCardWidth:     c.EndCardWidth,
		EndCardHeight:    c.EndCardHeight,
		EndCardExtension: c.EndCardExtension,
		EndCardLink:      c.EndCardLink,
	}
}
