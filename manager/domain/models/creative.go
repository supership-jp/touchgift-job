package models

import "time"

type Delivery struct {
	ID      int    `db:"id" json:"id"`
	Name    string `db:"name" json:"name"`
	Video   string `db:"video" json:"video"`
	EndCard string `db:"end_card" json:"end_card"`
}

type Creative struct {
	ID               int             `db:"id" json:"id"`
	OrganizationCode string          `db:"organization_code" json:"organization_code"`
	Name             string          `db:"name" json:"name"`
	Status           string          `db:"status" json:"status"`
	ClickURL         string          `db:"click_url" json:"click_url"`
	CreativeType     string          `db:"creative_type" json:"creative_type"`
	CreatedAt        time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time       `db:"updated_at" json:"updated_at"`
	LastUpdatedBy    int             `db:"last_updated_by" json:"last_updated_by"`
	Videos           []CreativeVideo `json:"videos"` // このクリエイティブに関連するビデオ
}

type CreativeVideo struct {
	CreativeID    int       `db:"creative_id" json:"creative_id"`
	VideoID       int       `db:"video_id" json:"video_id"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
	LastUpdatedBy int       `db:"last_updated_by" json:"last_updated_by"`
	Video         Video     `json:"video"` // Video 構造体への参照
}

type Video struct {
	ID               int       `db:"id" json:"id"`
	VideoURL         string    `db:"video_url" json:"video_url"`
	EndcardURL       string    `db:"endcard_url" json:"endcard_url"`
	VideoXID         string    `db:"video_xid" json:"video_xid"`
	EndcardXID       string    `db:"endcard_xid" json:"endcard_xid"`
	Height           int       `db:"height" json:"height"`
	Width            int       `db:"width" json:"width"`
	Extension        string    `db:"extension" json:"extension"`
	EndcardHeight    int       `db:"endcard_height" json:"endcard_height"`
	EndcardWidth     int       `db:"endcard_width" json:"endcard_width"`
	EndcardExtension string    `db:"endcard_extension" json:"endcard_extension"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
	LastUpdatedBy    int       `db:"last_updated_by" json:"last_updated_by"`
}
