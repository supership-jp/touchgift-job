package repository

import "time"

type UpdateCondition struct {
	CampaignID int
	Status     string
	UpdatedAt  time.Time
}
